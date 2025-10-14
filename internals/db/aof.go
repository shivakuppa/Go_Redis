package db

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/shivakuppa/Go_Redis/config"
	myio "github.com/shivakuppa/Go_Redis/internals/io"
)

type Aof struct {
	Writer *myio.RespWriter
	File   *os.File
	Config *config.Config
}

func NewAOF(conf *config.Config) *Aof {
	aof := Aof{Config: conf}

	fp := path.Join(conf.Dir, conf.AOFfn)
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Cannot open file: ", fp)
		return &aof
	}

	aof.Writer = myio.NewRespWriter(file)
	aof.File = file

	go func() {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()

		for range t.C {
			if err := aof.Writer.Flush(); err != nil {
				fmt.Println("AOF flush error:", err)
			}
		}
	}()

	return &aof
}
