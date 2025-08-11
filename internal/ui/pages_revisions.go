package ui

import (
	"fmt"

	models "github.com/IAL32/az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// createRevisionsTable creates a table component for revisions
func (m model) createRevisionsTable() table.Model {
	// Create dynamic column builder
	builder := NewDynamicColumnBuilder().
		AddColumn(columnKeyRevName, "Revision", 15, true).        // Dynamic width, min 15
		AddColumn(columnKeyRevActive, "Active", 8, false).        // Fixed width
		AddColumn(columnKeyRevTraffic, "Traffic", 10, false).     // Fixed width
		AddColumn(columnKeyRevReplicas, "Replicas", 10, false).   // Fixed width
		AddColumn(columnKeyRevScaling, "Scaling", 12, false).     // Fixed width
		AddColumn(columnKeyRevResources, "Resources", 15, false). // Fixed width
		AddColumn(columnKeyRevHealth, "Health", 12, true).        // Fixed width
		AddColumn(columnKeyRevRunning, "Running", 15, true).      // Fixed width
		AddColumn(columnKeyRevCreated, "Created", 20, false).     // Fixed width
		AddColumn(columnKeyRevStatus, "Status", 15, true).        // Fixed width
		AddColumn(columnKeyRevFQDN, "FQDN", 60, false)            // Fixed width (longest content)

	// Update dynamic column widths based on actual content
	for _, rev := range m.revs {
		builder.UpdateWidthFromString(columnKeyRevName, rev.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(m.revs) > 0 {
		rows = make([]table.Row, len(m.revs))
		for i, rev := range m.revs {
			activeMark := "·"
			if rev.Active {
				activeMark = "✓"
			}

			created := "-"
			if !rev.CreatedAt.IsZero() {
				created = rev.CreatedAt.Format("2006-01-02 15:04")
			}

			// Status priority: HealthState > RunningState > ProvisioningState
			status := rev.HealthState
			if status == "" {
				status = rev.RunningState
			}
			if status == "" {
				status = rev.ProvisioningState
			}
			if status == "" {
				status = "-"
			}

			// Current replicas
			replicas := fmt.Sprintf("%d", rev.Replicas)

			// Scaling range
			scaling := fmt.Sprintf("%d-%d", rev.MinReplicas, rev.MaxReplicas)
			if rev.MinReplicas == 0 && rev.MaxReplicas == 0 {
				scaling = "-"
			}

			// Resources
			resources := fmt.Sprintf("%.2gC/%.1s", rev.CPU, rev.Memory)
			if rev.CPU == 0 {
				resources = "-"
			}

			// Health state
			health := rev.HealthState
			if health == "" {
				health = "-"
			}

			// Running state
			running := rev.RunningState
			if running == "" {
				running = "-"
			}

			// FQDN
			fqdn := rev.FQDN
			if fqdn == "" {
				fqdn = "-"
			}

			rows[i] = table.NewRow(table.RowData{
				columnKeyRevName:      rev.Name,
				columnKeyRevActive:    activeMark,
				columnKeyRevTraffic:   fmt.Sprintf("%d%%", rev.Traffic),
				columnKeyRevReplicas:  replicas,
				columnKeyRevScaling:   scaling,
				columnKeyRevResources: resources,
				columnKeyRevHealth:    table.NewStyledCell(health, lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor(health)))),
				columnKeyRevRunning:   table.NewStyledCell(running, lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor(running)))),
				columnKeyRevCreated:   created,
				columnKeyRevStatus:    table.NewStyledCell(status, lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor(status)))),
				columnKeyRevFQDN:      fqdn,
			})
		}
	} else {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				columnKeyRevName:      "No revisions",
				columnKeyRevActive:    "",
				columnKeyRevTraffic:   "",
				columnKeyRevReplicas:  "",
				columnKeyRevScaling:   "",
				columnKeyRevResources: "",
				columnKeyRevHealth:    "",
				columnKeyRevRunning:   "",
				columnKeyRevCreated:   "",
				columnKeyRevStatus:    "",
				columnKeyRevFQDN:      "",
			}),
		}
	}

	t := table.New(columns).
		WithRows(rows).
		BorderRounded().
		WithBaseStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38"))).
		WithMaxTotalWidth(m.termW).
		WithHorizontalFreezeColumnCount(1).
		Filtered(true).
		WithFilterInput(m.revisionsFilterInput).
		Focused(true).
		WithFilterFunc(NewFuzzyFilter(columns))

	// Only sort if we have actual data
	if len(m.revs) > 0 {
		t = t.SortByDesc(columnKeyRevTraffic)
	}

	// Calculate height dynamically based on actual help and status bar heights
	helpBar := m.createHelpBar()
	statusBar := m.createStatusBar()
	helpBarHeight := lipgloss.Height(helpBar)
	statusBarHeight := lipgloss.Height(statusBar)

	// Available height = total height - help bar - status bar
	availableHeight := m.termH - helpBarHeight - statusBarHeight
	if availableHeight > 0 {
		t = t.WithPageSize(availableHeight)
	}

	return t
}

// getRevisionsHelpKeys returns the key bindings for revisions mode
func (m model) getRevisionsHelpKeys() keyMap {
	return keyMap{
		Enter:       m.keys.Enter,
		Refresh:     m.keys.Refresh,
		RestartRev:  m.keys.RestartRev,
		Filter:      m.keys.Filter,
		ScrollLeft:  m.keys.ScrollLeft,
		ScrollRight: m.keys.ScrollRight,
		Help:        m.keys.Help,
		Back:        m.keys.Back,
		Quit:        m.keys.Quit,
	}
}

// Navigation functions
func (m *model) enterRevsFor(a models.ContainerApp) tea.Cmd {
	m.mode = modeRevs
	m.currentAppID = appID(a)

	return LoadRevsCmd(m.dataProvider, a)
}

func (m *model) leaveRevs() {
	m.mode = modeApps
	m.currentAppID = ""

	// Clear revisions state
	m.revs = nil
	m.revisionsTable = m.createRevisionsTable()
}

// Message handlers
func (m model) handleLoadedRevsMsg(msg loadedRevsMsg) (model, tea.Cmd) {
	m.err = msg.err

	if msg.err != nil {
		m.revs = nil
		return m, nil
	}

	m.revs = msg.revs
	// Update the revisions table with new data
	if len(m.revs) > 0 {
		m.revisionsTable = m.createRevisionsTable()
	}

	return m, nil
}

func (m model) handleRevisionRestartedMsg(msg revisionRestartedMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.statusLine = fmt.Sprintf("Restart failed: %v", msg.err)
		return m, nil
	}
	m.statusLine = "Revision restart triggered."
	// Optional: refresh revs/containers after a short delay or immediately
	if a, ok := m.currentApp(); ok && appID(a) == msg.appID && m.currentRevName == msg.revName {
		// Reload revisions to reflect status changes after restart
		return m, LoadRevsCmd(m.dataProvider, a)
	}
	return m, nil
}

// Key handlers
func (m model) handleRevsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	// Handle filter input when focused
	if m.revisionsFilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.revisionsFilterInput.Blur()
			// Sync the filter with the table after applying
			m.revisionsTable = m.revisionsTable.WithFilterInput(m.revisionsFilterInput)
			return m, nil, true
		case "esc":
			m.revisionsFilterInput.SetValue("")
			m.revisionsFilterInput.Blur()
			m.revisionsTable = m.revisionsTable.WithFilterInput(m.revisionsFilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.revisionsFilterInput, cmd = m.revisionsFilterInput.Update(msg)
			m.revisionsTable = m.revisionsTable.WithFilterInput(m.revisionsFilterInput)
			return m, cmd, true
		}
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "?":
		m.help.ShowAll = !m.help.ShowAll
		return m, nil, true
	case "/":
		m.revisionsFilterInput.SetValue("") // Clear any existing value
		m.revisionsFilterInput.Focus()
		m.revisionsTable = m.revisionsTable.WithFilterInput(m.revisionsFilterInput)
		return m, nil, true
	case "enter":
		if len(m.revs) == 0 {
			return m, nil, true
		}

		// Get selected revision from table
		selectedRow := m.revisionsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		revName, ok := selectedRow.Data[columnKeyRevName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the revision by name
		var selectedRev models.Revision
		found := false
		for _, rev := range m.revs {
			if rev.Name == revName {
				selectedRev = rev
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}

		m.mode = modeContainers
		m.currentRevName = selectedRev.Name

		// Clear containers and load new ones
		m.ctrs = nil
		return m, LoadContainersCmd(m.dataProvider, a, selectedRev.Name), true

	case "r":
		// Refresh revisions list - clear data and show loading state
		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}
		m.revs = nil
		m.revisionsTable = m.createRevisionsTable()
		return m, LoadRevsCmd(m.dataProvider, a), true

	case "R":
		if len(m.revs) == 0 {
			return m, nil, true
		}

		// Get selected revision from table
		selectedRow := m.revisionsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		revName, ok := selectedRow.Data[columnKeyRevName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the revision by name
		var selectedRev models.Revision
		found := false
		for _, rev := range m.revs {
			if rev.Name == revName {
				selectedRev = rev
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}

		containerNames := make([]string, 0, len(m.ctrs))
		for _, c := range m.ctrs {
			containerNames = append(containerNames, c.Name)
		}

		txt := fmt.Sprintf("Restart revision?\n\nApp: %s\nRevision: %s\n(affects all containers incl. %q)",
			a.Name, selectedRev.Name, containerNames)

		m = m.withConfirm(
			txt,
			func(mm model) (model, tea.Cmd) {
				mm.statusLine = "Restarting revision..."
				return mm, mm.commandProvider.RestartRevision(a, selectedRev.Name)
			},
			nil, // no action on cancel
		)
		return m, nil, true

	case "s":
		if len(m.revs) == 0 {
			return m, nil, true
		}

		// Get selected revision from table
		selectedRow := m.revisionsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		revName, ok := selectedRow.Data[columnKeyRevName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the revision by name
		var selectedRev models.Revision
		found := false
		for _, rev := range m.revs {
			if rev.Name == revName {
				selectedRev = rev
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}

		return m, m.commandProvider.ExecIntoRevision(a, selectedRev.Name), true

	case "l":
		if len(m.revs) == 0 {
			return m, nil, true
		}

		// Get selected revision from table
		selectedRow := m.revisionsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		revName, ok := selectedRow.Data[columnKeyRevName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the revision by name
		var selectedRev models.Revision
		found := false
		for _, rev := range m.revs {
			if rev.Name == revName {
				selectedRev = rev
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}

		return m, m.commandProvider.ShowRevisionLogs(a, selectedRev.Name), true

	case "esc":
		m.leaveRevs()
		return m, nil, true
	}

	return m, nil, false
}

// View functions
func (m model) viewRevs() string {
	if len(m.revs) == 0 && m.err == nil {
		// Show loading state using generalized layout
		return m.createLoadingLayout("Loading revisions...")
	}

	if m.err != nil && len(m.revs) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(m.err.Error(), "[esc] back  [q] quit")
	}

	// Show table view using generalized layout
	tableView := m.revisionsTable.View()
	return m.createTableLayout(tableView)
}
