package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/shivakuppa/Go_Redis/internals/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	port := flag.String("port", "6379", "Port to start the Redis server on")
	flag.Parse()

	s := server.NewServer(":" + *port)
	s.Start()
}
