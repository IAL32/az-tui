package ui

import (
	models "az-tui/internal/models"

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
	for _, app := range m.apps {
		builder.UpdateWidthFromString(columnKeyAppName, app.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(m.apps) > 0 {
		rows = make([]table.Row, len(m.apps))
		for i, app := range m.apps {
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

	t := table.New(columns).
		WithRows(rows).
		BorderRounded().
		WithBaseStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38"))).
		WithMaxTotalWidth(m.termW).
		WithHorizontalFreezeColumnCount(1).
		Filtered(true).
		WithFilterInput(m.appsFilterInput).
		Focused(true)

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
			// Sync the filter with the table after applying
			m.appsTable = m.appsTable.WithFilterInput(m.appsFilterInput)
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
		m.appsFilterInput.SetValue("") // Clear any existing value
		m.appsFilterInput.Focus()
		m.appsTable = m.appsTable.WithFilterInput(m.appsFilterInput)
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

	case "esc":
		// Go back to resource groups mode
		m.mode = modeResourceGroups
		m.loading = true
		m.rg = "" // Clear the selected resource group
		m.resourceGroups = nil
		m.resourceGroupsTable = m.createResourceGroupsTable()
		return m, LoadResourceGroupsCmd(), true
	}

	return m, nil, false
}

// View functions
func (m model) viewApps() string {
	if m.loading && len(m.apps) == 0 {
		// Show loading state using generalized layout
		return m.createLoadingLayout("Loading container apps...")
	}

	if m.err != nil && len(m.apps) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(m.err.Error(), "Press r to retry or q to quit.")
	}

	// Show table view using generalized layout
	tableView := m.appsTable.View()
	return m.createTableLayout(tableView)
}
