package commands

import (
	"strconv"
	"time"

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

func flushdb(value *resp.Value, state *db.AppState) *resp.Value {
	db.DB.Reset()
	return &resp.Value{
		Type: resp.SimpleString,
		String: "OK",
	}
}

func dbsize(value *resp.Value, state *db.AppState) *resp.Value {
	length := db.DB.GetLen()
	return &resp.Value{
		Type: resp.Integer,
		Integer: int64(length),
	}
}

func expire(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	if len(args) != 2 {
		return &resp.Value{Type: resp.SimpleError, String: "ERR invalid number of arguments for 'EXPIRE' command"}
	}

	k := args[0].String
	exp := args[1].String

	expSecs, err := strconv.Atoi(exp)
	if err != nil {
		return &resp.Value{Type: resp.SimpleError, String: "ERR invalid expiry value"}
	}

	key, ok := db.DB.Get(k)
	if !ok {
		return &resp.Value{Type: resp.Integer, Integer: 0}
	}
	key.Expires = time.Now().Add(time.Second * time.Duration(expSecs))

	return &resp.Value{Type: resp.Integer, Integer: 1}
}

func ttl(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	if len(args) != 1 {
		return &resp.Value{Type: resp.SimpleError, String: "ERR invalid number of arguments for 'TTL' command"}
	}

	k := args[0].String

	item, ok := db.DB.Get(k)
	if !ok {
		return &resp.Value{Type: resp.Integer, Integer: -2}
	}

	expires := item.Expires
	if expires.Unix() == db.UNIX_TS_EPOCH {
		return &resp.Value{Type: resp.Integer, Integer: -1}
	}

	expired := db.DB.TryExpire(k, item)
	if expired {
		return &resp.Value{Type: resp.Integer, Integer: -2}
	}

	expSecs := int(time.Until(expires).Seconds())
	return &resp.Value{Type: resp.Integer, Integer: int64(expSecs)}
}