package main

import (
	"net"
	"log"
)

type Config struct {

}

type Server struct {
	Config
	ln net.Listener
}

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Cannot listen on Port 6379")
	}
	if l != nil {
		return
	}
}