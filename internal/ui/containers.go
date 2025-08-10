package ui

import (
	"fmt"
	"os"
	"os/exec"

	models "az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

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
