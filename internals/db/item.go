package db

import "time"

var UNIX_TS_EPOCH int64 = -62135596800

type Item struct {
	Value      string
	Expires    time.Time
	LastAccess time.Time
	Accesses   int
}

// makeItem creates a new Item with the given value and TTL (in seconds)
func makeItem(value string) *Item {
	now := time.Now()
	item := &Item{
		Value:      value,
		LastAccess: now,
		Accesses:   0,
	}

	return item
}

func (item *Item) shouldExpire() bool {
	return item.Expires.Unix() != UNIX_TS_EPOCH && time.Until(item.Expires).Seconds() <= 0
}

func (item *Item) approxMemUsage(name string) int64 {
	stringHeader := 16
	expHeader := 24
	mapEntrySize := 32

	return int64(stringHeader + len(name) + stringHeader + len(item.Value) + expHeader + mapEntrySize)
}