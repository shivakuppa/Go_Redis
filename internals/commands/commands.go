package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/shivakuppa/Go_Redis/internals/resp"
)

type CmdHandler func(*resp.Value) *resp.Value

var CmdHandlers = map[string]CmdHandler{
	CMD_COMMAND: command,
	CMD_SET: set,
	CMD_GET: get,
}

func HandleCommand(conn net.Conn, v *resp.Value) *resp.Value{
	cmd := v.Array[0].String
	handler, ok := CmdHandlers[strings.ToUpper(cmd)]
	if !ok {
		fmt.Println("Invalid command: ", cmd)
		return &resp.Value{
			Type: 	resp.Null,
			IsNull: true,
		}
	}

	reply := handler(v)
	return reply
}