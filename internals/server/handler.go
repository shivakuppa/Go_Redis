package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/shivakuppa/Go_Redis/internals/client"
	"github.com/shivakuppa/Go_Redis/internals/commands"
	"github.com/shivakuppa/Go_Redis/internals/db"
	myio "github.com/shivakuppa/Go_Redis/internals/io"
	"github.com/shivakuppa/Go_Redis/internals/resp"
)

func (s *Server) handleConnection(c *client.Client, state *db.AppState) {
	defer c.Conn.Close()
	w := myio.NewRespWriter(c.Conn)

	if state.Config.Requirepass {
		log.Println(state.Config.Password)
		authenticate(c, state, w)
	}

	if state.Config.AOFenabled {
		log.Println("syncing AOF records")
		aofSync(state.Aof)
	}

	if len(state.Config.RDB) > 0 {
		db.SyncRDB(state)
		db.InitRDBTrackers(state)
	}

	for {
		value, err := resp.Deserialize(c.Conn)
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

		reply := commands.HandleCommand(c.Conn, value, state)
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

func authenticate(c *client.Client, state *db.AppState, w *myio.RespWriter) {
	log.Println(state.Config.Password)

	value, err := resp.Deserialize(c.Conn)
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

	reply := commands.HandleCommand(c.Conn, value, state)
	w.Write(reply)
	w.Flush()

	askPassword := &resp.Value{
		Type:   resp.SimpleString,
		String: "Please enter password:",
	}

	authenticatedMsg := &resp.Value{
		Type:   resp.SimpleString,
		String: "OK - Password authenticated",
	}

	retryPassword := &resp.Value{
		Type:   resp.SimpleError,
		String: "ERR invalid password, please try again",
	}

	w.Write(askPassword)
	w.Flush()

	for {
		// Wait for password input
		reply, err := resp.Deserialize(c.Conn)
		password := reply.Array[0].String
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("Client disconnected during authentication")
				return
			}
			log.Printf("Error while reading password: %v\n", err)
			continue
		}

		log.Println(password)

		if password == state.Config.Password {
			c.Authenticated = true
			w.Write(authenticatedMsg)
			w.Flush()
			log.Printf("Client %v authenticated successfully\n", c.Conn.RemoteAddr())
			return
		}

		// Invalid password
		w.Write(retryPassword)
		w.Flush()
	}
}
