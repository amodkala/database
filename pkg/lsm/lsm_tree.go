package lsm

// LSMTree is an implementation of a multicomponent LSM Tree 
type LSMTree struct {
    id int
}

// New accepts an identifier which is used as a prefix for any files associated
// with this LSMTree (logs, SSTables, etc.) and returns a pointer to an LSM Tree
func New(id int) *LSMTree {
    return &LSMTree{
        id,
    }
}

// Write accepts information about a record; its key and the data to be
// associated with that key. Due to the immutability of LSM Trees this is more
// akin to an "upsert" operation from the client's point of view
func (t *LSMTree) Write(key int, data []byte) error { return nil }

// Read accepts a record's key and returns the data most recently associated 
// with that key if it exists in storage, else it returns an empty slice of
// bytes and "false" for the second return value.
func (t *LSMTree) Read(key int) (data []byte, ok bool) { return []byte{}, true }

// Delete accepts a record's key and writes to storage a record associated with
// that key where the data value is a tombstone. During compaction this
// tombstone will be propagated and subsequent calls to Read using this key will
// not return data until another record is written with this key.
func (t *LSMTree) Delete(key int) error { return nil }
