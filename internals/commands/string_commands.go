package commands

import (
	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func set(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	if len(args) != 2 {
		return &resp.Value{Type: resp.SimpleError, String: "ERR Invalid number of arguments for SET"}
	}

	key := args[0].String
	val := args[1].String

	db.DB.Mu.Lock()
	db.DB.Store[key] = val

	if state.Config.AOFenabled {
		state.Aof.Writer.Write(value)

		if state.Config.AOFfsync == "always" {
			state.Aof.Writer.Flush()
		}
	}

	if len(state.Config.RDB) > 0 {
		db.IncrRDBTrackers()
	}

	db.DB.Mu.Unlock()

	return &resp.Value{
		Type:   resp.SimpleString,
		String: "OK",
	}
}

func get(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
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
