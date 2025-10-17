package commands

import (
	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func set(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	if len(args) != 2 {
		return &resp.Value{
			Type: resp.SimpleError, 
			String: "ERR Invalid number of arguments for SET",
		}
	}

	key := args[0].String
	val := args[1].String

	db.DB.Set(key, val)

	if state.Config.AOFenabled {
		state.Aof.Writer.Write(value)

		if state.Config.AOFfsync == "always" {
			state.Aof.Writer.Flush()
		}
	}

	if len(state.Config.RDB) > 0 {
		db.IncrRDBTrackers()
	}

	return &resp.Value{
		Type:   resp.SimpleString,
		String: "OK",
	}
}

func get(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	if len(args) != 1 {
		return &resp.Value{
			Type: resp.SimpleError, 
			String: "ERR Invalid number of arguments for GET",
		}
	}

	key := args[0].String
	val, ok := db.DB.Get(key)

	if !ok {
		return &resp.Value{
			Type:   resp.Null,
			IsNull: true,
		}
	}

	return &resp.Value{
		Type:   resp.BulkString,
		String: val.Value,
	}
}
