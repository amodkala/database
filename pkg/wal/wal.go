package wal

import (
    "os"
    "sync"
)

type WAL struct {
    mu sync.RWMutex
    file *os.File
}

func New(filepath string) *WAL 

func (w *WAL) Write(data []byte) error
