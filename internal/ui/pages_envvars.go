package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Navigation functions
func (m *model) leaveEnvVars() {
	m.mode = modeContainers
	m.currentContainerName = ""
}

// Key handlers
func (m model) handleEnvVarsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	// Handle filter input when focused
	if m.envVarsFilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.envVarsFilterInput.Blur()
			return m, nil, true
		case "esc":
			m.envVarsFilterInput.SetValue("")
			m.envVarsFilterInput.Blur()
			m.envVarsTable = m.envVarsTable.WithFilterInput(m.envVarsFilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.envVarsFilterInput, cmd = m.envVarsFilterInput.Update(msg)
			m.envVarsTable = m.envVarsTable.WithFilterInput(m.envVarsFilterInput)
			return m, cmd, true
		}
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "?":
		m.help.ShowAll = !m.help.ShowAll
		return m, nil, true
	case "/":
		m.envVarsFilterInput.Focus()
		return m, nil, true
	case "esc":
		m.leaveEnvVars()
		return m, nil, true
	}

	return m, nil, false
}

// View functions
func (m model) viewEnvVars() string {
	// Show table view using generalized layout
	tableView := m.envVarsTable.View()
	return m.createTableLayout(tableView)
}
