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
		}

	case tea.WindowSizeMsg:
		w, h := msg.Width, msg.Height
		m.termW, m.termH = w, h
		return m, nil

	case loadedAppsMsg:
		return m.handleLoadedAppsMsg(msg)

	case loadedDetailsMsg:
		return m.handleLoadedDetailsMsg(msg)

	case loadedRevsMsg:
		return m.handleLoadedRevsMsg(msg)

	case loadedContainersMsg:
		return m.handleLoadedContainersMsg(msg)

	case revisionRestartedMsg:
		return m.handleRevisionRestartedMsg(msg)
	}
	return m, nil
}
