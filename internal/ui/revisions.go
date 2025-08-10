package ui

import (
	"fmt"
	"os"
	"os/exec"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Handle key events when in Revisions mode.
func (m model) handleRevsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "enter", "e":
		it := m.revList.SelectedItem()
		if it == nil {
			return m, nil, true
		}
		ri := it.(models.RevItem)
		// derive app: prefer current cursor; fall back to matching revAppID
		a := m.appForRevContext()
		if a.Name == "" {
			return m, nil, true
		}

		cmd := exec.Command("az", "containerapp", "exec",
			"-n", a.Name, "-g", a.ResourceGroup,
			"--revision", ri.Name, "--command", "/bin/sh",
		)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return m, tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} }), true

	case "l":
		it := m.revList.SelectedItem()
		if it == nil {
			return m, nil, true
		}
		ri := it.(models.RevItem)
		a := m.appForRevContext()
		if a.Name == "" {
			return m, nil, true
		}
		cmd := exec.Command("az", "containerapp", "logs", "show",
			"-n", a.Name, "-g", a.ResourceGroup,
			"--revision", ri.Name, "--follow",
		)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		fmt.Println("--- Ctrl+C to stop logs ---")
		return m, tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} }), true

	case "esc":
		m.syncRevsCursorFromList()
		m.leaveRevs()
		return m, nil, true
	}
	return m, nil, false
}

// Update list/spinner when in Revisions mode.
func (m model) updateRevsLists(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd

	var lcmd tea.Cmd
	m.revList, lcmd = m.revList.Update(msg)
	if lcmd != nil {
		cmds = append(cmds, lcmd)
	}
	m.syncRevsCursorFromList()

	var scmd tea.Cmd
	m.spin, scmd = m.spin.Update(msg)
	if scmd != nil {
		cmds = append(cmds, scmd)
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

// Helper: find the app that revisions pane is showing.
func (m model) appForRevContext() models.ContainerApp {
	if len(m.apps) == 0 || m.revAppID == "" {
		return models.ContainerApp{}
	}
	a := m.apps[m.appsCursor]
	if appID(a) == m.revAppID {
		return a
	}
	for _, x := range m.apps {
		if appID(x) == m.revAppID {
			return x
		}
	}
	return models.ContainerApp{}
}

// When revisions load, also seed the left revisionList items if we’re in this context.
func (m *model) seedRevisionListFromRevisions() {
	if m.mode != modeRevs || m.revAppID == "" {
		return
	}
	items := make([]list.Item, 0, len(m.revs))
	for _, r := range m.revs {
		items = append(items, models.RevItem{r})
	}
	m.revList.SetItems(items)
	// Restore the per-app cursor if valid, else select 0.
	sel := m.revCursorByAppID[m.revAppID]
	if sel < 0 || sel >= len(items) {
		sel = 0
	}
	m.revsCursor = sel
	m.revList.Select(sel)
}
