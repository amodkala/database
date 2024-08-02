package wal

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "os"

    "github.com/amodkala/raft/pkg/common"

    "google.golang.org/protobuf/proto"
)

var (
    emptyBuf = make([]byte, binary.MaxVarintLen64)
)

func (w *WAL) Length() uint32 {
    return uint32(len(w.offset))
}

func (w *WAL) Clear() error {
    w.offset = []uint32{}
    w.currOffset = 0
    return os.Remove(w.filepath)
}

func (w *WAL) Write(entries ...*common.Entry) error {
    txBuf, relOffsets, err := marshalEntries(entries...)
    if err != nil {
        return fmt.Errorf("%v", err)
    }

    w.mu.Lock()
    defer w.mu.Unlock()

    f, err := os.OpenFile(w.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
    if err != nil {
        return fmt.Errorf("%v", err)
    }
    defer f.Close()

    if n, err := f.Write(txBuf); n != len(txBuf) || err != nil { 
        return fmt.Errorf("%v", err)
    }
    f.Sync()

    for _, offset := range relOffsets {
        w.offset = append(w.offset, w.currOffset + offset)
    }

    // log.Println(fmt.Sprintf("wrote %d bytes to log", len(txBuf)))

    w.currOffset += uint32(len(txBuf))

    return nil
}

func marshalEntries(entries ...*common.Entry) ([]byte, []uint32, error) {
    var txBuf []byte
    relOffset := make([]uint32, len(entries))

    opts := proto.MarshalOptions{
        Deterministic: true,
    }

    prevSize := 0
    for i, entry := range entries {

        // encode the message's size as a byte array
        sizeBuf := make([]byte, binary.MaxVarintLen64)
        size := opts.Size(entry)
        binary.PutUvarint(sizeBuf, uint64(size))

        // encode the message itself to the protobuf wire format
        b, err := opts.Marshal(entry)
        if err != nil { 
            return nil, nil, fmt.Errorf("%v", err)
        }

        txBuf = append(txBuf, sizeBuf...)
        txBuf = append(txBuf, b...)

        relOffset[i] = uint32(prevSize)
        prevSize += binary.MaxVarintLen64 + size
    }

    return txBuf, relOffset, nil
}


func (w *WAL) Read(indexes ...uint32) ([]*common.Entry, error) {

    if len(indexes) != 1 && len(indexes) != 2 {
        return nil, fmt.Errorf("Read only accepts one or two indexes")
    }

    w.mu.RLock()
    defer w.mu.RUnlock()

    if w.Length() == 0 {
        return nil, fmt.Errorf("cannot Read from empty WAL")
    }

    start := indexes[0]
    if start >= w.Length() || start < 0 {
        return nil, fmt.Errorf("start index out of range")
    }
    initOffset := w.offset[start]

    end := start
    if len(indexes) == 2 {
        end = indexes[1]
        if end >= w.Length() || end < start {
            return nil, fmt.Errorf("end index out of range")
        }
    }

    var bufSize uint32
    switch end {
    case w.Length() - 1:
        bufSize = w.currOffset - initOffset
    default:
        bufSize = w.offset[end+1] - initOffset
    }
    buf := make([]byte, bufSize)

    f, err := os.OpenFile(w.filepath, os.O_RDONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("%v", err)
    }
    defer f.Close()

    if _, err = f.ReadAt(buf, int64(initOffset)); err != nil {
        return nil, fmt.Errorf("%v", err)
    }

    return unmarshalEntries(buf)
}

// have a byte slice which we know is a sequence of alternating uint64 byte
// encodings and proto.Marshal-ed Entry types. we want to 
func unmarshalEntries(entryBytes []byte) ([]*common.Entry, error) {
    var entries []*common.Entry

    opts := proto.UnmarshalOptions{}
    for nextIndex := 0; nextIndex < len(entryBytes); {
        
        if nextIndex + binary.MaxVarintLen64 > len(entryBytes) {
            return nil, fmt.Errorf("attempted to read out of range")
        }
        sizeBuf := entryBytes[nextIndex: nextIndex + binary.MaxVarintLen64]
        if bytes.Equal(sizeBuf, emptyBuf) {
            // buffer size greater than data contained within, no need to
            // continue
            break
        }
        size, n := binary.Uvarint(sizeBuf)
        if n < 0 {
            return nil, fmt.Errorf("error reading size bytes")
        }

        b := entryBytes[nextIndex + binary.MaxVarintLen64: nextIndex + int(size) + binary.MaxVarintLen64]

        var entry common.Entry

        if err := opts.Unmarshal(b, &entry); err != nil {
            return nil, fmt.Errorf("%v", err)
        }

        entries = append(entries, &entry)
        nextIndex = nextIndex + binary.MaxVarintLen64 + int(size)
    }


    return entries, nil
}
