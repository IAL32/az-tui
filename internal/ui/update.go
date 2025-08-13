package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------- Update ----------------------------

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		m.termW, m.termH = w, h
		m.help.Width = w

		// Update context list size
		listWidth := min(70, w-10)
		listHeight := min(12, h-6)
		m.contextList.SetSize(listWidth, listHeight)

		return m, nil

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

		// Handle context list when visible
		if m.showContextList {
			return m.handleContextListKey(msg)
		}

		// Intercept ":" to show context list
		if msg.String() == ":" {
			// Don't show if in filter mode or confirm dialog
			if m.isAnyFilterActive() || m.confirm.Visible {
				return m, nil
			}
			m.showContextList = true
			m.statusLine = "" // Clear status line when entering context mode
			// Recreate context list based on current mode
			m.contextList = m.createContextList()
			return m, nil
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
			m.appsPage.Table, cmd = m.appsPage.Table.Update(msg)
		case modeRevs:
			m.revisionsPage.Table, cmd = m.revisionsPage.Table.Update(msg)
		case modeContainers:
			m.containersPage.Table, cmd = m.containersPage.Table.Update(msg)
		case modeEnvVars:
			m.envVarsPage.Table, cmd = m.envVarsPage.Table.Update(msg)
		case modeResourceGroups:
			m.resourceGroupsPage.Table, cmd = m.resourceGroupsPage.Table.Update(msg)
		}

		if cmd != nil {
			return m, cmd
		}

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
