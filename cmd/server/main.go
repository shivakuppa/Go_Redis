package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/shivakuppa/Go_Redis/config"
	"github.com/shivakuppa/Go_Redis/internals/db"
	"github.com/shivakuppa/Go_Redis/internals/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	port := flag.String("port", "6379", "Port to start the Redis server on")
	flag.Parse()

	conf := config.ReadConfig("./config/redis.conf")
	state := db.NewAppState(conf)

	s := server.NewServer(":" + *port)
	s.Start(state)
}
