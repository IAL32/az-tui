package ui

import (
	"fmt"
	"os"
	"os/exec"

	azure "az-tui/internal/azure"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------- Update ----------------------------

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.activePane == paneDetails {
				m.activePane = paneRevisions
			} else {
				m.activePane = paneDetails
			}
			return m, nil
		case "r":
			m.loading, m.err = true, nil
			return m, tea.Batch(LoadAppsCmd(m.rg), m.spin.Tick)
		case "R":
			if len(m.apps) == 0 {
				return m, nil
			}
			return m, LoadRevsCmd(m.apps[m.selected])
		case "l":
			if len(m.apps) == 0 {
				return m, nil
			}
			a := m.apps[m.selected]
			cmd := exec.Command("az", "containerapp", "logs", "show", "-n", a.Name, "-g", a.ResourceGroup, "--follow")
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			fmt.Println("--- Ctrl+C to stop logs ---")
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return noop{} })
		case "e":
			if len(m.apps) == 0 {
				return m, nil
			}
			a := m.apps[m.selected]
			fmt.Println("--- Exec shell; exit to return ---")
			_ = azure.RunAzInteractive("containerapp", "exec", "-n", a.Name, "-g", a.ResourceGroup, "--command", "/bin/sh")
			return m, tea.Suspend
		}

	case tea.WindowSizeMsg:
		// resize components
		w, h := msg.Width, msg.Height
		m.list.SetSize(32, h-2)
		m.jsonView.Width = max(20, w-34)
		m.jsonView.Height = (h - 4) / 2
		m.revTable.SetHeight(h - 4 - m.jsonView.Height)
		return m, nil

	case loadedAppsMsg:
		m.loading = false
		m.err = msg.err
		if msg.err != nil {
			return m, nil
		}
		m.apps = msg.apps
		items := make([]list.Item, len(m.apps))
		for i, a := range m.apps {
			items[i] = item(a)
		}
		m.list.SetItems(items)
		if len(m.apps) > 0 {
			m.selected = 0
			a := m.apps[0]
			return m, tea.Batch(LoadDetailsCmd(a), LoadRevsCmd(a))
		}
		m.jsonView.SetContent("No container apps found.")
		m.revs = nil
		return m, nil

	case loadedDetailsMsg:
		m.err = msg.err
		if msg.err != nil {
			m.json = ""
			m.jsonView.SetContent(StyleError.Render(msg.err.Error()))
			return m, nil
		}
		m.json = msg.json
		m.jsonView.SetContent(m.headerForCurrent() + " " + m.json)
		return m, nil

	case loadedRevsMsg:
		m.err = msg.err
		if msg.err != nil {
			m.revs = nil
			m.revTable.SetRows(nil)
			return m, nil
		}
		m.revs = msg.revs
		rows := make([]table.Row, 0, len(m.revs))
		for _, r := range m.revs {
			act := ""
			if r.Active {
				act = "yes"
			} else {
				act = ""
			}
			rows = append(rows, table.Row{r.Name, act, fmt.Sprintf("%d", r.Traffic)})
		}
		m.revTable.SetRows(rows)
		return m, nil
	}

	// delegate to nested components
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	m.spin, _ = m.spin.Update(msg)
	return m, nil
}

func (m model) headerForCurrent() string {
	if len(m.apps) == 0 {
		return ""
	}
	a := m.apps[m.selected]
	fqdn := a.IngressFQDN
	if fqdn == "" {
		fqdn = "-"
	}
	return fmt.Sprintf("Name: %s  |  RG: %s  |  Loc: %s  |  FQDN: %s  |  Latest: %s",
		a.Name, a.ResourceGroup, a.Location, fqdn, a.LatestRevision)
}
