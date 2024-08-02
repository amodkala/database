package transaction 

import (
    "github.com/amodkala/raft/pkg/common"
)

type Tx struct {
    id string
    entries []*common.Entry
}

func New(id string, entries ...*common.Entry) Tx {

    for _, entry := range entries {
        entry.TxId = id
    }

    return Tx {
        id,
        entries,
    }

}
