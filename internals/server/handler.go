package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/shivakuppa/Go_Redis/config"
	"github.com/shivakuppa/Go_Redis/internals/commands"
	"github.com/shivakuppa/Go_Redis/internals/db"
	myio "github.com/shivakuppa/Go_Redis/internals/io"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func (s *Server) handleConnection(conn net.Conn) {
	config := config.ReadConfig("./config/redis.conf")
	state := db.NewAppState(config)
	if config.AOFenabled {
		log.Println("syncing AOF records")
		aofSync(state.Aof)
	}

	defer conn.Close()
	w := myio.NewRespWriter(conn)

	for {
		value, err := resp.Deserialize(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Client disconnected")
				return
			}

			fmt.Printf("Error deserializing request: %v\n", err)
			errVal := &resp.Value{
				Type:   resp.SimpleError,
				String: "ERR invalid request",
			}
			_ = w.Write(errVal)
			_ = w.Flush()
			return
		}

		reply := commands.HandleCommand(conn, value, state)
		w.Write(reply)
		w.Flush()
	}
}

func aofSync(aof *db.Aof) {
	file := aof.File
	reader := bufio.NewReader(file)
	replayState := db.NewAppState(aof.Config)
	replayState.Config.AOFenabled = false

	for {
		value, err := resp.Deserialize(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // reached end of file
			}
			fmt.Println("Error reading AOF file:", err)
			continue // skip bad entries instead of breaking everything
		}

		commands.ResolveCommand(value, replayState)
	}

	replayState.Config.AOFenabled = true
	fmt.Println("AOF replay complete — state restored successfully.")
}
