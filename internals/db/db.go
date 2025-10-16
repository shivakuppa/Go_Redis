package db

import "sync"

type Database struct {
	store map[string]string
	mu    sync.RWMutex
}

type DatabaseInterface interface {
	Get(key string) (string, bool)
	Set(key string, val string)
	Del(key string)
	GetKeys() []string
	GetItems() map[string]string
	GetLen() int 
	Reset()
}

func (d *Database) Get(key string) (string, bool) {
	d.mu.RLock()
	val, ok := d.store[key]
	d.mu.RUnlock()
	return val, ok
}

func (d *Database) Set(key string, val string) {
	d.mu.Lock()
	d.store[key] = val
	d.mu.Unlock()
}

func (d *Database) Del(key string) {
	d.mu.Lock()
	delete(d.store, d.store[key])
	d.mu.Unlock()
}

func NewDatabase() *Database {
	return &Database{
		store: map[string]string{},
		mu:    sync.RWMutex{},
	}
}

func (d *Database) GetKeys() *[]string {
    d.mu.RLock()
    defer d.mu.RUnlock()

    keys := make([]string, 0, len(d.store))
    for k := range d.store {
        keys = append(keys, k)
    }

    return &keys
}

func (d *Database) GetItems() *map[string]string {
    d.mu.RLock()
    defer d.mu.RUnlock()

    items := make(map[string]string, len(d.store))
    for k, v := range d.store {
        items[k] = v
    }

    return &items
}

func (d *Database) GetLen() int {
	d.mu.RLock()
	length := len(d.store)
	d.mu.RUnlock()
	return length
}

func (d *Database) Reset() {
	d.mu.Lock()
	d.store = map[string]string{}
	d.mu.Unlock()
}

var DB = NewDatabase()
