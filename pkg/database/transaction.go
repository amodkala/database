package database

import (
    "github.com/amodkala/database/pkg/common"
    ts "google.golang.org/protobuf/types/known/timestamppb"
)

type Tx struct {
    entries common.Entry
}

func NewTx()
