package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

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

// createEnvVarsTable creates a table component for environment variables
func (m model) createEnvVarsTable() table.Model {
	// Calculate dynamic widths for Name and Value columns based on content
	nameWidth := len("Name")
	valueWidth := len("Value")

	// Find maximum width needed for Name and Value columns
	if m.currentContainerName != "" {
		for _, ctr := range m.ctrs {
			if ctr.Name == m.currentContainerName {
				for name, value := range ctr.Env {
					if len(name) > nameWidth {
						nameWidth = len(name)
					}
					if len(value) > valueWidth {
						valueWidth = len(value)
					}
				}
				break
			}
		}
	}

	// Add padding and set minimum widths
	nameWidth += 2
	valueWidth += 2
	if nameWidth < 15 {
		nameWidth = 15
	}
	if valueWidth < 20 {
		valueWidth = 20
	}

	columns := []table.Column{
		table.NewColumn("name", "Name", nameWidth).WithFiltered(true),    // Dynamic width based on content
		table.NewColumn("value", "Value", valueWidth).WithFiltered(true), // Dynamic width based on content
	}

	var rows []table.Row
	if m.currentContainerName != "" {
		// Find the current container and get its environment variables
		for _, ctr := range m.ctrs {
			if ctr.Name == m.currentContainerName {
				if len(ctr.Env) > 0 {
					rows = make([]table.Row, 0, len(ctr.Env))
					for name, value := range ctr.Env {
						rows = append(rows, table.NewRow(table.RowData{
							"name":  name,
							"value": value,
						}))
					}
				}
				break
			}
		}
	}

	if len(rows) == 0 {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				"name":  "No environment variables",
				"value": "",
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
		WithFilterInput(m.envVarsFilterInput).
		Focused(true)

	// Calculate height dynamically based on actual help and status bar heights
	helpBar := m.createHelpBar()
	statusBar := m.createStatusBar()
	helpBarHeight := lipgloss.Height(helpBar)
	statusBarHeight := lipgloss.Height(statusBar)

	// Available height = total height - help bar - status bar - table overhead (6 lines)
	// Conservative calculation to ensure table header stays visible
	// Accounts for: borders, header, filter area, and extra margin
	availableHeight := m.termH - helpBarHeight - statusBarHeight - 6
	if availableHeight > 0 {
		t = t.WithPageSize(availableHeight)
	}

	return t
}
