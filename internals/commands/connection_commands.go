package commands

import (
	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func command(v *resp.Value, state *db.AppState) *resp.Value {
	return &resp.Value{
		Type:   resp.SimpleString,
		String: "OK",
	}
}
