package menu

import (
	"fmt"
	"huego/internal/config"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type authModel struct {
	state   *config.ProgramState
	spinner spinner.Model
}

func InitAuthenticationModel(state *config.ProgramState) authModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	return authModel{
		state:   state,
		spinner: spin,
	}
}

func (m authModel) Init() tea.Cmd {
	return tea.Batch(
		m.state.Conn.Authenticate,
		m.spinner.Tick,
	)
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case string:
		if msg == "Success" {
			return m, func() tea.Msg {
				return InitDevicesModel(m.state)
			}
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m authModel) View() string {
	content := m.spinner.View()
	content = fmt.Sprintf("%s Please press the button on your hue bridge to authenticate...", content)
	return content
}
