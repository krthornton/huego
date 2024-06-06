package menu

import (
	"fmt"

	"huego-cli/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

type devicesModel struct {
	cursor int
	state  *config.ProgramState
}

func (m devicesModel) Init() tea.Cmd {
	return func() tea.Msg {
		// attempt to fetch device data
		m.state.DeviceDb.FetchDevices()

		return "Devices fetched."
	}
}

func (m devicesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 1 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.state.DeviceDb.GetDevices()) {
				m.cursor++
			}
		case " ":
			return m, func() tea.Msg {
				light := m.state.DeviceDb.GetDevices()[m.cursor-1]
				light.ChangePowerState(!light.IsPoweredOn())
				return light
			}
		}
	}

	return m, nil
}

func (m devicesModel) View() string {
	var header string
	var content string

	devices := m.state.DeviceDb.GetDevices()
	if len(devices) > 0 {
		header = "Devices discovered:"

		item := 1
		for _, light := range devices {
			cursorText := " "
			if m.cursor == item {
				cursorText = ">"
			}
			powerText := "On"
			if !light.IsPoweredOn() {
				powerText = "Off"
			}
			content = fmt.Sprintf("%s %s %d. %s - %s\n", content, cursorText, item, light.Name(), powerText)
			item++
		}
	} else {
		header = "Fetching devices from hue bridge..."
	}
	footer := "Use arrows to change selection || Press space to toggle power || Press 'CTRL+C' or 'q' to quit"

	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
}

func InitDevicesModel(state *config.ProgramState) devicesModel {
	return devicesModel{
		state:  state,
		cursor: 1,
	}
}
