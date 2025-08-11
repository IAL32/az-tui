package ui

import (
	models "az-tui/internal/models"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Table column keys for Containers mode
const (
	columnKeyCtrName    = "name"
	columnKeyCtrImage   = "image"
	columnKeyCtrCommand = "command"
	columnKeyCtrArgs    = "args"
	columnKeyCtrStatus  = "status"
)

// Navigation functions
func (m *model) leaveContainers() {
	m.mode = modeRevs
	m.currentRevName = ""

	// Clear containers state
	m.ctrs = nil
	m.containersTable = m.createContainersTable()
}

// Container-related messages
type loadedContainersMsg struct {
	appID   string
	revName string
	ctrs    []models.Container
	err     error
}

// Message handlers
func (m model) handleLoadedContainersMsg(msg loadedContainersMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.ctrs = nil
		return m, nil
	}

	// cache
	m.containersByRev[revKey(msg.appID, msg.revName)] = msg.ctrs
	m.ctrs = msg.ctrs

	// Create the containers table with the loaded data
	if len(m.ctrs) > 0 {
		m.containersTable = m.createContainersTable()
	}

	return m, nil
}

// Key handlers
func (m model) handleContainersKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "s":
		if len(m.ctrs) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		containerName, ok := selectedRow.Data[columnKeyCtrName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the container by name
		var selectedContainer models.Container
		found := false
		for _, ctr := range m.ctrs {
			if ctr.Name == containerName {
				selectedContainer = ctr
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}

		return m, m.azureCommands.ExecIntoContainer(a, m.currentRevName, selectedContainer.Name), true

	case "l":
		if len(m.ctrs) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		containerName, ok := selectedRow.Data[columnKeyCtrName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the container by name
		var selectedContainer models.Container
		found := false
		for _, ctr := range m.ctrs {
			if ctr.Name == containerName {
				selectedContainer = ctr
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}

		return m, m.azureCommands.ShowContainerLogs(a, m.currentRevName, selectedContainer.Name), true

	case "esc":
		m.leaveContainers()
		return m, nil, true
	}

	return m, nil, false
}

// View functions
func (m model) viewContainers() string {
	if m.err != nil && !m.loading {
		return StyleError.Render("Error: ") + m.err.Error() + "  [b/esc] back"
	}

	// Show table view
	tableView := m.containersTable.View()
	help := styleAccent.Render("[s] exec  [l] logs  [shift+←/→] scroll  [b/esc] back  [q] quit")

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		styleTitle.Render(fmt.Sprintf("Containers — %s@%s", m.getCurrentAppName(), m.currentRevName)),
		tableView,
		help,
		m.statusLine,
	)

	if m.confirm.Visible {
		return lipgloss.Place(m.termW, m.termH, lipgloss.Center, lipgloss.Center, m.confirmBox())
	}
	return body
}
