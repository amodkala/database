package transaction

import (
    "github.com/amodkala/database/pkg/common"
)

func (t Tx) ID() string { return t.id }

func (t Tx) Length() uint32 { return uint32(len(t.entries)) }

func (t Tx) Entries() []*common.Entry { return t.entries }
