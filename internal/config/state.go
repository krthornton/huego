package config

import "huego/internal/hue"

type ProgramState struct {
	Config   *Configuration
	Conn     *hue.HueConnection
	DeviceDb *DeviceDb
}

func NewProgramState(conf *Configuration) *ProgramState {
	// create new connection and device db objects
	conn := hue.NewHueConnection()
	db := MakeNewDeviceDb(conn)

	// instantiate app state object to pass around
	state := &ProgramState{
		Config:   conf,
		Conn:     conn,
		DeviceDb: db,
	}

	return state
}
