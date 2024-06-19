package wal

import (
    "errors"
    "io/fs"
    "os"
    "testing"
    "time"

    "github.com/amodkala/database/pkg/common"

	ts "google.golang.org/protobuf/types/known/timestamppb"
    "google.golang.org/protobuf/proto"

)

var testEntries = []*common.Entry{
    {
        RaftTerm: 0,
        Key: 1,
        CreateTime: ts.New(time.Now()),
        TxId: "tx-1",
        Value: "hello",
    },
    {
        RaftTerm: 1,
        Key: 2,
        CreateTime: ts.New(time.Now()),
        TxId: "tx-1",
        Value: "test value",
    },
    {
        RaftTerm: 1,
        Key: 3,
        CreateTime: ts.New(time.Now()),
        TxId: "tx-1",
        Value: "another test value",
    },
}

// TestMarshalEntries ensures that we can correctly encode a stream of entries
// to binary format and back again
func TestMarshalEntries(t *testing.T) {
    b, _, err := marshalEntries(testEntries...)
    if err != nil {
        t.Errorf("error marshalling entries: %v", err)
    }

    e, err := unmarshalEntries(b)
    if err != nil {
        t.Errorf("error unmarshalling entries: %v", err)
    }

    for i, entry := range e {
        if !proto.Equal(testEntries[i], entry) {
            t.Errorf("entry %d did not match", i)
        }
    }
}

func TestClear(t *testing.T) {
    filepath := "/tmp/clear"
    testWAL := New(filepath)

    if err := testWAL.Write(testEntries...); err != nil {
        t.Errorf("error while writing entries: %v", err)
    }

    testWAL.Clear()

    _, err := os.Open(filepath)

    if !errors.Is(err.(*fs.PathError).Unwrap(), fs.ErrNotExist) {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestRWSingleEntry(t *testing.T) {

    testWAL := New("/tmp/singlerw")
    defer testWAL.Clear()

    if err := testWAL.Write(testEntries[0]); err != nil {
        t.Errorf("error while writing entry: %v", err)
    }

    entries, err := testWAL.Read(0, 1)
    if err != nil {
        t.Errorf("error while reading entry: %v", err)
    }

    if len(entries) > 1 {
        t.Errorf("unexpected entries read from WAL")
    }

    if !proto.Equal(testEntries[0], entries[0]) {
        t.Errorf("entries were not equal")
    }
}

func TestMultipleWrites(t *testing.T) {

    testWAL := New("/tmp/multiwrite")
    defer testWAL.Clear()

    for _, entry := range testEntries {
        if err := testWAL.Write(entry); err != nil {
            t.Errorf("error writing entry to wal: %v", err)
        }
    }

    if int(testWAL.Length()) != len(testEntries) {
        t.Errorf("not all entries written")
    }

    entries, err := testWAL.Read(0, uint32(len(testEntries)))
    if err != nil {
        t.Errorf("error reading entries from wal: %v", err)
    }

    for i := range len(testEntries) {
        if !proto.Equal(testEntries[i], entries[i]) {
            t.Errorf("entries at index %d were not equal", i)
        }
    }
}

func TestMultipleReads(t *testing.T) {

    testWAL := New("/tmp/multiread")
    defer testWAL.Clear()

    if err := testWAL.Write(testEntries...); err != nil {
        t.Errorf("error while writing entry: %v", err)
    }

    entries := []*common.Entry{}
    for i := range uint32(len(testEntries)) {
        entry, err := testWAL.Read(i, 1)
        if err != nil {
            t.Errorf("error while reading entry: %v", err)
        }
        entries = append(entries, entry...)
    }

    for i, entry := range testEntries {
        if !proto.Equal(testEntries[i], entry) {
            t.Errorf("entries were not equal")
        }
    }
}

