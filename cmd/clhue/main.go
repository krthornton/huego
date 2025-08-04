package main

import (
	"errors"
	"huego"

	tea "github.com/charmbracelet/bubbletea"
)

type programState struct {
	config *configuration
	conn   *huego.HueConnection
}

func newProgramState(conf *configuration) *programState {
	// create new connection and device db objects
	conn := huego.NewHueConnection()

	// instantiate app state object to pass around
	state := &programState{
		config: conf,
		conn:   conn,
	}

	return state
}

func main() {
	// attempt to load configuration
	conf, err := loadConfiguration()
	if err != nil {
		var confErr *configFileNotExists
		if !errors.As(err, &confErr) {
			panic(err.Error())
		}

		// simply create a blank new config if none exists
		conf = newConfiguration()
	}

	// init program state object to pass between menus
	state := newProgramState(&conf)

	// setup TUI and start its main loop
	mainModel := initMainModel(state)
	program := tea.NewProgram(mainModel)
	if _, err := program.Run(); err != nil {
		panic(err.Error())
	}

	// main loop has exited, let's save config back to disk
	conf.saveConfiguration()
}
