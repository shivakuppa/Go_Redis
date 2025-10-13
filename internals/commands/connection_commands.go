package commands

import (
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func command(v *resp.Value) *resp.Value {
	return &resp.Value{
		Type:   resp.SimpleString,
		String: "OK",
	}
}
