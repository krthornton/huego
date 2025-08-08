package main

import (
	"fmt"
	hue "huego/internal/hue"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type devicesModel struct {
	cursor int
	state  *programState
	lights []*hue.Light
}

type DbInitCompleteEvent struct{}

func (m devicesModel) initDb() tea.Msg {
	m.state.db.Initialize()

	return DbInitCompleteEvent{}
}

type DbUpdateTickEvent struct{}

func dbUpdateTick() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return DbUpdateTickEvent{}
	})
}

func (m devicesModel) Init() tea.Cmd {
	return tea.Batch(
		m.initDb,
		dbUpdateTick(),
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
			if m.cursor < len(m.lights) {
				m.cursor++
			}
		case "left":
			return m, func() tea.Msg {
				light := m.lights[m.cursor-1]
				if !light.GetPowerState() {
					return nil
				}
				currentBrightness := int(light.GetBrightnessLevel())
				desiredBrightness := ((currentBrightness / 10) - 1) * 10
				if desiredBrightness < 0 {
					desiredBrightness = 0
				}
				light.SetBrightnessLevel(desiredBrightness)
				return desiredBrightness
			}
		case "right":
			return m, func() tea.Msg {
				light := m.lights[m.cursor-1]
				if !light.GetPowerState() {
					return nil
				}
				currentBrightness := int(light.GetBrightnessLevel())
				desiredBrightness := ((currentBrightness / 10) + 1) * 10
				if desiredBrightness > 100 {
					desiredBrightness = 100
				}
				light.SetBrightnessLevel(desiredBrightness)
				return desiredBrightness
			}
		case "home":
			return m, func() tea.Msg {
				light := m.lights[m.cursor-1]
				if !light.GetPowerState() {
					return nil
				}
				desiredBrightness := 100
				light.SetBrightnessLevel(desiredBrightness)
				return desiredBrightness
			}
		case "end":
			return m, func() tea.Msg {
				light := m.lights[m.cursor-1]
				if !light.GetPowerState() {
					return nil
				}
				desiredBrightness := 10
				light.SetBrightnessLevel(desiredBrightness)
				return desiredBrightness
			}
		case " ":
			return m, func() tea.Msg {
				light := m.lights[m.cursor-1]
				light.SetPowerState(!light.GetPowerState())
				return light.GetPowerState()
			}
		}
	case DbInitCompleteEvent:
		resources := m.state.db.GetResourcesByType("light")
		m.lights = make([]*hue.Light, 0)
		for _, res := range resources {
			if light, ok := res.(*hue.Light); ok {
				m.lights = append(m.lights, light)
			}
		}

		return m, nil
	case DbUpdateTickEvent:
		m.state.db.Update()
		return m, dbUpdateTick()
	}

	return m, nil
}

func (m devicesModel) View() string {
	var header string
	var content string
	var footer string

	if len(m.lights) > 0 {
		header = "Devices discovered:"

		item := 1
		for _, light := range m.lights {
			cursorText := " "
			if m.cursor == item {
				cursorText = ">"
			}
			powerText := "On"
			if !light.GetPowerState() {
				powerText = "Off"
			}
			content = fmt.Sprintf("%s %s %d. %s - %s", content, cursorText, item, light.Name(), powerText)
			if light.GetPowerState() {
				content = fmt.Sprintf("%s - %d%%", content, int(light.GetBrightnessLevel()))
			}
			content = fmt.Sprintf("%s\n", content)
			item++
		}
		footer = "↑↓ to change selection || space to toggle power || ←→ to change brightness || 'q' to quit"
	} else {
		header = "Fetching devices from hue bridge..."
		footer = "press 'q' to quit"
	}

	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
}

func initDevicesModel(state *programState) devicesModel {
	return devicesModel{
		state:  state,
		cursor: 1,
		lights: nil,
	}
}
