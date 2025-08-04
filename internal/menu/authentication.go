package menu

import (
	"errors"
	"fmt"
	"huego/internal/config"
	"huego/internal/hue"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type authModel struct {
	state          *config.ProgramState
	spinner        spinner.Model
	buttonRequired bool
}

func InitAuthenticationModel(state *config.ProgramState) authModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	return authModel{
		state:          state,
		spinner:        spin,
		buttonRequired: false,
	}
}

type authTickMsg struct {
	success bool
}

func (m authModel) authTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		err := m.state.Conn.Authenticate(m.state.Config)

		if err != nil {
			var unauthErr hue.UnauthenticatedError
			if errors.As(err, &unauthErr) {
				return authTickMsg{success: false}
			} else {
				panic(err)
			}
		}

		return authTickMsg{success: true}
	})
}

func (m authModel) Init() tea.Cmd {
	return tea.Batch(
		m.authTick(),
		m.spinner.Tick,
	)
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case authTickMsg:
		if msg.success {
			return m, func() tea.Msg {
				return InitDevicesModel(m.state)
			}
		} else {
			m.buttonRequired = true
			cmd = m.authTick()
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m authModel) View() string {
	var content string
	if m.buttonRequired {
		content = fmt.Sprintf("%s Please press the button on your Hue bridge to authenticate...\n", m.spinner.View())
	} else {
		content = fmt.Sprintf("%s Authenticating...\n", m.spinner.View())
	}
	content = fmt.Sprintf("%s\n press 'q' to quit", content)
	return content
}
