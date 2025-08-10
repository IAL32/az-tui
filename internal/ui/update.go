package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
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
			break
		}
		switch m.mode {
		case modeApps:
			m, cmd, handled := m.handleAppsKey(msg)
			if handled {
				return m, cmd
			}
		case modeRevs:
			m, cmd, handled := m.handleRevsKey(msg)
			if handled {
				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		// resize components
		w, h := msg.Width, msg.Height
		leftW := 34
		if m.mode == modeRevs {
			m.revList.SetSize(leftW, h-2)
		} else {
			m.list.SetSize(leftW, h-2)
		}
		m.jsonView.Width = max(20, w-leftW-2)
		m.jsonView.Height = (h - 4) / 2
		m.revTable.SetWidth(m.jsonView.Width)
		m.revTable.SetHeight(h - 4 - m.jsonView.Height)

	case loadedAppsMsg:
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
		m.seedRevisionListFromRevisions()
		return m, nil
	}

	if m.mode == modeRevs {
		return m.updateRevsLists(msg)
	}
	return m.updateAppsLists(msg)
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
