package data_persistance

import (
	"fmt"
	"os"
	"path"

	"github.com/shivakuppa/Go_Redis/config"
	"github.com/shivakuppa/Go_Redis/internals/server"
)

type Aof struct {
	Writer 		*server.Writer
	File 		*os.File
	Config 		*config.Config
}

func NewAOF(conf *config.Config) *Aof {
	aof := Aof{Config: conf}

	fp := path.Join(conf.Dir, conf.AOFfn)
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Cannot open file: ", fp)
		return &aof
	}

	aof.Writer = server.NewWriter(file)
	aof.File = file

	return &aof
}

