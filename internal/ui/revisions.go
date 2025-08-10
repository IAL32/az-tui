package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Handle key events when in Revisions mode.
func (m model) handleRevsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "enter":
		it := m.revList.SelectedItem()
		if it == nil {
			return m, nil, true
		}
		ri := it.(models.RevItem)
		a, ok := m.currentApp()
		if !ok {
			return m, nil, true
		}

		m.mode = modeContainers
		m.revName = ri.Name
		// title + size
		m.ctrList.Title = fmt.Sprintf("Containers — %s@%s", a.Name, ri.Name)
		m.ctrList.SetSize(m.list.Width(), m.list.Height())

		// no cache: clear right pane and load
		m.ctrs = nil
		m.jsonView.SetContent(m.containerHeader(a, ri.Name) + "\n\nLoading containers…")
		return m, LoadContainersCmd(a, ri.Name), true
	case "r":
		it := m.revList.SelectedItem()
		if it == nil {
			// If the list hasn't populated yet, nothing to do.
			// (Optionally: set a statusLine like "Revisions still loading…")
			return m, nil, true
		}
		ri := it.(models.RevItem)
		a, ok := m.currentApp()
		if !ok {
			return m, nil, true
		}

		containerNames := make([]string, 0, len(m.ctrs))
		for _, c := range m.ctrs {
			containerNames = append(containerNames, c.Name)
		}

		txt := fmt.Sprintf("Restart revision?\n\nApp: %s\nRevision: %s\n(affects all containers incl. %q)",
			a.Name, ri.Name, containerNames)

		m = m.withConfirm(
			txt,
			func(mm model) (model, tea.Cmd) {
				mm.statusLine = "Restarting revision..."
				return mm, RestartRevisionCmd(a, ri.Name)
			},
			nil, // no action on cancel
		)
		return m, nil, true
	case "s":
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
func (m model) containerHeader(a models.ContainerApp, rev string) string {
	return fmt.Sprintf("App: %s  |  RG: %s  |  Rev: %s", a.Name, a.ResourceGroup, rev)
}
func (m model) prettyContainerJSON(c models.Container) string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

func (m model) viewRevs() string {
	if m.err != nil && !m.loading {
		return StyleError.Render("Error: ") + m.err.Error() + "  [b/esc] back"
	}

	left := m.revList.View()
	right := styleTitle.Render("Details") + "\n" + m.jsonView.View()
	help := styleAccent.Render("[enter] containers  [e] exec  [l] logs  [r] restart revision  [b/esc] back  [q] quit  (/ filter)")

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
