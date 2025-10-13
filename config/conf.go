package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Dir        string
	RDB        []RDBSnapshot
	RDBfn      string
	AOFenabled bool
	AOFfn      string
	AOFfsync   FSyncMode
}

type RDBSnapshot struct {
	Secs        int
	KeysChanged int
}

type FSyncMode string

const (
	Always   FSyncMode = "always"
	EverySec FSyncMode = "everysec"
	No       FSyncMode = "no"
)

func NewConfig() *Config {
	return &Config{}
}

func ReadConf(fn string) *Config {
	conf := NewConfig()

	file, err := os.Open(fn)
	if err != nil {
		fmt.Printf("cannot read file: %s - using default config\n", fn)
		return conf
	}

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		line := scan.Text()
		parseLine(line, conf)
	}

	if err := scan.Err(); err != nil {
		fmt.Println("error scanning config file: ", err)
		return conf
	}

	if conf.Dir != "" {
		os.MkdirAll(conf.Dir, 0755)
	}

	fmt.Printf("%+v\n", conf)

	return conf
}

func parseLine(line string, conf *Config) {
	args := strings.Split(line, " ")
	cmd := args[0]
	switch cmd {
	case "save":
		secs, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("invalid secs")
			return
		}

		keysChanged, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("invalid keys")
			return
		}

		snapshot := RDBSnapshot{
			Secs:        secs,
			KeysChanged: keysChanged,
		}
		conf.RDB = append(conf.RDB, snapshot)

	case "dbfilename":
		conf.RDBfn = args[1]

	case "appendfilename":
		conf.AOFfn = args[1]

	case "appendfsync":
		conf.AOFfsync = FSyncMode(args[1])

	case "dir":
		conf.Dir = args[1]

	case "appendonly":
		if args[1] == "yes" {
			conf.AOFenabled = true
		} else {
			conf.AOFenabled = false
		}
	}

}
