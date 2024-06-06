package menu

import (
	"fmt"
	"huego-cli/internal/config"
	"huego-cli/internal/hue"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type discoveryModel struct {
	state   *config.ProgramState
	spinner spinner.Model
}

func InitDiscoveryModel(state *config.ProgramState) discoveryModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	return discoveryModel{
		state:   state,
		spinner: spin,
	}
}

func (m discoveryModel) initDiscovery() tea.Msg {
	ipAddr := hue.DiscoverIpAddress()

	var apiKey string
	for _, savedHub := range m.state.Config.Hubs {
		if savedHub.IpAddress == ipAddr {
			apiKey = savedHub.ApiKey
			break
		}
	}

	if apiKey == "" {
		// TODO: implement method for authenticating with new hub
		panic("API key not found for discovered hub")
	}

	m.state.Conn.SetApiKey(apiKey)
	m.state.Conn.SetIpAddress(ipAddr)

	return "IP address discovered"
}

func (m discoveryModel) Init() tea.Cmd {
	return tea.Batch(m.initDiscovery, m.spinner.Tick)
}

func (m discoveryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ipAddr := m.state.Conn.GetIpAddress(); ipAddr != "" {
		return m, func() tea.Msg {
			return InitDevicesModel(m.state)
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m discoveryModel) View() string {
	content := m.spinner.View()
	content = fmt.Sprintf("%s Discovering local hue bridge...", content)
	return content
}
