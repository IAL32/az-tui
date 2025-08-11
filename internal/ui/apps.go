package ui

import (
	models "az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Table column keys for Apps mode
const (
	columnKeyAppName      = "name"
	columnKeyAppRG        = "rg"
	columnKeyAppLocation  = "location"
	columnKeyAppRevision  = "revision"
	columnKeyAppFQDN      = "fqdn"
	columnKeyAppStatus    = "status"
	columnKeyAppReplicas  = "replicas"
	columnKeyAppResources = "resources"
	columnKeyAppIngress   = "ingress"
	columnKeyAppIdentity  = "identity"
	columnKeyAppWorkload  = "workload"
)

// App-related messages
type loadedAppsMsg struct {
	apps []models.ContainerApp
	err  error
}

// Message handlers
func (m model) handleLoadedAppsMsg(msg loadedAppsMsg) (model, tea.Cmd) {
	m.loading = false
	m.err = msg.err
	if msg.err != nil {
		return m, nil
	}
	m.apps = msg.apps

	if len(m.apps) == 0 {
		return m, nil
	}

	// Create the apps table with the loaded data
	m.appsTable = m.createAppsTable()

	return m, nil
}

// Key handlers
func (m model) handleAppsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "enter":
		if len(m.apps) == 0 {
			return m, nil, true
		}

		// Get selected app from table
		selectedRow := m.appsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		appName, ok := selectedRow.Data[columnKeyAppName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the app by name
		var selectedApp models.ContainerApp
		found := false
		for _, app := range m.apps {
			if app.Name == appName {
				selectedApp = app
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		return m, m.enterRevsFor(selectedApp), true

	case "R":
		if len(m.apps) == 0 {
			return m, nil, true
		}

		// Get selected app from table
		selectedRow := m.appsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		appName, ok := selectedRow.Data[columnKeyAppName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the app by name
		var selectedApp models.ContainerApp
		found := false
		for _, app := range m.apps {
			if app.Name == appName {
				selectedApp = app
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		return m, LoadRevsCmd(selectedApp), true

	case "l":
		if len(m.apps) == 0 {
			return m, nil, true
		}

		// Get selected app from table
		selectedRow := m.appsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		appName, ok := selectedRow.Data[columnKeyAppName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the app by name
		var selectedApp models.ContainerApp
		found := false
		for _, app := range m.apps {
			if app.Name == appName {
				selectedApp = app
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		return m, m.azureCommands.ShowAppLogs(selectedApp), true

	case "s", "e":
		if len(m.apps) == 0 {
			return m, nil, true
		}

		// Get selected app from table
		selectedRow := m.appsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		appName, ok := selectedRow.Data[columnKeyAppName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the app by name
		var selectedApp models.ContainerApp
		found := false
		for _, app := range m.apps {
			if app.Name == appName {
				selectedApp = app
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		return m, m.azureCommands.ExecIntoApp(selectedApp), true
	}

	return m, nil, false
}

// View functions
func (m model) viewApps() string {
	if m.loading {
		spinner := m.createSpinner()
		return styleTitle.Render("Loading apps… ") + spinner.View()
	}
	if m.err != nil {
		return StyleError.Render("Error: ") + m.err.Error() + " Press r to retry or q to quit."
	}

	// Show table view
	tableView := m.appsTable.View()
	help := styleAccent.Render("[enter] revisions  [l] logs  [s] exec  [r] refresh  [R] reload revs  [shift+←/→] scroll  [q] quit")

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		styleTitle.Render("Container Apps"),
		tableView,
		help,
		m.statusLine,
	)

	if m.confirm.Visible {
		return lipgloss.Place(m.termW, m.termH, lipgloss.Center, lipgloss.Center, m.confirmBox())
	}
	return body
}
