package ui

import (
	"fmt"

	models "az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// Table column keys for Revisions mode
const (
	columnKeyRevName      = "name"
	columnKeyRevActive    = "active"
	columnKeyRevTraffic   = "traffic"
	columnKeyRevCreated   = "created"
	columnKeyRevStatus    = "status"
	columnKeyRevReplicas  = "replicas"
	columnKeyRevScaling   = "scaling"
	columnKeyRevResources = "resources"
	columnKeyRevHealth    = "health"
	columnKeyRevRunning   = "running"
	columnKeyRevFQDN      = "fqdn"
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

// Navigation functions
func (m *model) enterRevsFor(a models.ContainerApp) tea.Cmd {
	m.mode = modeRevs
	m.currentAppID = appID(a)

	return LoadRevsCmd(a)
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
		// you can choose to reload containers/revisions; often not needed immediately
		// return m, LoadRevsCmd(a) // if you want to reflect status changes
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
		m.revisionsFilterInput.Focus()
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
		return m, LoadContainersCmd(a, selectedRev.Name), true

	case "r":
		// Refresh revisions list - clear data and show loading state
		a := m.getCurrentApp()
		if a.Name == "" {
			return m, nil, true
		}
		m.revs = nil
		m.revisionsTable = m.createRevisionsTable()
		return m, LoadRevsCmd(a), true

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
				return mm, RestartRevisionCmd(a, selectedRev.Name)
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

		return m, m.azureCommands.ExecIntoRevision(a, selectedRev.Name), true

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

		return m, m.azureCommands.ShowRevisionLogs(a, selectedRev.Name), true

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
