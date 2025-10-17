package db

import (
	"github.com/shivakuppa/Go_Redis/config"
)

type AppState struct {
	Config 			*config.Config
	Aof    			*Aof
	BgSaveRunning  	bool
	DBCopy			map[string]*Item
}

func NewAppState(config *config.Config) *AppState {
	state := AppState{
		Config: config,
	}

	if config.AOFenabled {
		state.Aof = NewAOF(config)
	}

	return &state
}
