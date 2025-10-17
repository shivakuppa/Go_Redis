package db

import "sync"

type Database struct {
	store map[string]*Item
	mu    sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		store: map[string]*Item{},
		mu:    sync.RWMutex{},
	}
}

type DatabaseInterface interface {
	Get(key string) (*Item, bool)
	Set(key string, val string)
	Del(key string)
	GetKeys() []string
	GetItems() map[string]*Item
	GetLen() int 
	Reset()
}

func (d *Database) Get(key string) (*Item, bool) {
	d.mu.RLock()
	val, ok := d.store[key]
	d.mu.RUnlock()
	return val, ok
}

func (d *Database) Set(key string, value string) {
	d.mu.Lock()
	d.store[key] = makeItem(value)
	d.mu.Unlock()
}
func (d *Database) Del(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.store, key)
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

func (d *Database) GetItems() *map[string]*Item {
    d.mu.RLock()
    defer d.mu.RUnlock()

    items := make(map[string]*Item, len(d.store))
    for k, v := range d.store {
        if v != nil {
            copyItem := *v
            items[k] = &copyItem
        } else {
            items[k] = nil
        }
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
	d.store = map[string]*Item{}
	d.mu.Unlock()
}

func (db *Database) TryExpire(k string, i *Item) bool {
	if i.shouldExpire() {
		DB.mu.Lock()
		DB.Del(k)
		DB.mu.Unlock()
		// state.generalStats.expired_keys++
		return true
	}
	return false
}

var DB = NewDatabase()
