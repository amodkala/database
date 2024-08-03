package wal

import "sync"

type WAL struct {
	mu         sync.RWMutex
	filepath   string
	offset     []uint32
	currOffset uint32
}

func New(filepath string) *WAL {
	return &WAL{
		mu:         sync.RWMutex{},
		filepath:   filepath,
		offset:     []uint32{},
		currOffset: 0,
	}
}
