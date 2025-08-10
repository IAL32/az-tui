package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App-related messages
type loadedAppsMsg struct {
	apps []models.ContainerApp
	err  error
}

type loadedDetailsMsg struct {
	json string
	err  error
}

// Handle loadedAppsMsg
func (m model) handleLoadedAppsMsg(msg loadedAppsMsg) (model, tea.Cmd) {
	m.loading = false
	m.err = msg.err
	if msg.err != nil {
		return m, nil
	}
	m.apps = msg.apps
	m.lastAppsIndex = -1 // force detail load on first render
	items := make([]list.Item, len(m.apps))
	for i, a := range m.apps {
		items[i] = item(a)
	}
	m.list.SetItems(items)
	if len(items) == 0 {
		m.jsonView.SetContent("No container apps found.")
		m.revTable.SetRows(nil)
		return m, nil
	}
	// Trigger initial load
	return m, tea.Batch(
		LoadDetailsCmd(m.apps[m.appsCursor]),
		LoadRevsCmd(m.apps[m.appsCursor]),
	)
}

// Handle loadedDetailsMsg
func (m model) handleLoadedDetailsMsg(msg loadedDetailsMsg) (model, tea.Cmd) {
	m.err = msg.err
	if msg.err != nil {
		m.json = ""
		m.jsonView.SetContent(StyleError.Render(msg.err.Error()))
		return m, nil
	}

	// Pretty-print so indentation is stable
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(msg.json), "", "  "); err != nil {
		// fall back to raw if indent fails
		m.json = msg.json
	} else {
		m.json = buf.String()
	}

	// Ensure indentation starts at col 0 on its own line
	// m.jsonView.SetWrap(false) // don't reflow JSON
	m.jsonView.SetContent(m.headerForCurrent() + "\n\n" + m.json)
	return m, nil
}

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
func (m model) viewApps() string {
	if m.loading {
		return styleTitle.Render("Loading apps… ") + m.spin.View()
	}
	if m.err != nil {
		return StyleError.Render("Error: ") + m.err.Error() + " Press r to retry or q to quit."
	}

	left := m.list.View()
	right := lipgloss.JoinVertical(
		lipgloss.Left,
		styleTitle.Render("Details")+"\n"+m.jsonView.View(),
		styleTitle.Render("Revisions")+"\n"+m.revTable.View(),
	)
	help := styleAccent.Render("[enter] revisions  [l] logs  [s] exec  [r] refresh  [R] reload revs  [q] quit  (/ filter)")

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(34).Render(left),
		lipgloss.NewStyle().Padding(0, 1).Render(right),
	) + "\n" + help + "\n" + m.statusLine

	if m.confirm.Visible {
		return lipgloss.Place(m.termW, m.termH, lipgloss.Center, lipgloss.Center, m.confirmBox())
	}
	return body
}
