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

func ReadConfig(fn string) *Config {
	config := NewConfig()

	file, err := os.Open(fn)
	if err != nil {
		fmt.Printf("cannot read file: %s - using default config\n", fn)
		return config
	}

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		line := scan.Text()
		parseLine(line, config)
	}

	if err := scan.Err(); err != nil {
		fmt.Println("error scanning config file: ", err)
		return config
	}

	if config.Dir != "" {
		os.MkdirAll(config.Dir, 0755)
	}

	fmt.Printf("%+v\n", config)

	return config
}

func parseLine(line string, config *Config) {
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
		config.RDB = append(config.RDB, snapshot)

	case "dbfilename":
		config.RDBfn = args[1]

	case "appendfilename":
		config.AOFfn = args[1]

	case "appendfsync":
		config.AOFfsync = FSyncMode(args[1])

	case "dir":
		config.Dir = args[1]

	case "appendonly":
		if args[1] == "yes" {
			config.AOFenabled = true
		} else {
			config.AOFenabled = false
		}
	}

}
