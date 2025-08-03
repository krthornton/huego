package menu

import (
	"fmt"
	"huego/internal/config"
	"huego/internal/hue"
	"net"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type discoveryModel struct {
	state       *config.ProgramState
	spinner     spinner.Model
	hubs        []*net.IP
	cursor      int
	ipCh        chan *net.IP
	discovering bool
}

func InitDiscoveryModel(state *config.ProgramState) discoveryModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	ipCh := make(chan *net.IP, 10)
	hue.DiscoverHueBridges(ipCh)

	return discoveryModel{
		state:       state,
		spinner:     spin,
		hubs:        nil,
		cursor:      1,
		ipCh:        ipCh,
		discovering: true,
	}
}

type discoveryPollMsg struct {
	t time.Time
}

func (m discoveryModel) discoveryPollTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return discoveryPollMsg{
			t: t,
		}
	})
}

type discoveryDoneMsg struct {
	t time.Time
}

func (m discoveryModel) discoveryDoneTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return discoveryDoneMsg{
			t: t,
		}
	})
}

func (m discoveryModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.discoveryPollTick(),
		m.discoveryDoneTick(),
	)
}

func (m discoveryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			if !m.discovering {
				m.discovering = true
				m.hubs = make([]*net.IP, 0)
				m.cursor = 1
				hue.DiscoverHueBridges(m.ipCh)
				cmd = m.discoveryDoneTick()
			}
		case "enter":
			if !m.discovering {
				m, cmd = m.selectHub(m.hubs[m.cursor-1].String())
			}
		case "up":
			if !m.discovering && m.cursor > 1 {
				m.cursor--
			}
		case "down":
			if !m.discovering && m.cursor < len(m.hubs)-1 {
				m.cursor++
			}
		}

	case discoveryPollMsg:
		select {
		case ip := <-m.ipCh:
			m.hubs = append(m.hubs, ip)
		default:
		}
		cmd = m.discoveryPollTick()
	case discoveryDoneMsg:
		m.discovering = false
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m discoveryModel) View() string {
	var content string
	if m.discovering {
		content = fmt.Sprintf("%s Discovering local Hue bridges\n", m.spinner.View())
	} else {
		if len(m.hubs) == 0 {
			content = "No local Hue hubs found. Press 'r' to refresh."
		} else {
			content = "Please select a hub with which to interact to proceed:\n"
			item := 1
			for _, hubIp := range m.hubs {
				cursorText := " "
				if m.cursor == item {
					cursorText = ">"
				}
				content = fmt.Sprintf("%s %s %d. %s\n", content, cursorText, item, hubIp.String())
				item++
			}
			content = fmt.Sprintf("%s\n↑↓ to change selection || ENTER to select hub || 'r' to refresh || 'q' to quit", content)
		}
	}
	return content
}

func (m discoveryModel) selectHub(ipAddr string) (discoveryModel, tea.Cmd) {
	m.state.Conn.SetIpAddress(ipAddr)

	var apiKey string
	for _, savedHub := range m.state.Config.Hubs {
		if savedHub.IpAddress == ipAddr {
			apiKey = savedHub.ApiKey
			break
		}
	}

	if apiKey == "" {
		// we have not authenticated with this hub yet
		return m, func() tea.Msg {
			return InitAuthenticationModel(m.state)
		}
	}

	// we've already authenticated with the discovered hub
	m.state.Conn.SetApiKey(apiKey)

	return m, func() tea.Msg {
		return InitDevicesModel(m.state)
	}
}
