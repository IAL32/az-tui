package ui

import (
	models "az-tui/internal/models"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Table column keys for Containers mode
const (
	columnKeyCtrName      = "name"
	columnKeyCtrImage     = "image"
	columnKeyCtrCommand   = "command"
	columnKeyCtrArgs      = "args"
	columnKeyCtrResources = "resources"
	columnKeyCtrEnvCount  = "envcount"
	columnKeyCtrProbes    = "probes"
	columnKeyCtrVolumes   = "volumes"
	columnKeyCtrStatus    = "status"
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
	// Handle filter input when focused
	if m.containersFilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.containersFilterInput.Blur()
			return m, nil, true
		case "esc":
			m.containersFilterInput.SetValue("")
			m.containersFilterInput.Blur()
			m.containersTable = m.containersTable.WithFilterInput(m.containersFilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.containersFilterInput, cmd = m.containersFilterInput.Update(msg)
			m.containersTable = m.containersTable.WithFilterInput(m.containersFilterInput)
			return m, cmd, true
		}
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "/":
		m.containersFilterInput.Focus()
		return m, nil, true
	case "r":
		// Refresh containers list - clear data and show loading state
		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}
		m.ctrs = nil
		m.containersTable = m.createContainersTable()
		return m, LoadContainersCmd(a, m.currentRevName), true
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
	title := fmt.Sprintf("Containers — %s@%s", m.getCurrentAppName(), m.currentRevName)

	if len(m.ctrs) == 0 && m.err == nil {
		// Show loading state using generalized layout
		return m.createLoadingLayout(title, "Loading containers...")
	}

	if m.err != nil && len(m.ctrs) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(title, m.err.Error(), "[esc] back  [q] quit")
	}

	// Show table view using generalized layout
	tableView := m.containersTable.View()
	return m.createTableLayout(title, tableView)
}
