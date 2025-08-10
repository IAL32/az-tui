package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------- Update ----------------------------

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.isFiltering() {
			break
		}
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

		leftW := 34
		rightW := max(20, w-leftW-2)

		// Always size all three lists so titles don’t disappear when switching
		m.list.SetSize(leftW, h-2)
		m.revList.SetSize(leftW, h-2)
		m.ctrList.SetSize(leftW, h-2)

		m.jsonView.Width = rightW
		if m.mode == modeApps {
			// Split right pane: Details (top) + Revisions table (bottom)
			m.jsonView.Height = (h - 4) / 2
			m.revTable.SetWidth(rightW)
			m.revTable.SetHeight(h - 4 - m.jsonView.Height)
		} else {
			// Revisions/Containers mode: Details uses full right pane height
			m.jsonView.Height = h - 4
			m.revTable.SetWidth(rightW)
			m.revTable.SetHeight(0) // hidden
		}
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
	switch m.mode {
	case modeContainers:
		return m.updateContainersList(msg)
	case modeRevs:
		return m.updateRevsLists(msg)
	default:
		return m.updateAppsLists(msg)
	}
}

func (m model) isFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

func (m model) headerForCurrent() string {
	if len(m.apps) == 0 || m.appsCursor < 0 || m.appsCursor >= len(m.apps) {
		return ""
	}
	curr := m.apps[m.appsCursor]
	fqdn := curr.IngressFQDN
	if fqdn == "" {
		fqdn = "-"
	}
	return fmt.Sprintf("Name: %s  |  RG: %s  |  Loc: %s  |  FQDN: %s  |  Latest: %s",
		curr.Name, curr.ResourceGroup, curr.Location, fqdn, curr.LatestRevision)
}
