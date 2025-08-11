package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------- Update ----------------------------

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.confirm.Visible {
			switch msg.String() {
			case "y", "enter":
				m.confirm.Visible = false
				if m.confirm.OnYes != nil {
					return m.confirm.OnYes(m)
				}
				return m, nil
			case "n", "esc":
				m.confirm.Visible = false
				if m.confirm.OnNo != nil {
					return m.confirm.OnNo(m)
				}
				return m, nil
			}
			return m, nil // swallow all other keys when modal visible
		}
		switch m.mode {
		case modeApps:
			if nm, cmd, handled := m.handleAppsKey(msg); handled {
				return nm, cmd
			}
		case modeRevs:
			if nm, cmd, handled := m.handleRevsKey(msg); handled {
				return nm, cmd
			}
		case modeContainers:
			if nm, cmd, handled := m.handleContainersKey(msg); handled {
				return nm, cmd
			}
		case modeEnvVars:
			if nm, cmd, handled := m.handleEnvVarsKey(msg); handled {
				return nm, cmd
			}
		case modeResourceGroups:
			if nm, cmd, handled := m.handleResourceGroupsKey(msg); handled {
				return nm, cmd
			}
		}

		// Handle table navigation for unhandled keys
		var cmd tea.Cmd
		switch m.mode {
		case modeApps:
			m.appsTable, cmd = m.appsTable.Update(msg)
		case modeRevs:
			m.revisionsTable, cmd = m.revisionsTable.Update(msg)
		case modeContainers:
			m.containersTable, cmd = m.containersTable.Update(msg)
		case modeEnvVars:
			m.envVarsTable, cmd = m.envVarsTable.Update(msg)
		case modeResourceGroups:
			m.resourceGroupsTable, cmd = m.resourceGroupsTable.Update(msg)
		}

		if cmd != nil {
			return m, cmd
		}

	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		m.termW, m.termH = w, h
		m.help.Width = w
		return m, nil

	case loadedAppsMsg:
		return m.handleLoadedAppsMsg(msg)

	case loadedRevsMsg:
		return m.handleLoadedRevsMsg(msg)

	case loadedContainersMsg:
		return m.handleLoadedContainersMsg(msg)

	case revisionRestartedMsg:
		return m.handleRevisionRestartedMsg(msg)

	case loadedResourceGroupsMsg:
		return m.handleLoadedResourceGroupsMsg(msg)

	default:
		// Handle spinner updates
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}
