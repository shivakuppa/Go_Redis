package commands

import (
	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func set(v *resp.Value) *resp.Value {
	args := v.Array[1:]
	if len(args) != 2 {
		return &resp.Value{Type: resp.SimpleError, String: "ERR Invalid number of arguments for SET"}
	}

	key := args[0].String
	val := args[1].String

	db.DB.Mu.Lock()
	db.DB.Store[key] = val
	db.DB.Mu.Unlock()

	return &resp.Value{
		Type:   resp.SimpleString,
		String: "OK",
	}
}

func get(v *resp.Value) *resp.Value {
	args := v.Array[1:]
	if len(args) != 1 {
		return &resp.Value{Type: resp.SimpleError, String: "ERR Invalid number of arguments for GET"}
	}

	key := args[0].String

	db.DB.Mu.RLock()
	val, ok := db.DB.Store[key]
	db.DB.Mu.RUnlock()

	if !ok {
		return &resp.Value{
			Type:   resp.Null,
			IsNull: true,
		}
	}

	return &resp.Value{
		Type:   resp.BulkString,
		String: val,
	}
}
