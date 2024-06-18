package wal

import (
    "log"
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
}

func TestMarshalEntries(t *testing.T) {
    b, _, _ := marshalEntries(testEntries...)
    e, _ := unmarshalEntries(b)
    for i, entry := range e {
        if !proto.Equal(testEntries[i], entry) {
            log.Println("entry did not match")
        }
    }
}
