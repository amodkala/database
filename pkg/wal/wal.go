package wal

import (
    "os"
    "sync"
)

type WAL struct {
    mu sync.Mutex
    file *os.File
}

func New(filepath string) *WAL 
