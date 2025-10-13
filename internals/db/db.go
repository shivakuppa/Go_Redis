package db

import "sync"

type Database struct {
	Store map[string]string
	Mu    sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		Store: map[string]string{},
		Mu:    sync.RWMutex{},
	}
}

var DB = NewDatabase()
