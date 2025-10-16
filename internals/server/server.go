package server

import (
	"log/slog"
	"net"
	"sync"

	"github.com/shivakuppa/Go_Redis/internals/client"
	"github.com/shivakuppa/Go_Redis/internals/db"
)

const defaultListenAddr = ":6379"

type Server struct {
	ListenAddr string
	Listener   net.Listener
}

func NewServer(listenAddr string) *Server {
	if len(listenAddr) == 0 {
		listenAddr = defaultListenAddr
	}

	return &Server{
		ListenAddr: listenAddr,
	}
}

func (s *Server) Start(state *db.AppState) error {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("Cannot listen on port", "addr", s.ListenAddr, "error", err)
		return err
	}

	slog.Info("goredis server running", "listenAddr", s.ListenAddr)
	defer listener.Close()
	s.Listener = listener

	return s.acceptLoop(state)
}

func (s *Server) acceptLoop(state *db.AppState) error {
	var wg sync.WaitGroup
	defer wg.Wait() // ✅ wait for all connections only when shutting down

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			slog.Error("Error during accepting", "error", err)
			continue
		}

		c := client.NewClient(conn)

		wg.Add(1)
		go func(c *client.Client, appstate *db.AppState) {
			defer wg.Done()
			s.handleConnection(c, appstate)
		}(c, state)
	}
}
