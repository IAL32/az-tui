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
	// Initialize column widths with header lengths
	nameWidth := len("Name")
	rgWidth := len("Resource Group")
	locationWidth := len("Location")
	statusWidth := len("Status")
	replicasWidth := len("Replicas")
	resourcesWidth := len("Resources")
	ingressWidth := len("Ingress")
	identityWidth := len("Identity")
	workloadWidth := len("Workload")
	revisionWidth := len("Latest Revision")
	fqdnWidth := len("FQDN")

	// Find maximum width needed for each column based on actual data
	for _, app := range m.apps {
		if len(app.Name) > nameWidth {
			nameWidth = len(app.Name)
		}
		if len(app.ResourceGroup) > rgWidth {
			rgWidth = len(app.ResourceGroup)
		}
		if len(app.Location) > locationWidth {
			locationWidth = len(app.Location)
		}
		if len(app.LatestRevision) > revisionWidth {
			revisionWidth = len(app.LatestRevision)
		}

		// Calculate formatted values and check their lengths
		fqdn := app.IngressFQDN
		if fqdn == "" {
			fqdn = "-"
		}
		if len(fqdn) > fqdnWidth {
			fqdnWidth = len(fqdn)
		}

		status := app.RunningStatus
		if status == "" {
			status = app.ProvisioningState
		}
		if status == "" {
			status = "Unknown"
		}
		if len(status) > statusWidth {
			statusWidth = len(status)
		}

		// Format replicas and check length
		replicas := fmt.Sprintf("%d-%d", app.MinReplicas, app.MaxReplicas)
		if app.MinReplicas == 0 && app.MaxReplicas == 0 {
			replicas = "-"
		}
		if len(replicas) > replicasWidth {
			replicasWidth = len(replicas)
		}

		// Format resources and check length
		resources := fmt.Sprintf("%.2gC/%.1s", app.CPU, app.Memory)
		if app.CPU == 0 {
			resources = "-"
		}
		if len(resources) > resourcesWidth {
			resourcesWidth = len(resources)
		}

		// Format ingress and check length
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
		if len(ingress) > ingressWidth {
			ingressWidth = len(ingress)
		}

		// Format identity and check length
		identity := app.IdentityType
		if identity == "" {
			identity = "None"
		}
		if len(identity) > identityWidth {
			identityWidth = len(identity)
		}

		// Format workload profile and check length
		workload := app.WorkloadProfile
		if workload == "" {
			workload = "Consumption"
		}
		if len(workload) > workloadWidth {
			workloadWidth = len(workload)
		}
	}

	// Add padding to each column for readability
	nameWidth += 2
	rgWidth += 2
	locationWidth += 2
	statusWidth += 2
	replicasWidth += 2
	resourcesWidth += 2
	ingressWidth += 2
	identityWidth += 2
	workloadWidth += 2
	revisionWidth += 2
	fqdnWidth += 2

	// Set reasonable minimum widths to ensure readability
	if nameWidth < 12 {
		nameWidth = 12
	}
	if rgWidth < 15 {
		rgWidth = 15
	}
	if locationWidth < 12 {
		locationWidth = 12
	}
	if statusWidth < 10 {
		statusWidth = 10
	}
	if replicasWidth < 8 {
		replicasWidth = 8
	}
	if resourcesWidth < 10 {
		resourcesWidth = 10
	}
	if ingressWidth < 8 {
		ingressWidth = 8
	}
	if identityWidth < 8 {
		identityWidth = 8
	}
	if workloadWidth < 12 {
		workloadWidth = 12
	}
	if revisionWidth < 16 {
		revisionWidth = 16
	}
	if fqdnWidth < 20 {
		fqdnWidth = 20
	}

	columns := []table.Column{
		table.NewColumn(columnKeyAppName, "Name", nameWidth),
		table.NewColumn(columnKeyAppRG, "Resource Group", rgWidth),
		table.NewColumn(columnKeyAppLocation, "Location", locationWidth),
		table.NewColumn(columnKeyAppStatus, "Status", statusWidth),
		table.NewColumn(columnKeyAppReplicas, "Replicas", replicasWidth),
		table.NewColumn(columnKeyAppResources, "Resources", resourcesWidth),
		table.NewColumn(columnKeyAppIngress, "Ingress", ingressWidth),
		table.NewColumn(columnKeyAppIdentity, "Identity", identityWidth),
		table.NewColumn(columnKeyAppWorkload, "Workload", workloadWidth),
		table.NewColumn(columnKeyAppRevision, "Latest Revision", revisionWidth),
		table.NewColumn(columnKeyAppFQDN, "FQDN", fqdnWidth),
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
	// Use fixed column widths optimized for content, allowing horizontal scroll
	columns := []table.Column{
		table.NewColumn(columnKeyRevName, "Revision", 35),
		table.NewColumn(columnKeyRevActive, "Active", 8),
		table.NewColumn(columnKeyRevTraffic, "Traffic", 10),
		table.NewColumn(columnKeyRevCreated, "Created", 20),
		table.NewColumn(columnKeyRevStatus, "Status", 15),
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

			status := rev.Status
			if status == "" {
				status = "-"
			}

			rows[i] = table.NewRow(table.RowData{
				columnKeyRevName:    rev.Name,
				columnKeyRevActive:  activeMark,
				columnKeyRevTraffic: fmt.Sprintf("%d%%", rev.Traffic),
				columnKeyRevCreated: created,
				columnKeyRevStatus:  status,
			})
		}
	} else {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				columnKeyRevName:    "No revisions",
				columnKeyRevActive:  "",
				columnKeyRevTraffic: "",
				columnKeyRevCreated: "",
				columnKeyRevStatus:  "",
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
	// Use fixed column widths optimized for content, allowing horizontal scroll
	columns := []table.Column{
		table.NewColumn(columnKeyCtrName, "Container", 20),
		table.NewColumn(columnKeyCtrImage, "Image", 50),
		table.NewColumn(columnKeyCtrCommand, "Command", 30),
		table.NewColumn(columnKeyCtrArgs, "Args", 30),
		table.NewColumn(columnKeyCtrStatus, "Status", 12),
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

			rows[i] = table.NewRow(table.RowData{
				columnKeyCtrName:    ctr.Name,
				columnKeyCtrImage:   ctr.Image,
				columnKeyCtrCommand: command,
				columnKeyCtrArgs:    args,
				columnKeyCtrStatus:  "Running", // Default status
			})
		}
	} else {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				columnKeyCtrName:    "No containers",
				columnKeyCtrImage:   "",
				columnKeyCtrCommand: "",
				columnKeyCtrArgs:    "",
				columnKeyCtrStatus:  "",
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
