package client

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/shivakuppa/Go_Redis/internals/resp"
	myio "github.com/shivakuppa/Go_Redis/internals/io"
)

type Client struct {
	Conn          net.Conn
	Authenticated bool
}

func NewClient(conn net.Conn) *Client {
	return &Client{Conn: conn}
}

func (c *Client) writeMonitorLog(value *resp.Value) {
	log.Println("relaying command to monitor: ", c.Conn.LocalAddr().String())

	msg := fmt.Sprintf("%d [%s]", time.Now().Unix(), c.Conn.LocalAddr().String())
	for _, val := range value.Array {
		msg += fmt.Sprintf("\"%s\"", val.String)
	}

	reply := resp.Value{Type: resp.BulkString, String: msg}
	w := myio.NewRespWriter(c.Conn)
	w.Write(&reply)
	w.Flush()
}