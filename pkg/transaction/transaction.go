package transaction 

import (
    "github.com/google/uuid"
    "github.com/amodkala/raft/pkg/common"
)

type Tx struct {
    id string
    entries []*common.Entry
}

func New(entries ...*common.Entry) Tx {

    id := uuid.NewString()

    for _, entry := range entries {
        entry.TxId = id
    }

    return Tx {
        id,
        entries,
    }

}
