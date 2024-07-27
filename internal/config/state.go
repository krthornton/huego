package config

import "huego/internal/hue"

type ProgramState struct {
	Config *Configuration
	Conn   *hue.HueConnection
}

func NewProgramState(conf *Configuration) *ProgramState {
	// create new connection and device db objects
	conn := hue.NewHueConnection()

	// instantiate app state object to pass around
	state := &ProgramState{
		Config: conf,
		Conn:   conn,
	}

	return state
}
