package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// createEnvVarsTable creates a table component for environment variables
func (m model) createEnvVarsTable() table.Model {
	// Create dynamic column builder
	builder := NewDynamicColumnBuilder().
		AddColumn("name", "Name", 15, true).  // Min width 15, with filter
		AddColumn("value", "Value", 20, true) // Min width 20, with filter

	// Update column widths based on actual content
	if m.currentContainerName != "" {
		for _, ctr := range m.ctrs {
			if ctr.Name == m.currentContainerName {
				for name, value := range ctr.Env {
					builder.UpdateWidthFromString("name", name)
					builder.UpdateWidthFromString("value", value)
				}
				break
			}
		}
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if m.currentContainerName != "" {
		// Find the current container and get its environment variables
		for _, ctr := range m.ctrs {
			if ctr.Name == m.currentContainerName {
				if len(ctr.Env) > 0 {
					rows = make([]table.Row, 0, len(ctr.Env))
					for name, value := range ctr.Env {
						rows = append(rows, table.NewRow(table.RowData{
							"name":  name,
							"value": value,
						}))
					}
				}
				break
			}
		}
	}

	if len(rows) == 0 {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				"name":  "No environment variables",
				"value": "",
			}),
		}
	}

	return m.createUnifiedTable(columns, rows, m.envVarsFilterInput)
}

// getEnvVarsHelpKeys returns the key bindings for environment variables mode
func (m model) getEnvVarsHelpKeys() keyMap {
	// Environment variables mode is read-only - no action keys
	return keyMap{
		Filter:      m.keys.Filter,
		ScrollLeft:  m.keys.ScrollLeft,
		ScrollRight: m.keys.ScrollRight,
		Help:        m.keys.Help,
		Back:        m.keys.Back,
		Quit:        m.keys.Quit,
	}
}

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
			// Sync the filter with the table after applying
			m.envVarsTable = m.envVarsTable.WithFilterInput(m.envVarsFilterInput)
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
		m.envVarsFilterInput.SetValue("") // Clear any existing value
		m.envVarsFilterInput.Focus()
		m.envVarsTable = m.envVarsTable.WithFilterInput(m.envVarsFilterInput)
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
