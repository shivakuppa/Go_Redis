package server

import (
	// "errors"
	"fmt"
	"net"
	// "io"
)

func (s *Server) handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		// if errors.Is(err, io.EOF) {
		// 	return
		// }
		fmt.Println(string(buffer))
		conn.Write([]byte("+OK\r\n"))
	}
}