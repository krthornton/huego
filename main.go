package main

import (
	"errors"
	"huego/internal/config"
	"huego/internal/menu"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// attempt to load configuration
	conf, err := config.LoadConfiguration()
	if err != nil {
		var confErr *config.ConfigFileNotExists
		if !errors.As(err, &confErr) {
			panic(err.Error())
		}

		// simply create a blank new config if none exists
		conf = config.NewConfiguration()
	}

	// init program state object to pass between menus
	state := config.NewProgramState(&conf)

	// setup TUI and start its main loop
	mainModel := menu.InitMainModel(state)
	program := tea.NewProgram(mainModel)
	if _, err := program.Run(); err != nil {
		panic(err.Error())
	}

	// main loop has exited, let's save config back to disk
	config.SaveConfiguration(conf)
}
