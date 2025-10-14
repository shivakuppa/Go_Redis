package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

type CmdHandler func(*resp.Value, *db.AppState) *resp.Value

var CmdHandlers = map[string]CmdHandler{
	CMD_COMMAND: command,
	CMD_SET:     set,
	CMD_GET:     get,
}

func HandleCommand(conn net.Conn, value *resp.Value, state *db.AppState) *resp.Value {
	cmd := value.Array[0].String
	handler, ok := CmdHandlers[strings.ToUpper(cmd)]
	if !ok {
		fmt.Println("Invalid command: ", cmd)
		return &resp.Value{
			Type:   resp.Null,
			IsNull: true,
		}
	}

	reply := handler(value, state)
	return reply
}

func ResolveCommand(value *resp.Value, state *db.AppState) {
	cmd := value.Array[0].String
	handler, ok := CmdHandlers[strings.ToUpper(cmd)]
	if !ok {
		fmt.Println("Invalid command: ", cmd)
	}
	handler(value, state)
}
