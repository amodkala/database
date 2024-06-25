package database

import (
    "fmt"
    "log"
    "time"

    "github.com/amodkala/database/pkg/common"
    tx "github.com/amodkala/database/pkg/transaction"

    "google.golang.org/protobuf/types/known/timestamppb"
)

func (c Client) Write(records ...Record) error {

    entries := []*common.Entry{}

    for _, record := range records {
        entries = append(entries, &common.Entry{
            Key: record.Key,
            Value: record.Value,
            CreateTime: timestamppb.New(time.Now()),
        })
    }

    tx := tx.New("tx-1", entries...)

    if leaderId, err := c.cm.Replicate(tx); err != nil {
        if leaderId == "" {
            return fmt.Errorf("error writing records to database: %v", err)
        }
        // TODO: implement redirection
        log.Printf("redirect request to %s\n", leaderId)
    }

    // entries have been flagged for commitment, write to LSM tree

    return nil
}

func (c Client) Read(key uint32) (value string, err error)
