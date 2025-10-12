package server

import (
	// "errors"
	// "bufio"
	"fmt"
	"net"
	"errors"
	"io"

	"github.com/shivakuppa/Go_Redis/internals/resp"
	"github.com/shivakuppa/Go_Redis/internals/commands"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	w := NewWriter(conn)

	for {
		value, err := resp.Deserialize(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Client disconnected")
				return
			}

			fmt.Printf("Error deserializing request: %v\n", err)
			// Send RESP error back to client
			errVal := &resp.Value{
				Type:   resp.SimpleError,
				String: "ERR invalid request",
			}
			_ = w.Write(errVal)
			_ = w.Flush()
			return
		}

		reply := commands.HandleCommand(conn, value)
		w.Write(reply)
		w.Flush()
	}
}
