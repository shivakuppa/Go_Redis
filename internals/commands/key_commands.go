package commands

import (
	"log"
	"path/filepath"

	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func del(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	var keysDeleted int = 0

	for _, arg := range args {
		_, exists := db.DB.Get(arg.String)
		if exists {
			db.DB.Del(arg.String)
			keysDeleted++
		}
	}

	return &resp.Value{
		Type:    resp.Integer,
		Integer: int64(keysDeleted),
	}
}

func exists(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	var keysDetected int = 0

	for _, arg := range args {
		_, exists := db.DB.Get(arg.String)
		if exists {
			keysDetected++
		}
	}

	return &resp.Value{
		Type:    resp.Integer,
		Integer: int64(keysDetected),
	}
}

func keys(value *resp.Value, state *db.AppState) *resp.Value {
	args := value.Array[1:]
	if len(args) != 1 {
		return &resp.Value{
			Type:   resp.BulkString,
			String: "ERR Invalid number of arguments for KEYS",
		}
	}

	pattern := args[0].String
	var matches []string
	for _, key := range *db.DB.GetKeys() {
		match, err := filepath.Match(pattern, key)
		if err != nil {
			log.Printf("error matching keys: (pattern: %s), (key: %s) - %v\n", pattern, key, err)
			continue
		}

		if match {
			matches = append(matches, key)
		}
	}

	reply := resp.Value{Type: resp.Array}
	for _, m := range matches {
		reply.Array = append(reply.Array, &resp.Value{Type: resp.BulkString, String: m})
	}
	return &reply
}
