package db

import (
	"github.com/shivakuppa/Go_Redis/config"
)

type AppState struct {
	Config *config.Config
	Aof    *Aof
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
