package server

import (
	"net"
	"log/slog"
)

const defaultListenAddr = ":6379"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln net.Listener
}

func NewServer(config Config) *Server {
	if len(config.ListenAddr) == 0 {
		config.ListenAddr = defaultListenAddr
	}
	
	return &Server{
		Config: config,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("Cannot listen on port", "addr", s.ListenAddr, "error", err)
		return err
	}
	
	slog.Info("goredis server running", "listenAddr", s.ListenAddr)
	defer ln.Close()
	s.ln = ln
	
	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Error during accepting")
			continue
		}

		go s.handleConnection(conn)
	}
}
