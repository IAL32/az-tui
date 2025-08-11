package ui

import (
	models "az-tui/internal/models"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Container-related messages
type loadedContainersMsg struct {
	appID   string
	revName string
	ctrs    []models.Container
	err     error
}

// Handle loadedContainersMsg
func (m model) handleLoadedContainersMsg(msg loadedContainersMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.ctrs = nil
		return m, nil
	}

	// cache
	m.containersByRev[revKey(msg.appID, msg.revName)] = msg.ctrs
	m.ctrs = msg.ctrs

	// Set initial selection
	if len(m.ctrs) > 0 {
		m.selectedContainer = 0
	}

	return m, nil
}

type ctrItem struct{ Container models.Container }

func (ci ctrItem) Title() string       { return ci.Container.Name }
func (ci ctrItem) Description() string { return ci.Container.Image }
func (ci ctrItem) FilterValue() string { return ci.Container.Name + " " + ci.Container.Image }

func (m model) handleContainersKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "j", "down":
		if m.selectedContainer < len(m.ctrs)-1 {
			m.selectedContainer++
		}
		return m, nil, true
	case "k", "up":
		if m.selectedContainer > 0 {
			m.selectedContainer--
		}
		return m, nil, true
	case "s":
		if len(m.ctrs) == 0 || m.selectedContainer >= len(m.ctrs) {
			return m, nil, true
		}
		container := m.ctrs[m.selectedContainer]
		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}

		return m, m.azureCommands.ExecIntoContainer(a, m.currentRevName, container.Name), true

	case "l":
		if len(m.ctrs) == 0 || m.selectedContainer >= len(m.ctrs) {
			return m, nil, true
		}
		container := m.ctrs[m.selectedContainer]
		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}

		return m, m.azureCommands.ShowContainerLogs(a, m.currentRevName, container.Name), true

	case "esc":
		m.mode = modeRevs
		return m, nil, true
	}
	return m, nil, false
}

func (m model) viewContainers() string {
	if m.err != nil && !m.loading {
		return StyleError.Render("Error: ") + m.err.Error() + "  [b/esc] back"
	}

	// Create components from current data
	containersList := m.createContainersList()
	detailsView := m.createDetailsView()

	left := containersList.View()
	right := styleTitle.Render("Details") + "\n" + detailsView.View()
	help := styleAccent.Render("[s] exec  [l] logs  [r] restart  [b/esc] back  [q] quit  (j/k navigate)")

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(34).Render(left),
		lipgloss.NewStyle().Padding(0, 1).Render(right),
	) + "\n" + help + "\n" + m.statusLine

	if m.confirm.Visible {
		return lipgloss.Place(m.termW, m.termH, lipgloss.Center, lipgloss.Center, m.confirmBox())
	}
	return body
}

// Component factory methods for containers mode

// createContainersList creates a list component for containers
func (m model) createContainersList() list.Model {
	items := make([]list.Item, len(m.ctrs))
	for i, ctr := range m.ctrs {
		items[i] = ctrItem{Container: ctr}
	}

	l := list.New(items, list.NewDefaultDelegate(), 32, m.termH-2)
	l.Title = fmt.Sprintf("Containers — %s@%s", m.getCurrentAppName(), m.currentRevName)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	if m.selectedContainer >= 0 && m.selectedContainer < len(items) {
		l.Select(m.selectedContainer)
	}

	return l
}

// containerHeader returns header information for containers mode
func (m model) containerHeader(a models.ContainerApp, rev string) string {
	return fmt.Sprintf("App: %s  |  RG: %s  |  Rev: %s", a.Name, a.ResourceGroup, rev)
}

// prettyContainerJSON returns formatted JSON for a container
func (m model) prettyContainerJSON(c models.Container) string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}
