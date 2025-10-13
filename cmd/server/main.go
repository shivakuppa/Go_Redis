package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/shivakuppa/Go_Redis/config"
	"github.com/shivakuppa/Go_Redis/internals/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	config.ReadConf("./config/redis.conf")

	port := flag.String("port", "6379", "Port to start the Redis server on")
	flag.Parse()

	sConfig := server.Config{ListenAddr: ":" + *port}

	s := server.NewServer(sConfig)
	s.Start()
}
