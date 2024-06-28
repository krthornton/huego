package menu

import (
	"fmt"
	"huego/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

type mainModel struct {
	state       *config.ProgramState
	currentMenu tea.Model
}

func (m mainModel) Init() tea.Cmd {
	return m.currentMenu.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {

	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "q":
			// handle quitting at the highest model level
			return m, tea.Quit
		}

	case tea.Model:
		// redirect to specified model/menu
		m.currentMenu = msg.(tea.Model)
		return m, m.currentMenu.Init()
	}

	// propogate updates to child menu
	newChild, childCmd := m.currentMenu.Update(msg)
	m.currentMenu = newChild
	return m, childCmd
}

func (m mainModel) View() string {
	childContent := m.currentMenu.View()
	content := fmt.Sprintf("huego\n\n%s\n", childContent)
	return content
}

func InitMainModel(state *config.ProgramState) mainModel {
	return mainModel{
		state:       state,
		currentMenu: InitDiscoveryModel(state),
	}
}
