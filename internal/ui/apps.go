package ui

import (
	models "az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
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
	// Handle filter input when focused
	if m.appsFilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.appsFilterInput.Blur()
			return m, nil, true
		case "esc":
			m.appsFilterInput.SetValue("")
			m.appsFilterInput.Blur()
			m.appsTable = m.appsTable.WithFilterInput(m.appsFilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.appsFilterInput, cmd = m.appsFilterInput.Update(msg)
			m.appsTable = m.appsTable.WithFilterInput(m.appsFilterInput)
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
		m.appsFilterInput.Focus()
		return m, nil, true
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

	case "r":
		// Refresh apps list - clear data and show loading state
		m.loading = true
		m.apps = nil
		m.appsTable = m.createAppsTable()
		return m, LoadAppsCmd(m.rg), true

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
	title := "Container Apps"

	if m.loading && len(m.apps) == 0 {
		// Show loading state using generalized layout
		return m.createLoadingLayout(title, "Loading container apps...")
	}

	if m.err != nil && len(m.apps) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(title, m.err.Error(), "Press r to retry or q to quit.")
	}

	// Show table view using generalized layout
	tableView := m.appsTable.View()
	return m.createTableLayout(title, tableView)
}
