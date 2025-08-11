package ui

import (
	models "az-tui/internal/models"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// Spinner component factory
func (m model) createSpinner() spinner.Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	return sp
}

// Table component factories

// createAppsTable creates a table component for container apps
func (m model) createAppsTable() table.Model {
	// Calculate dynamic widths for Name and Resource Group columns based on content
	nameWidth := len("Name")
	rgWidth := len("Resource Group")

	// Find maximum width needed for Name and Resource Group columns
	for _, app := range m.apps {
		if len(app.Name) > nameWidth {
			nameWidth = len(app.Name)
		}
		if len(app.ResourceGroup) > rgWidth {
			rgWidth = len(app.ResourceGroup)
		}
	}

	// Add padding and set minimum widths
	nameWidth += 2
	rgWidth += 2
	if nameWidth < 15 {
		nameWidth = 15
	}
	if rgWidth < 18 {
		rgWidth = 18
	}

	columns := []table.Column{
		table.NewColumn(columnKeyAppName, "Name", nameWidth).WithFiltered(true),       // Dynamic width based on content
		table.NewColumn(columnKeyAppRG, "Resource Group", rgWidth).WithFiltered(true), // Dynamic width based on content
		table.NewColumn(columnKeyAppLocation, "Location", 15).WithFiltered(true),      // Fixed width
		table.NewColumn(columnKeyAppStatus, "Status", 12).WithFiltered(true),          // Fixed width
		table.NewColumn(columnKeyAppReplicas, "Replicas", 10),                         // Fixed width
		table.NewColumn(columnKeyAppResources, "Resources", 12),                       // Fixed width
		table.NewColumn(columnKeyAppIngress, "Ingress", 18),                           // Fixed width
		table.NewColumn(columnKeyAppIdentity, "Identity", 15),                         // Fixed width
		table.NewColumn(columnKeyAppWorkload, "Workload", 15),                         // Fixed width
		table.NewColumn(columnKeyAppRevision, "Latest Revision", 30),                  // Fixed width
		table.NewColumn(columnKeyAppFQDN, "FQDN", 60),                                 // Fixed width (longest content)
	}

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
				columnKeyAppRG:        app.ResourceGroup,
				columnKeyAppLocation:  app.Location,
				columnKeyAppStatus:    status,
				columnKeyAppReplicas:  replicas,
				columnKeyAppResources: resources,
				columnKeyAppIngress:   ingress,
				columnKeyAppIdentity:  identity,
				columnKeyAppWorkload:  workload,
				columnKeyAppRevision:  app.LatestRevision,
				columnKeyAppFQDN:      fqdn,
			})
		}
	} else {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				columnKeyAppName:      "Loading...",
				columnKeyAppRG:        "",
				columnKeyAppLocation:  "",
				columnKeyAppStatus:    "",
				columnKeyAppReplicas:  "",
				columnKeyAppResources: "",
				columnKeyAppIngress:   "",
				columnKeyAppIdentity:  "",
				columnKeyAppWorkload:  "",
				columnKeyAppRevision:  "",
				columnKeyAppFQDN:      "",
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
		WithFilterInput(m.appsFilterInput).
		Focused(true)

	// Calculate height: total height - title (1) - help (1) - status (1) - margins (2)
	tableHeight := m.termH - 5
	if tableHeight > 0 {
		t = t.WithPageSize(tableHeight)
	}

	return t
}

// createRevisionsTable creates a table component for revisions
func (m model) createRevisionsTable() table.Model {
	// Calculate dynamic width for Revision column based on content
	revisionWidth := len("Revision")

	// Find maximum width needed for Revision column
	for _, rev := range m.revs {
		if len(rev.Name) > revisionWidth {
			revisionWidth = len(rev.Name)
		}
	}

	// Add padding and set minimum width
	revisionWidth += 2
	if revisionWidth < 15 {
		revisionWidth = 15
	}

	columns := []table.Column{
		table.NewColumn(columnKeyRevName, "Revision", revisionWidth).WithFiltered(true), // Dynamic width based on content
		table.NewColumn(columnKeyRevActive, "Active", 8),                                // Fixed width
		table.NewColumn(columnKeyRevTraffic, "Traffic", 10),                             // Fixed width
		table.NewColumn(columnKeyRevReplicas, "Replicas", 10),                           // Fixed width
		table.NewColumn(columnKeyRevScaling, "Scaling", 12),                             // Fixed width
		table.NewColumn(columnKeyRevResources, "Resources", 15),                         // Fixed width
		table.NewColumn(columnKeyRevHealth, "Health", 12).WithFiltered(true),            // Fixed width
		table.NewColumn(columnKeyRevRunning, "Running", 15).WithFiltered(true),          // Fixed width
		table.NewColumn(columnKeyRevCreated, "Created", 20),                             // Fixed width
		table.NewColumn(columnKeyRevStatus, "Status", 15).WithFiltered(true),            // Fixed width
		table.NewColumn(columnKeyRevFQDN, "FQDN", 60),                                   // Fixed width (longest content)
	}

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
				columnKeyRevHealth:    health,
				columnKeyRevRunning:   running,
				columnKeyRevCreated:   created,
				columnKeyRevStatus:    status,
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
		Focused(true)

	// Only sort if we have actual data
	if len(m.revs) > 0 {
		t = t.SortByDesc(columnKeyRevTraffic)
	}

	// Calculate height: total height - title (1) - help (1) - status (1) - margins (2)
	tableHeight := m.termH - 5
	if tableHeight > 0 {
		t = t.WithPageSize(tableHeight)
	}

	return t
}

// createContainersTable creates a table component for containers
func (m model) createContainersTable() table.Model {
	// Calculate dynamic width for Container column based on content
	containerWidth := len("Container")

	// Find maximum width needed for Container column
	for _, ctr := range m.ctrs {
		if len(ctr.Name) > containerWidth {
			containerWidth = len(ctr.Name)
		}
	}

	// Add padding and set minimum width
	containerWidth += 2
	if containerWidth < 12 {
		containerWidth = 12
	}

	columns := []table.Column{
		table.NewColumn(columnKeyCtrName, "Container", containerWidth).WithFiltered(true), // Dynamic width based on content
		table.NewColumn(columnKeyCtrStatus, "Status", 10).WithFiltered(true),              // Fixed width - moved to second position
		table.NewColumn(columnKeyCtrImage, "Image", 50).WithFiltered(true),                // Fixed width (longest content)
		table.NewColumn(columnKeyCtrCommand, "Command", 25),                               // Fixed width
		table.NewColumn(columnKeyCtrArgs, "Args", 25),                                     // Fixed width
		table.NewColumn(columnKeyCtrResources, "Resources", 15),                           // Fixed width
		table.NewColumn(columnKeyCtrEnvCount, "Env", 8),                                   // Fixed width
		table.NewColumn(columnKeyCtrProbes, "Probes", 12),                                 // Fixed width
		table.NewColumn(columnKeyCtrVolumes, "Volumes", 10),                               // Fixed width
	}

	var rows []table.Row
	if len(m.ctrs) > 0 {
		rows = make([]table.Row, len(m.ctrs))
		for i, ctr := range m.ctrs {
			command := strings.Join(ctr.Command, " ")
			if command == "" {
				command = "-"
			}

			args := strings.Join(ctr.Args, " ")
			if args == "" {
				args = "-"
			}

			// Resources
			resources := fmt.Sprintf("%.2gC/%.1s", ctr.CPU, ctr.Memory)
			if ctr.CPU == 0 {
				resources = "-"
			}

			// Environment variables count
			envCount := fmt.Sprintf("%d", len(ctr.Env))
			if len(ctr.Env) == 0 {
				envCount = "-"
			}

			// Probes
			probes := strings.Join(ctr.Probes, ",")
			if probes == "" {
				probes = "-"
			}

			// Volume mounts count
			volumes := fmt.Sprintf("%d", len(ctr.VolumeMounts))
			if len(ctr.VolumeMounts) == 0 {
				volumes = "-"
			}

			rows[i] = table.NewRow(table.RowData{
				columnKeyCtrName:      ctr.Name,
				columnKeyCtrImage:     ctr.Image,
				columnKeyCtrCommand:   command,
				columnKeyCtrArgs:      args,
				columnKeyCtrResources: resources,
				columnKeyCtrEnvCount:  envCount,
				columnKeyCtrProbes:    probes,
				columnKeyCtrVolumes:   volumes,
				columnKeyCtrStatus:    "Running", // Default status
			})
		}
	} else {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				columnKeyCtrName:      "No containers",
				columnKeyCtrImage:     "",
				columnKeyCtrCommand:   "",
				columnKeyCtrArgs:      "",
				columnKeyCtrResources: "",
				columnKeyCtrEnvCount:  "",
				columnKeyCtrProbes:    "",
				columnKeyCtrVolumes:   "",
				columnKeyCtrStatus:    "",
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
		WithFilterInput(m.containersFilterInput).
		Focused(true)

	// Calculate height: total height - title (1) - help (1) - status (1) - margins (2)
	tableHeight := m.termH - 5
	if tableHeight > 0 {
		t = t.WithPageSize(tableHeight)
	}

	return t
}

// Helper methods for component factories

func (m model) getCurrentAppName() string {
	if app, ok := m.currentApp(); ok {
		return app.Name
	}
	return ""
}

func (m model) getCurrentApp() models.ContainerApp {
	if app, ok := m.currentApp(); ok {
		return app
	}
	return models.ContainerApp{}
}

// Confirmation dialog helper
func (m model) confirmBox() string {
	if !m.confirm.Visible {
		return ""
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center).
		Render(m.confirm.Text + "\n\n[y] Yes  [n] No")

	return box
}

// Status bar component factory
func (m model) createStatusBar() string {
	w := lipgloss.Width

	// Mode indicator
	var modeIndicator string
	switch m.mode {
	case modeApps:
		modeIndicator = modeAppsStyle.Render("APPS")
	case modeRevs:
		modeIndicator = modeRevisionsStyle.Render("REVISIONS")
	case modeContainers:
		modeIndicator = modeContainersStyle.Render("CONTAINERS")
	}

	// Status indicator
	var statusIndicator string
	var filterActive bool

	// Check if any filter input is focused
	switch m.mode {
	case modeApps:
		filterActive = m.appsFilterInput.Focused()
	case modeRevs:
		filterActive = m.revisionsFilterInput.Focused()
	case modeContainers:
		filterActive = m.containersFilterInput.Focused()
	}

	if filterActive {
		statusIndicator = statusLoadingStyle.Render("Filtering")
	} else if m.loading {
		spinner := m.createSpinner()
		statusIndicator = statusLoadingStyle.Render("Loading " + spinner.View())
	} else if m.err != nil {
		statusIndicator = statusErrorStyle.Render("Error")
	} else {
		statusIndicator = statusReadyStyle.Render("Ready")
	}

	// Count indicator
	var countIndicator string
	switch m.mode {
	case modeApps:
		if len(m.apps) == 1 {
			countIndicator = countStyle.Render("1 App")
		} else {
			countIndicator = countStyle.Render(fmt.Sprintf("%d Apps", len(m.apps)))
		}
	case modeRevs:
		if len(m.revs) == 1 {
			countIndicator = countStyle.Render("1 Revision")
		} else {
			countIndicator = countStyle.Render(fmt.Sprintf("%d Revisions", len(m.revs)))
		}
	case modeContainers:
		if len(m.ctrs) == 1 {
			countIndicator = countStyle.Render("1 Container")
		} else {
			countIndicator = countStyle.Render(fmt.Sprintf("%d Containers", len(m.ctrs)))
		}
	}

	// Context indicator (for deeper navigation levels)
	var contextIndicator string
	switch m.mode {
	case modeRevs:
		if appName := m.getCurrentAppName(); appName != "" {
			contextIndicator = contextStyle.Render("App: " + appName)
		}
	case modeContainers:
		if appName := m.getCurrentAppName(); appName != "" && m.currentRevName != "" {
			contextIndicator = contextStyle.Render(fmt.Sprintf("App: %s@%s", appName, m.currentRevName))
		}
	}

	// Resource group indicator
	var rgIndicator string
	if m.rg != "" {
		rgIndicator = rgStyle.Render("RG: " + m.rg)
	}

	// Calculate widths for fixed elements
	fixedWidth := w(modeIndicator) + w(statusIndicator) + w(countIndicator) + w(rgIndicator)
	if contextIndicator != "" {
		fixedWidth += w(contextIndicator)
	}

	// Status message (expandable middle section)
	statusMessage := m.statusLine
	if statusMessage == "" {
		if m.err != nil {
			statusMessage = m.err.Error()
		} else if m.loading {
			switch m.mode {
			case modeApps:
				statusMessage = "Loading container apps..."
			case modeRevs:
				statusMessage = "Loading revisions..."
			case modeContainers:
				statusMessage = "Loading containers..."
			}
		} else {
			// When ready, show empty message since status indicator already shows "Ready"
			statusMessage = ""
		}
	}

	// Create expandable status text
	statusVal := statusText.
		Width(max(0, m.termW-fixedWidth-4)). // Leave some margin
		Render(statusMessage)

	// Build the status bar
	var elements []string
	elements = append(elements, modeIndicator)
	elements = append(elements, statusIndicator)
	if contextIndicator != "" {
		elements = append(elements, contextIndicator)
	}
	elements = append(elements, statusVal)
	elements = append(elements, countIndicator)
	if rgIndicator != "" {
		elements = append(elements, rgIndicator)
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top, elements...)
	return statusBarStyle.Width(m.termW).Render(bar)
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Help bar component factory - now uses bubble help
func (m model) createHelpBar() string {
	return m.help.View(m.keys)
}

// General layout manager for consistent UI structure
func (m model) createLayout(title string, content string) string {
	// Create the main content area (title + content only)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		styleTitle.Render(title),
		content,
	)

	// Create the bottom bars
	helpBar := m.createHelpBar()
	statusBar := m.createStatusBar()

	// Calculate available height for main content (total height - help bar - status bar)
	helpBarHeight := lipgloss.Height(helpBar)
	statusBarHeight := lipgloss.Height(statusBar)
	mainContentHeight := m.termH - helpBarHeight - statusBarHeight

	// Position main content at top, help bar and status bar at bottom
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Height(mainContentHeight).Render(mainContent),
		helpBar,
		statusBar,
	)

	if m.confirm.Visible {
		return lipgloss.Place(m.termW, m.termH, lipgloss.Center, lipgloss.Center, m.confirmBox())
	}
	return body
}

// Loading state layout
func (m model) createLoadingLayout(title string, message string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styleAccent.Render(message),
		"",
	)
	return m.createLayout(title, content)
}

// Error state layout
func (m model) createErrorLayout(title string, errorMsg string, helpMsg string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		StyleError.Render("Error: ")+errorMsg,
		styleAccent.Render(helpMsg),
		"",
	)
	return m.createLayout(title, content)
}

// Table layout
func (m model) createTableLayout(title string, tableView string) string {
	return m.createLayout(title, tableView)
}
