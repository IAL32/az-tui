package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------- Update ----------------------------

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.isFiltering() {
			break // filtering is active, do not process global key commands
		}
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.apps)-1 {
				m.cursor++
			}
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
			return m, LoadRevsCmd(m.apps[m.cursor])
		case "l":
			a := m.apps[m.cursor]
			cmd := exec.Command("az", "containerapp", "logs", "show", "-n", a.Name, "-g", a.ResourceGroup, "--follow")
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			fmt.Println("--- Ctrl+C to stop logs ---")
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return noop{} })
		case "s":
			a := m.apps[m.cursor]
			cmd := exec.Command("az", "containerapp", "exec", "-n", a.Name, "-g", a.ResourceGroup, "--command", "/bin/sh")
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return noop{} })
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
		m.lastSelectedIndex = -1 // force detail load on first render
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
			LoadDetailsCmd(m.apps[m.cursor]),
			LoadRevsCmd(m.apps[m.cursor]),
		)

	case loadedDetailsMsg:
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

	case loadedRevsMsg:
		m.err = msg.err

		pad := func(cells ...string) table.Row {
			row := make([]string, 5)
			copy(row, cells)
			return row
		}

		if msg.err != nil {
			m.revs = nil
			m.revTable.SetRows([]table.Row{pad("Error", msg.err.Error())})
			return m, nil
		}

		m.revs = msg.revs
		if len(m.revs) == 0 {
			m.revTable.SetRows([]table.Row{pad("No revisions found")})
			return m, nil
		}

		// Optional: sort by traffic desc
		sort.Slice(m.revs, func(i, j int) bool { return m.revs[i].Traffic > m.revs[j].Traffic })

		rows := make([]table.Row, 0, len(m.revs))
		for _, r := range m.revs {
			activeMark := "·"
			if r.Active {
				activeMark = "✓"
			}

			created := "-"
			if !r.CreatedAt.IsZero() {
				created = r.CreatedAt.Format("2006-01-02 15:04")
			}

			status := r.Status
			if status == "" {
				status = "-"
			}

			rows = append(rows, pad(
				r.Name,
				activeMark,
				fmt.Sprintf("%3d%%", r.Traffic), // number only, right-aligned
				created,
				status,
			))
		}

		m.revTable.SetRows(rows)
		return m, nil
	}

	return updateSelectedItems(m, msg)
}

func (m model) isFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

/**
 * updateSelectedItems updates the selected items in the model based on the given message.
 * It delegates to the list component for handling list-specific messages.
 * Finally, it updates the model with any commands returned by the list component.
 */
func updateSelectedItems(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	appendCmd := func(cmd tea.Cmd) {
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Always update the list first
	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	appendCmd(listCmd)

	// Always update spinner
	var spinCmd tea.Cmd
	m.spin, spinCmd = m.spin.Update(msg)
	appendCmd(spinCmd)

	// If filtering, stop here (don't trigger global commands)
	if m.isFiltering() {
		return m, tea.Batch(cmds...)
	}

	// Guard against invalid cursor or empty list
	if len(m.apps) == 0 || m.cursor < 0 || m.cursor >= len(m.apps) {
		return m, tea.Batch(cmds...)
	}

	currIndex := m.cursor
	currApp := m.apps[currIndex]

	// Only load details/revs if the selection has changed
	if currIndex != m.lastSelectedIndex {
		m.lastSelectedIndex = currIndex

		// Show a loading placeholder
		m.jsonView.SetContent(m.headerForCurrent() + "\n\n" + "Loading details...")
		m.revTable.SetRows(nil) // clear rev table until loaded

		appendCmd(LoadDetailsCmd(currApp))
		appendCmd(LoadRevsCmd(currApp))
	}

	return m, tea.Batch(cmds...)
}

func (m model) headerForCurrent() string {
	if len(m.apps) == 0 || m.cursor < 0 || m.cursor >= len(m.apps) {
		return ""
	}
	curr := m.apps[m.cursor]
	fqdn := curr.IngressFQDN
	if fqdn == "" {
		fqdn = "-"
	}
	return fmt.Sprintf("Name: %s  |  RG: %s  |  Loc: %s  |  FQDN: %s  |  Latest: %s",
		curr.Name, curr.ResourceGroup, curr.Location, fqdn, curr.LatestRevision)
}
