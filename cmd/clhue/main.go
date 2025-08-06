package main

import (
	"errors"
	"huego"
	keyman "huego/internal/keyman"

	tea "github.com/charmbracelet/bubbletea"
)

type programState struct {
	keyMan *keyman.KeyManager
	conn   *huego.HueConnection
}

func newProgramState(keyMan *keyman.KeyManager) *programState {
	// create new connection and device db objects
	conn := huego.NewHueConnection()

	// instantiate app state object to pass around
	state := &programState{
		keyMan: keyMan,
		conn:   conn,
	}

	return state
}

func main() {
	// attempt to load saved API keys
	keyMan := keyman.NewKeyManager()
	keyStoreFilePath := keyman.GetDefaultKeyStoreFilePath()
	err := keyMan.LoadFromKeyStore(keyStoreFilePath)
	if err != nil {
		var confErr *keyman.KeyStoreNotExistsError
		if !errors.As(err, &confErr) {
			panic(err.Error())
		}

		// simply create a blank new key store if none exists
		keyMan = keyman.NewKeyManager()
	}

	// init program state object to pass between menus
	state := newProgramState(&keyMan)

	// setup TUI and start its main loop
	mainModel := initMainModel(state)
	program := tea.NewProgram(mainModel)
	if _, err := program.Run(); err != nil {
		panic(err.Error())
	}

	// main loop has exited, let's save key store back to disk
	err = keyMan.SaveToKeyStore(keyStoreFilePath)
	if err != nil {
		panic(err.Error())
	}
}
