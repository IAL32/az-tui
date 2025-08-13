package ui

import (
	models "github.com/IAL32/az-tui/internal/models"

	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// createAppsTable creates a table component for container apps
func (m model) createAppsTable() table.Model {
	// Create dynamic column builder
	builder := NewDynamicColumnBuilder().
		AddColumn(columnKeyAppName, "Name", 15, true).                 // Dynamic width, min 15
		AddColumn(columnKeyAppLocation, "Location", 15, true).         // Fixed width
		AddColumn(columnKeyAppStatus, "Status", 12, true).             // Fixed width
		AddColumn(columnKeyAppReplicas, "Replicas", 10, false).        // Fixed width
		AddColumn(columnKeyAppResources, "Resources", 12, false).      // Fixed width
		AddColumn(columnKeyAppIngress, "Ingress", 18, false).          // Fixed width
		AddColumn(columnKeyAppIdentity, "Identity", 15, false).        // Fixed width
		AddColumn(columnKeyAppWorkload, "Workload", 15, false).        // Fixed width
		AddColumn(columnKeyAppRevision, "Latest Revision", 30, false). // Fixed width
		AddColumn(columnKeyAppFQDN, "FQDN", 60, false)                 // Fixed width (longest content)

	// Update dynamic column widths based on actual content
	for _, app := range m.appsPage.Data {
		builder.UpdateWidthFromString(columnKeyAppName, app.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(m.appsPage.Data) > 0 {
		rows = make([]table.Row, len(m.appsPage.Data))
		for i, app := range m.appsPage.Data {
			fqdn := app.IngressFQDN
			if fqdn == "" {
				fqdn = "-"
			}

			status := app.RunningStatus
			if status == "" {
				status = app.ProvisioningState
			}
			if status == "" {
				status = "Unknown"
			}

			// Format replicas
			replicas := fmt.Sprintf("%d-%d", app.MinReplicas, app.MaxReplicas)
			if app.MinReplicas == 0 && app.MaxReplicas == 0 {
				replicas = "-"
			}

			// Format resources
			resources := fmt.Sprintf("%.2gC/%.1s", app.CPU, app.Memory)
			if app.CPU == 0 {
				resources = "-"
			}

			// Format ingress
			ingress := "None"
			if app.IngressFQDN != "" {
				if app.IngressExternal {
					ingress = "External"
				} else {
					ingress = "Internal"
				}
				if app.TargetPort > 0 {
					ingress += fmt.Sprintf(":%d", app.TargetPort)
				}
			}

			// Format identity
			identity := app.IdentityType
			if identity == "" {
				identity = "None"
			}

			// Format workload profile
			workload := app.WorkloadProfile
			if workload == "" {
				workload = "Consumption"
			}

			rows[i] = table.NewRow(table.RowData{
				columnKeyAppName:      app.Name,
				columnKeyAppLocation:  app.Location,
				columnKeyAppStatus:    table.NewStyledCell(status, lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor(status)))),
				columnKeyAppReplicas:  replicas,
				columnKeyAppResources: resources,
				columnKeyAppIngress:   ingress,
				columnKeyAppIdentity:  identity,
				columnKeyAppWorkload:  workload,
				columnKeyAppRevision:  app.LatestRevision,
				columnKeyAppFQDN:      fqdn,
			})
		}
	}
	// Don't show any placeholder rows - empty table is fine

	return m.createUnifiedTable(columns, rows, m.appsPage.FilterInput)
}

// getAppsHelpKeys returns the key bindings for apps mode
func (m model) getAppsHelpKeys() keyMap {
	return keyMap{
		Enter:       m.keys.Enter,
		Refresh:     m.keys.Refresh,
		Filter:      m.keys.Filter,
		ScrollLeft:  m.keys.ScrollLeft,
		ScrollRight: m.keys.ScrollRight,
		Help:        m.keys.Help,
		Back:        m.keys.Back,
		Quit:        m.keys.Quit,
	}
}

// Message handlers
func (m model) handleLoadedAppsMsg(msg loadedAppsMsg) (model, tea.Cmd) {
	m.appsPage.IsLoading = false
	m.appsPage.Error = msg.err
	if msg.err != nil {
		m.appsPage.Data = nil
		return m, nil
	}
	m.appsPage.Data = msg.apps

	// Create the apps table with the loaded data
	m.appsPage.Table = m.createAppsTable()

	return m, nil
}

// Key handlers
func (m model) handleAppsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	// Handle filter input when focused
	if m.appsPage.FilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.appsPage.FilterInput.Blur()
			// Sync the filter with the table after applying
			m.appsPage.Table = m.appsPage.Table.WithFilterInput(m.appsPage.FilterInput)
			return m, nil, true
		case "esc":
			m.appsPage.FilterInput.SetValue("")
			m.appsPage.FilterInput.Blur()
			m.appsPage.Table = m.appsPage.Table.WithFilterInput(m.appsPage.FilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.appsPage.FilterInput, cmd = m.appsPage.FilterInput.Update(msg)
			m.appsPage.Table = m.appsPage.Table.WithFilterInput(m.appsPage.FilterInput)
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
		m.appsPage.FilterInput.SetValue("") // Clear any existing value
		m.appsPage.FilterInput.Focus()
		m.appsPage.Table = m.appsPage.Table.WithFilterInput(m.appsPage.FilterInput)
		return m, nil, true
	case "enter":
		if len(m.appsPage.Data) == 0 {
			return m, nil, true
		}

		// Get selected app from table
		selectedRow := m.appsPage.Table.HighlightedRow()
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
		for _, app := range m.appsPage.Data {
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
		m.appsPage.IsLoading = true
		m.appsPage.Error = nil
		m.appsPage.Data = nil
		m.appsPage.Table = m.createAppsTable()
		return m, LoadAppsCmd(m.dataProvider, m.currentRG), true

	case "l":
		if len(m.appsPage.Data) == 0 {
			return m, nil, true
		}

		// Get selected app from table
		selectedRow := m.appsPage.Table.HighlightedRow()
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
		for _, app := range m.appsPage.Data {
			if app.Name == appName {
				selectedApp = app
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		return m, m.commandProvider.ShowAppLogs(selectedApp), true

	case "s", "e":
		if len(m.appsPage.Data) == 0 {
			return m, nil, true
		}

		// Get selected app from table
		selectedRow := m.appsPage.Table.HighlightedRow()
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
		for _, app := range m.appsPage.Data {
			if app.Name == appName {
				selectedApp = app
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		return m, m.commandProvider.ExecIntoApp(selectedApp), true

	case "esc":
		// Go back to resource groups mode
		m.mode = modeResourceGroups
		m.resourceGroupsPage.IsLoading = true
		m.resourceGroupsPage.Error = nil
		m.currentRG = "" // Clear the selected resource group
		m.resourceGroupsPage.Data = nil
		m.resourceGroupsPage.Table = m.createResourceGroupsTable()
		return m, LoadResourceGroupsCmd(m.dataProvider), true
	}

	return m, nil, false
}

// View functions
func (m model) viewApps() string {
	if m.appsPage.IsLoading && len(m.appsPage.Data) == 0 {
		// Show loading state using generalized layout
		return m.createLoadingLayout("Loading container apps...")
	}

	if m.appsPage.Error != nil && len(m.appsPage.Data) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(m.appsPage.Error.Error(), "Press r to retry or q to quit.")
	}

	// Show table view using generalized layout
	tableView := m.appsPage.Table.View()
	return m.createTableLayout(tableView)
}
