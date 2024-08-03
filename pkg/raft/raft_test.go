package raft

import (
    "os"
	"testing"
	"time"

	"github.com/amodkala/raft/pkg/common"
	tx "github.com/amodkala/raft/pkg/transaction"

	ts "google.golang.org/protobuf/types/known/timestamppb"
)

var (

    raftAddress = os.Getenv("RAFT_ADDRESS")

	testEntries = []*common.Entry{
		{
			RaftTerm:   0,
			Key:        1,
			CreateTime: ts.New(time.Now()),
			TxId:       "tx-1",
			Value:      "hello",
		},
		{
			RaftTerm:   1,
			Key:        2,
			CreateTime: ts.New(time.Now()),
			TxId:       "tx-1",
			Value:      "test value",
		},
		{
			RaftTerm:   1,
			Key:        3,
			CreateTime: ts.New(time.Now()),
			TxId:       "tx-1",
			Value:      "another test value",
		},
	}

	testTx = tx.New(testEntries...)
)

func TestStart(t *testing.T) {
	cm := New("test-start", "localhost:8080")
	go func() {
		if err := cm.Start(raftAddress); err != nil {
			t.Errorf("error on raft cm start: %v", err)
		}
	}()
	time.Sleep(5 * time.Second)
}

func TestReplicate(t *testing.T) {
	cm := New("test-replicate", "localhost:8081")
	go func() {
		if err := cm.Start(raftAddress); err != nil {
			t.Errorf("error on raft cm start: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	if _, err := cm.Replicate(testTx); err != nil {
		t.Errorf("error during replication, %v", err)
	}
}
