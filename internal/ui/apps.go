package ui

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// Handle key events when in Apps mode.
func (m model) handleAppsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "enter":
		if len(m.apps) == 0 {
			return m, nil, true
		}
		a := m.apps[m.appsCursor]
		return m, m.enterRevsFor(a), true
	case "R":
		if len(m.apps) == 0 {
			return m, nil, true
		}
		return m, LoadRevsCmd(m.apps[m.appsCursor]), true
	case "l":
		a := m.apps[m.appsCursor]
		cmd := exec.Command("az", "containerapp", "logs", "show", "-n", a.Name, "-g", a.ResourceGroup, "--follow")
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		fmt.Println("--- Ctrl+C to stop logs ---")
		return m, tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} }), true
	case "s", "e":
		a := m.apps[m.appsCursor]
		cmd := exec.Command("az", "containerapp", "exec", "-n", a.Name, "-g", a.ResourceGroup, "--command", "/bin/sh")
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return m, tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} }), true
	case "tab":
		if m.activePane == paneDetails {
			m.activePane = paneRevisions
		} else {
			m.activePane = paneDetails
		}
		return m, nil, true
	}
	return m, nil, false
}

// Update list/spinner when in Apps mode.
// Returns updated model and an aggregated command (if any).
func (m model) updateAppsLists(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd

	var lcmd tea.Cmd
	m.list, lcmd = m.list.Update(msg) // ← arrows/j/k land here
	if lcmd != nil {
		cmds = append(cmds, lcmd)
	}

	// keep your spinner update...
	m.syncAppsCursorFromList()

	// fire loads when selection actually changed
	if m.appsCursor >= 0 && m.appsCursor < len(m.apps) && m.appsCursor != m.lastAppsIndex {
		m.lastAppsIndex = m.appsCursor
		m.jsonView.SetContent(m.headerForCurrent() + "\n\nLoading details...")
		m.revTable.SetRows(nil)
		cmds = append(cmds, tea.Batch(
			LoadDetailsCmd(m.apps[m.appsCursor]),
			LoadRevsCmd(m.apps[m.appsCursor]),
		))
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}
