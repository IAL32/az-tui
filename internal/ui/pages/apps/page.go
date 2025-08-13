package apps

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/IAL32/az-tui/internal/models"
	tablebuilder "github.com/IAL32/az-tui/internal/ui/components/table"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/IAL32/az-tui/internal/ui/pages"
)

// AppsPage represents the container apps page using the new page interface system.
// It displays container apps in an actionable table format with logs and exec actions.
type AppsPage struct {
	*pages.ActionablePage[models.ContainerApp]

	// Navigation context
	resourceGroupName string

	// Layout system
	layoutSystem *layouts.LayoutSystem

	// Key bindings
	keys AppsKeyMap

	// Action functions
	showLogsFunc    func(models.ContainerApp) tea.Cmd
	execIntoAppFunc func(models.ContainerApp) tea.Cmd

	// Navigation functions
	navigateToRevisionsFunc  func(models.ContainerApp) tea.Cmd
	backToResourceGroupsFunc func() tea.Cmd
}

// AppsKeyMap defines the key bindings for the apps page
type AppsKeyMap struct {
	Enter       key.Binding
	Logs        key.Binding
	Exec        key.Binding
	Refresh     key.Binding
	Filter      key.Binding
	ScrollLeft  key.Binding
	ScrollRight key.Binding
	Help        key.Binding
	Back        key.Binding
	Quit        key.Binding
}

// NewAppsPage creates a new apps page
func NewAppsPage(layoutSystem *layouts.LayoutSystem) *AppsPage {
	// Create the base actionable page
	basePage := pages.NewActionablePage[models.ContainerApp]("Filter container apps...")

	// Create the apps page
	page := &AppsPage{
		ActionablePage: basePage,
		layoutSystem:   layoutSystem,
		keys:           defaultAppsKeyMap(),
	}

	// Set the table creation function
	page.SetCreateTableFunc(page.createAppsTable)

	// Enable navigation
	page.SetNavigationFunc(page.handleNavigation)

	// Set up actions
	page.setupActions()

	return page
}

// defaultAppsKeyMap returns the default key bindings for apps
func defaultAppsKeyMap() AppsKeyMap {
	return AppsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "revisions"),
		),
		Logs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "logs"),
		),
		Exec: key.NewBinding(
			key.WithKeys("s", "e"),
			key.WithHelp("s/e", "exec"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ScrollLeft: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+←", "scroll left"),
		),
		ScrollRight: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+→", "scroll right"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// Configuration methods

// SetResourceGroupContext sets the resource group context for the apps page
func (p *AppsPage) SetResourceGroupContext(resourceGroupName string) {
	p.resourceGroupName = resourceGroupName
}

// SetShowLogsFunc sets the function to call for showing logs
func (p *AppsPage) SetShowLogsFunc(fn func(models.ContainerApp) tea.Cmd) {
	p.showLogsFunc = fn
}

// SetExecIntoAppFunc sets the function to call for exec into app
func (p *AppsPage) SetExecIntoAppFunc(fn func(models.ContainerApp) tea.Cmd) {
	p.execIntoAppFunc = fn
}

// SetNavigateToRevisionsFunc sets the function to call when navigating to revisions
func (p *AppsPage) SetNavigateToRevisionsFunc(fn func(models.ContainerApp) tea.Cmd) {
	p.navigateToRevisionsFunc = fn
}

// SetBackToResourceGroupsFunc sets the function to call when going back to resource groups
func (p *AppsPage) SetBackToResourceGroupsFunc(fn func() tea.Cmd) {
	p.backToResourceGroupsFunc = fn
	p.SetBackFunc(fn)
}

// Action setup

// setupActions configures the available actions for the apps page
func (p *AppsPage) setupActions() {
	// Add logs action
	p.AddAction("logs", p.keys.Logs, func(app models.ContainerApp) tea.Cmd {
		if p.showLogsFunc != nil {
			return p.showLogsFunc(app)
		}
		return nil
	})

	// Add exec action
	p.AddAction("exec", p.keys.Exec, func(app models.ContainerApp) tea.Cmd {
		if p.execIntoAppFunc != nil {
			return p.execIntoAppFunc(app)
		}
		return nil
	})
}

// Table creation methods

// createAppsTable creates a table for displaying container apps
func (p *AppsPage) createAppsTable(data []models.ContainerApp) table.Model {
	// Create dynamic column builder
	builder := tablebuilder.NewDynamicColumnBuilder().
		AddColumn("name", "Name", 15, true).                 // Dynamic width, min 15
		AddColumn("location", "Location", 15, true).         // Fixed width
		AddColumn("status", "Status", 12, true).             // Fixed width
		AddColumn("replicas", "Replicas", 10, false).        // Fixed width
		AddColumn("resources", "Resources", 12, false).      // Fixed width
		AddColumn("ingress", "Ingress", 18, false).          // Fixed width
		AddColumn("identity", "Identity", 15, false).        // Fixed width
		AddColumn("workload", "Workload", 15, false).        // Fixed width
		AddColumn("revision", "Latest Revision", 30, false). // Fixed width
		AddColumn("fqdn", "FQDN", 60, false)                 // Fixed width (longest content)

	// Update dynamic column widths based on actual content
	for _, app := range data {
		builder.UpdateWidthFromString("name", app.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(data) > 0 {
		rows = make([]table.Row, len(data))
		for i, app := range data {
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
				"name":      app.Name,
				"location":  app.Location,
				"status":    table.NewStyledCell(status, lipgloss.NewStyle().Foreground(getStatusColor(status))),
				"replicas":  replicas,
				"resources": resources,
				"ingress":   ingress,
				"identity":  identity,
				"workload":  workload,
				"revision":  app.LatestRevision,
				"fqdn":      fqdn,
			})
		}
	}

	// Get content dimensions
	contentWidth, contentHeight := p.layoutSystem.GetContentDimensions(layouts.LayoutOptions{})

	// Create the table using the unified table builder with theme styling
	config := tablebuilder.UnifiedTableConfig{
		Columns:     columns,
		Rows:        rows,
		FilterInput: p.GetFilterInput(),
		BaseStyle:   p.layoutSystem.GetStyle("tableBase"),
		MaxWidth:    contentWidth,
		MaxHeight:   contentHeight,
	}

	return tablebuilder.CreateUnifiedTable(config)
}

// Navigation methods

// handleNavigation handles navigation to the selected app's revisions
func (p *AppsPage) handleNavigation(app models.ContainerApp) tea.Cmd {
	if p.navigateToRevisionsFunc != nil {
		return p.navigateToRevisionsFunc(app)
	}
	return nil
}

// Event handling methods

// HandleKeyMsg handles key messages for the apps page
func (p *AppsPage) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// First, try base actionable page key handling
	if cmd, handled := p.ActionablePage.HandleKeyMsg(msg); handled {
		return cmd, handled
	}

	// Handle apps-specific keys
	switch msg.String() {
	case "esc":
		if p.backToResourceGroupsFunc != nil {
			return p.backToResourceGroupsFunc(), true
		}
		return nil, true
	case "?":
		// Help toggle - let the parent handle this
		return nil, false
	}

	return nil, false
}

// GetHelpKeys returns the help keys for the apps page
func (p *AppsPage) GetHelpKeys() []key.Binding {
	baseKeys := []key.Binding{
		p.keys.Enter,
		p.keys.Refresh,
		p.keys.Filter,
		p.keys.ScrollLeft,
		p.keys.ScrollRight,
		p.keys.Help,
		p.keys.Back,
		p.keys.Quit,
	}

	// Add action keys (includes logs and exec)
	actionKeys := p.GetActionKeys()
	return append(baseKeys, actionKeys...)
}

// View rendering methods

// View renders the apps page
func (p *AppsPage) View() string {
	// Use default help context (ShowAll = false)
	return p.ViewWithHelpContext(layouts.HelpContext{
		Mode: "Container Apps",
	})
}

// ViewWithHelpContext renders the apps page with help context
func (p *AppsPage) ViewWithHelpContext(helpContext layouts.HelpContext) string {
	// Ensure the mode is set correctly
	helpContext.Mode = "Container Apps"

	// Handle loading state
	if p.IsLoading() {
		return p.layoutSystem.CreateLoadingLayout(
			"Loading container apps...",
			layouts.StatusContext{
				Mode:        "Container Apps",
				ContextInfo: map[string]string{"resource_group": p.resourceGroupName},
			},
			helpContext,
		)
	}

	// Handle error state
	if err := p.GetError(); err != nil {
		return p.layoutSystem.CreateErrorLayout(
			err.Error(),
			"Press 'r' to retry or 'esc' to go back",
			layouts.StatusContext{
				Mode:        "Container Apps",
				Error:       err,
				ContextInfo: map[string]string{"resource_group": p.resourceGroupName},
			},
			helpContext,
		)
	}

	// Render the table view
	tableView := p.GetTable().View()
	return p.layoutSystem.CreateTableLayout(
		tableView,
		layouts.StatusContext{
			Mode:        "Container Apps",
			ContextInfo: map[string]string{"resource_group": p.resourceGroupName},
			Counters:    map[string]int{"count": len(p.GetData())},
		},
		helpContext,
	)
}

// Helper functions

// getStatusColor returns a color for the given status
func getStatusColor(status string) lipgloss.Color {
	switch status {
	case "Running", "Succeeded":
		return lipgloss.Color("#32CD32") // Success green
	case "Failed", "Error":
		return lipgloss.Color("#FF6B6B") // Error red
	case "Pending", "Creating", "Updating":
		return lipgloss.Color("#FFB347") // Warning yellow
	default:
		return lipgloss.Color("#808080") // Gray
	}
}
