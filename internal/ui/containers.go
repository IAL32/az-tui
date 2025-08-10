package ui

import (
	"fmt"
	"os"
	"os/exec"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Container-related messages
type loadedContainersMsg struct {
	appID   string
	revName string
	ctrs    []models.Container
	err     error
}

// Handle loadedContainersMsg
func (m model) handleLoadedContainersMsg(msg loadedContainersMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.ctrs = nil
		m.ctrList.SetItems([]list.Item{ctrItem{Container: models.Container{Name: "Error", Image: msg.err.Error()}}})
		return m, nil
	}

	// cache
	m.containersByRev[revKey(msg.appID, msg.revName)] = msg.ctrs
	m.ctrs = msg.ctrs

	// build left items
	items := make([]list.Item, 0, len(m.ctrs))
	for _, c := range m.ctrs {
		items = append(items, ctrItem{c})
	}
	m.ctrList.SetItems(items)

	// select previous or 0
	sel := m.ctrCursor
	if sel < 0 || sel >= len(items) {
		sel = 0
	}
	m.ctrList.Select(sel)
	m.lastCtrIndex = -1 // force right-pane refresh on first movement

	// render first container details if available
	if len(m.ctrs) > 0 {
		a, ok := m.currentApp()
		if ok {
			m.jsonView.SetContent(m.containerHeader(a, msg.revName) + "\n\n" + m.prettyContainerJSON(m.ctrs[sel]))
		}
	} else {
		a, ok := m.currentApp()
		if ok {
			m.jsonView.SetContent(m.containerHeader(a, msg.revName) + "\n\nNo containers found.")
		}
	}
	return m, nil
}

type ctrItem struct{ Container models.Container }

func (ci ctrItem) Title() string       { return ci.Container.Name }
func (ci ctrItem) Description() string { return ci.Container.Image }
func (ci ctrItem) FilterValue() string { return ci.Container.Name + " " + ci.Container.Image }

func (m model) handleContainersKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "s":
		it := m.ctrList.SelectedItem()
		if it == nil {
			return m, nil, true
		}
		ci := it.(ctrItem)
		a, ok := m.currentApp()
		if !ok || m.revName == "" {
			return m, nil, true
		}

		// exec into specific revision + container
		cmd := exec.Command("az", "containerapp", "exec",
			"-n", a.Name, "-g", a.ResourceGroup,
			"--revision", m.revName,
			"--container", ci.Container.Name,
			"--command", "/bin/sh",
		)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return m, tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} }), true

	case "l":
		it := m.ctrList.SelectedItem()
		if it == nil {
			return m, nil, true
		}
		ci := it.(ctrItem)
		a, ok := m.currentApp()
		if !ok || m.revName == "" {
			return m, nil, true
		}

		cmd := exec.Command("az", "containerapp", "logs", "show",
			"-n", a.Name, "-g", a.ResourceGroup,
			"--revision", m.revName,
			"--container", ci.Container.Name,
			"--follow",
		)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		fmt.Println("--- Ctrl+C to stop logs ---")
		return m, tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} }), true

	case "esc":
		// remember cursor in this revision context
		m.ctrCursor = m.ctrList.Index()
		m.lastCtrIndex = m.ctrCursor
		m.mode = modeRevs
		return m, nil, true
	}
	return m, nil, false
}

func (m model) updateContainersList(msg tea.Msg) (model, tea.Cmd) {
	var cmds []tea.Cmd
	var lcmd tea.Cmd
	m.ctrList, lcmd = m.ctrList.Update(msg)
	if lcmd != nil {
		cmds = append(cmds, lcmd)
	}

	// when selection changes, refresh right pane with container details (from m.ctrs)
	idx := m.ctrList.Index()
	if idx >= 0 && idx < len(m.ctrs) && idx != m.lastCtrIndex {
		m.lastCtrIndex = idx
		a, ok := m.currentApp()
		if ok {
			m.jsonView.SetContent(m.containerHeader(a, m.revName) + "\n\n" + m.prettyContainerJSON(m.ctrs[idx]))
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}
func (m model) viewContainers() string {
	if m.err != nil && !m.loading {
		return StyleError.Render("Error: ") + m.err.Error() + "  [b/esc] back"
	}

	left := m.ctrList.View()
	right := styleTitle.Render("Details") + "\n" + m.jsonView.View()
	help := styleAccent.Render("[enter/e] exec  [l] logs  [r] restart  [b/esc] back  [q] quit  (/ filter)")

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
