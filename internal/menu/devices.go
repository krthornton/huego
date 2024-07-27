package menu

import (
	"fmt"
	"time"

	"huego/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

type devicesModel struct {
	cursor int
	state  *config.ProgramState
}

type TickMsg time.Time

func (m devicesModel) fetchDevices() tea.Msg {
	m.state.Conn.FetchDevices()

	return nil
}

func (m devicesModel) startEventListener() tea.Msg {
	m.state.Conn.StartEventListener()

	return nil
}

func (m devicesModel) nextDeviceUpdateTick() tea.Cmd {
	// starting checking every second for events to process
	return tea.Every(time.Duration(250*time.Millisecond), func(t time.Time) tea.Msg {
		m.state.Conn.ProcessEvents()

		return TickMsg(t)
	})
}

func (m devicesModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchDevices,
		m.startEventListener,
		m.nextDeviceUpdateTick(),
	)
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
			if m.cursor < len(m.state.Conn.GetDevices()) {
				m.cursor++
			}
		case "left":
			return m, func() tea.Msg {
				light := m.state.Conn.GetDevice(m.cursor - 1)
				if !light.IsPoweredOn() {
					return nil
				}
				currentBrightness := int(light.Brightness())
				desiredBrightness := ((currentBrightness / 10) - 1) * 10
				if desiredBrightness < 0 {
					desiredBrightness = 0
				}
				light.ChangeBrightness(float64(desiredBrightness))
				return desiredBrightness
			}
		case "right":
			return m, func() tea.Msg {
				light := m.state.Conn.GetDevice(m.cursor - 1)
				if !light.IsPoweredOn() {
					return nil
				}
				currentBrightness := int(light.Brightness())
				desiredBrightness := ((currentBrightness / 10) + 1) * 10
				if desiredBrightness > 100 {
					desiredBrightness = 100
				}
				light.ChangeBrightness(float64(desiredBrightness))
				return desiredBrightness
			}
		case " ":
			return m, func() tea.Msg {
				light := m.state.Conn.GetDevice(m.cursor - 1)
				light.ChangePowerState(!light.IsPoweredOn())
				return light.IsPoweredOn()
			}
		}
	case TickMsg:
		// continue checking every second
		return m, m.nextDeviceUpdateTick()
	}

	return m, nil
}

func (m devicesModel) View() string {
	var header string
	var content string

	devices := m.state.Conn.GetDevices()
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
			content = fmt.Sprintf("%s %s %d. %s - %s", content, cursorText, item, light.Name(), powerText)
			if light.IsPoweredOn() {
				content = fmt.Sprintf("%s - %d%%", content, int(light.Brightness()))
			}
			content = fmt.Sprintf("%s\n", content)
			item++
		}
	} else {
		header = "Fetching devices from hue bridge..."
	}
	footer := "↑↓ to change selection || space to toggle power || ←→ to change brightness || 'q' to quit"

	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
}

func InitDevicesModel(state *config.ProgramState) devicesModel {
	return devicesModel{
		state:  state,
		cursor: 1,
	}
}
