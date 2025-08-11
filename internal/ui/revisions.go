package ui

import (
	"fmt"
	"sort"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Revision-related messages
type loadedRevsMsg struct {
	revs []models.Revision
	err  error
}

type revisionRestartedMsg struct {
	appID   string
	revName string
	err     error
	out     string
}

// Handle loadedRevsMsg
func (m model) handleLoadedRevsMsg(msg loadedRevsMsg) (model, tea.Cmd) {
	m.err = msg.err

	if msg.err != nil {
		m.revs = nil
		return m, nil
	}

	m.revs = msg.revs
	if len(m.revs) > 0 {
		m.selectedRevision = 0
	}

	return m, nil
}

// Handle revisionRestartedMsg
func (m model) handleRevisionRestartedMsg(msg revisionRestartedMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.statusLine = fmt.Sprintf("Restart failed: %v", msg.err)
		return m, nil
	}
	m.statusLine = "Revision restart triggered."
	// Optional: refresh revs/containers after a short delay or immediately
	if a, ok := m.currentApp(); ok && appID(a) == msg.appID && m.currentRevName == msg.revName {
		// you can choose to reload containers/revisions; often not needed immediately
		// return m, LoadRevsCmd(a) // if you want to reflect status changes
	}
	return m, nil
}

// Handle key events when in Revisions mode.
func (m model) handleRevsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "j", "down":
		if m.selectedRevision < len(m.revs)-1 {
			m.selectedRevision++
		}
		return m, nil, true
	case "k", "up":
		if m.selectedRevision > 0 {
			m.selectedRevision--
		}
		return m, nil, true
	case "enter":
		if len(m.revs) == 0 || m.selectedRevision >= len(m.revs) {
			return m, nil, true
		}
		rev := m.revs[m.selectedRevision]
		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}

		m.mode = modeContainers
		m.currentRevName = rev.Name

		// Clear containers and load new ones
		m.ctrs = nil
		m.selectedContainer = 0
		return m, LoadContainersCmd(a, rev.Name), true
	case "r":
		if len(m.revs) == 0 || m.selectedRevision >= len(m.revs) {
			return m, nil, true
		}
		rev := m.revs[m.selectedRevision]
		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}

		containerNames := make([]string, 0, len(m.ctrs))
		for _, c := range m.ctrs {
			containerNames = append(containerNames, c.Name)
		}

		txt := fmt.Sprintf("Restart revision?\n\nApp: %s\nRevision: %s\n(affects all containers incl. %q)",
			a.Name, rev.Name, containerNames)

		m = m.withConfirm(
			txt,
			func(mm model) (model, tea.Cmd) {
				mm.statusLine = "Restarting revision..."
				return mm, RestartRevisionCmd(a, rev.Name)
			},
			nil, // no action on cancel
		)
		return m, nil, true
	case "s":
		if len(m.revs) == 0 || m.selectedRevision >= len(m.revs) {
			return m, nil, true
		}
		rev := m.revs[m.selectedRevision]
		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}

		return m, m.azureCommands.ExecIntoRevision(a, rev.Name), true

	case "l":
		if len(m.revs) == 0 || m.selectedRevision >= len(m.revs) {
			return m, nil, true
		}
		rev := m.revs[m.selectedRevision]
		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}
		return m, m.azureCommands.ShowRevisionLogs(a, rev.Name), true

	case "esc":
		m.leaveRevs()
		return m, nil, true
	}
	return m, nil, false
}

func (m model) viewRevs() string {
	if m.err != nil && !m.loading {
		return StyleError.Render("Error: ") + m.err.Error() + "  [b/esc] back"
	}

	// Create components from current data
	revsList := m.createRevisionsList()
	detailsView := m.createDetailsView()

	left := revsList.View()
	right := styleTitle.Render("Details") + "\n" + detailsView.View()
	help := styleAccent.Render("[enter] containers  [s] exec  [l] logs  [r] restart revision  [b/esc] back  [q] quit  (j/k navigate)")

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

// Component factory methods for revisions mode

// createRevisionsList creates a list component for revisions
func (m model) createRevisionsList() list.Model {
	items := make([]list.Item, len(m.revs))
	for i, rev := range m.revs {
		items[i] = models.RevItem{Revision: rev}
	}

	l := list.New(items, list.NewDefaultDelegate(), 32, m.termH-2)
	l.Title = fmt.Sprintf("Revisions — %s", m.getCurrentAppName())
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	if m.selectedRevision >= 0 && m.selectedRevision < len(items) {
		l.Select(m.selectedRevision)
	}

	return l
}

// createRevisionsTable creates a table component for revisions
func (m model) createRevisionsTable() table.Model {
	columns := []table.Column{
		{Title: "Revision", Width: 28},
		{Title: "Active", Width: 6},
		{Title: "Traffic", Width: 8},
		{Title: "Created", Width: 18},
		{Title: "Status", Width: 12},
	}

	t := table.New(table.WithColumns(columns))

	if len(m.revs) == 0 {
		t.SetRows([]table.Row{{"No revisions found", "", "", "", ""}})
		return t
	}

	// Sort by traffic desc (same as original logic)
	sortedRevs := make([]models.Revision, len(m.revs))
	copy(sortedRevs, m.revs)
	sort.Slice(sortedRevs, func(i, j int) bool { return sortedRevs[i].Traffic > sortedRevs[j].Traffic })

	rows := make([]table.Row, len(sortedRevs))
	for i, r := range sortedRevs {
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

		rows[i] = table.Row{
			r.Name,
			activeMark,
			fmt.Sprintf("%3d%%", r.Traffic),
			created,
			status,
		}
	}

	t.SetRows(rows)

	// Size the table appropriately
	rightW := max(20, m.termW-34-2)
	t.SetWidth(rightW)
	if m.mode == modeApps {
		t.SetHeight(m.termH - 4 - (m.termH-4)/2) // Bottom half of right pane
	} else {
		t.SetHeight(0) // Hidden in other modes
	}

	return t
}

// getRevisionContext returns context information for revisions mode
func (m model) getRevisionContext() string {
	app := m.getCurrentApp()
	if app.Name == "" {
		return "No app context"
	}
	return fmt.Sprintf("App: %s\nResource Group: %s\nRevisions: %d",
		app.Name, app.ResourceGroup, len(m.revs))
}
