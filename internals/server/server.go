package server

import (
	"log/slog"
	"net"
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

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("Cannot listen on port", "addr", s.ListenAddr, "error", err)
		return err
	}

	slog.Info("goredis server running", "listenAddr", s.ListenAddr)
	defer listener.Close()
	s.Listener = listener

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			slog.Error("Error during accepting")
			continue
		}

		go s.handleConnection(conn)
	}
}
