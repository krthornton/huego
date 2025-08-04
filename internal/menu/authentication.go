package menu

import (
	"fmt"
	"huego/internal/config"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type authModel struct {
	state         *config.ProgramState
	spinner       spinner.Model
	authenticated bool
}

func InitAuthenticationModel(state *config.ProgramState) authModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	authState := false
	for _, savedHub := range state.Config.Hubs {
		if savedHub.IpAddress == state.Conn.GetIpAddress() {
			state.Conn.SetApiKey(savedHub.ApiKey)
			authState = true
		}
	}

	return authModel{
		state:         state,
		spinner:       spin,
		authenticated: authState,
	}
}

type authTickMsg struct {
	success bool
}

func (m authModel) authTick() tea.Cmd {
	if m.authenticated {
		return func() tea.Msg {
			return authTickMsg{
				success: true,
			}
		}
	}

	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return authTickMsg{
			success: m.state.Conn.Authenticate(),
		}
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
				m.state.Config.SetApiKeyForIpAddr(m.state.Conn.GetIpAddress(), m.state.Conn.GetApiKey())
				return InitDevicesModel(m.state)
			}
		} else {
			cmd = m.authTick()
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m authModel) View() string {
	var content string = ""
	if !m.authenticated {
		content = fmt.Sprintf("%s Please press the button on your Hue bridge to authenticate...", m.spinner.View())
	}
	return content
}
