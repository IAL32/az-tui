package ui

import (
	models "az-tui/internal/models"
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

	case loadedContainersMsg:
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
	case revisionRestartedMsg:
		if msg.err != nil {
			m.statusLine = fmt.Sprintf("Restart failed: %v", msg.err)
			return m, nil
		}
		m.statusLine = "Revision restart triggered."
		// Optional: refresh revs/containers after a short delay or immediately
		if a, ok := m.currentApp(); ok && appID(a) == msg.appID && m.revName == msg.revName {
			// you can choose to reload containers/revisions; often not needed immediately
			// return m, LoadRevsCmd(a) // if you want to reflect status changes
		}
		return m, nil
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
