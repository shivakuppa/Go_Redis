package commands

import (
	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func save(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	if len(args) != 0 {
		return &resp.Value{
			Type:   resp.BulkString,
			String: "ERR Invalid number of arguments for KEYS",
		}
	}

	db.SaveRDB(state)
	return &resp.Value{
		Type: resp.SimpleString,
		String: "OK",
	}
}

func bgsave(value *resp.Value, state *db.AppState) *resp.Value {
	if state.BgSaveRunning {
		return &resp.Value{
			Type: resp.SimpleError,
			String: "ERR background saving already in progress",
		}
	}

	args := value.Array[1:]
	if len(args) != 0 {
		return &resp.Value{
			Type:   resp.BulkString,
			String: "ERR Invalid number of arguments for KEYS",
		}
	}
	
	copy := db.DB.GetItems()
	state.DBCopy = *copy
	go func() {
		defer func() {
			state.BgSaveRunning = false
			state.DBCopy = nil
		}()

		db.SaveRDB(state)
	}()

	db.SaveRDB(state)
	return &resp.Value{
		Type: resp.SimpleString,
		String: "OK",
	}
}

func flushdb(value *resp.Value, state *db.AppState) *resp.Value{
	db.DB.Reset()
	return &resp.Value{
		Type: resp.SimpleString,
		String: "OK",
	}
}

func dbsize(value *resp.Value, state *db.AppState) *resp.Value{
	length := db.DB.GetLen()
	return &resp.Value{
		Type: resp.Integer,
		Integer: int64(length),
	}
}