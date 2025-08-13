package revisions

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

// RevisionsPage represents the revisions page using the new page interface system.
// It displays revisions in an actionable table format with restart, logs, and exec actions.
type RevisionsPage struct {
	*pages.ActionablePage[models.Revision]

	// Navigation context
	appName string
	appID   string

	// Layout system
	layoutSystem *layouts.LayoutSystem

	// Key bindings
	keys RevisionsKeyMap

	// Action functions
	restartRevisionFunc  func(models.Revision) tea.Cmd
	showLogsFunc         func(models.Revision) tea.Cmd
	execIntoRevisionFunc func(models.Revision) tea.Cmd

	// Navigation functions
	navigateToContainersFunc func(models.Revision) tea.Cmd
	backToAppsFunc           func() tea.Cmd
}

// RevisionsKeyMap defines the key bindings for the revisions page
type RevisionsKeyMap struct {
	Enter       key.Binding
	Restart     key.Binding
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

// NewRevisionsPage creates a new revisions page
func NewRevisionsPage(layoutSystem *layouts.LayoutSystem) *RevisionsPage {
	// Create the base actionable page
	basePage := pages.NewActionablePage[models.Revision]("Filter revisions...")

	// Create the revisions page
	page := &RevisionsPage{
		ActionablePage: basePage,
		layoutSystem:   layoutSystem,
		keys:           defaultRevisionsKeyMap(),
	}

	// Set the table creation function
	page.SetCreateTableFunc(page.createRevisionsTable)

	// Enable navigation
	page.SetNavigationFunc(page.handleNavigation)

	// Set up actions
	page.setupActions()

	return page
}

// defaultRevisionsKeyMap returns the default key bindings for revisions
func defaultRevisionsKeyMap() RevisionsKeyMap {
	return RevisionsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "containers"),
		),
		Restart: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "restart"),
		),
		Logs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "logs"),
		),
		Exec: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "exec"),
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

// SetAppContext sets the app context for the revisions page
func (p *RevisionsPage) SetAppContext(appName, appID string) {
	p.appName = appName
	p.appID = appID
}

// SetRestartRevisionFunc sets the function to call for restarting a revision
func (p *RevisionsPage) SetRestartRevisionFunc(fn func(models.Revision) tea.Cmd) {
	p.restartRevisionFunc = fn
}

// SetShowLogsFunc sets the function to call for showing logs
func (p *RevisionsPage) SetShowLogsFunc(fn func(models.Revision) tea.Cmd) {
	p.showLogsFunc = fn
}

// SetExecIntoRevisionFunc sets the function to call for exec into revision
func (p *RevisionsPage) SetExecIntoRevisionFunc(fn func(models.Revision) tea.Cmd) {
	p.execIntoRevisionFunc = fn
}

// SetNavigateToContainersFunc sets the function to call when navigating to containers
func (p *RevisionsPage) SetNavigateToContainersFunc(fn func(models.Revision) tea.Cmd) {
	p.navigateToContainersFunc = fn
}

// SetBackToAppsFunc sets the function to call when going back to apps
func (p *RevisionsPage) SetBackToAppsFunc(fn func() tea.Cmd) {
	p.backToAppsFunc = fn
	p.SetBackFunc(fn)
}

// Action setup

// setupActions configures the available actions for the revisions page
func (p *RevisionsPage) setupActions() {
	// Add restart action
	p.AddAction("restart", p.keys.Restart, func(rev models.Revision) tea.Cmd {
		if p.restartRevisionFunc != nil {
			return p.restartRevisionFunc(rev)
		}
		return nil
	})

	// Add logs action
	p.AddAction("logs", p.keys.Logs, func(rev models.Revision) tea.Cmd {
		if p.showLogsFunc != nil {
			return p.showLogsFunc(rev)
		}
		return nil
	})

	// Add exec action
	p.AddAction("exec", p.keys.Exec, func(rev models.Revision) tea.Cmd {
		if p.execIntoRevisionFunc != nil {
			return p.execIntoRevisionFunc(rev)
		}
		return nil
	})
}

// Table creation methods

// createRevisionsTable creates a table for displaying revisions
func (p *RevisionsPage) createRevisionsTable(data []models.Revision) table.Model {
	// Create dynamic column builder
	builder := tablebuilder.NewDynamicColumnBuilder().
		AddColumn("name", "Revision", 15, true).        // Dynamic width, min 15
		AddColumn("active", "Active", 8, false).        // Fixed width
		AddColumn("traffic", "Traffic", 10, false).     // Fixed width
		AddColumn("replicas", "Replicas", 10, false).   // Fixed width
		AddColumn("scaling", "Scaling", 12, false).     // Fixed width
		AddColumn("resources", "Resources", 15, false). // Fixed width
		AddColumn("health", "Health", 12, true).        // Fixed width
		AddColumn("running", "Running", 15, true).      // Fixed width
		AddColumn("created", "Created", 20, false).     // Fixed width
		AddColumn("status", "Status", 15, true).        // Fixed width
		AddColumn("fqdn", "FQDN", 60, false)            // Fixed width (longest content)

	// Update dynamic column widths based on actual content
	for _, rev := range data {
		builder.UpdateWidthFromString("name", rev.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(data) > 0 {
		rows = make([]table.Row, len(data))
		for i, rev := range data {
			activeMark := "❌"
			if rev.Active {
				activeMark = "✅"
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
				"name":      rev.Name,
				"active":    table.NewStyledCell(activeMark, lipgloss.NewStyle().Align(lipgloss.Center)),
				"traffic":   fmt.Sprintf("%d%%", rev.Traffic),
				"replicas":  replicas,
				"scaling":   scaling,
				"resources": resources,
				"health":    table.NewStyledCell(health, lipgloss.NewStyle().Foreground(getStatusColor(health))),
				"running":   table.NewStyledCell(running, lipgloss.NewStyle().Foreground(getStatusColor(running))),
				"created":   created,
				"status":    table.NewStyledCell(status, lipgloss.NewStyle().Foreground(getStatusColor(status))),
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

	return tablebuilder.CreateUnifiedTable(config).SortByDesc("traffic")
}

// Navigation methods

// handleNavigation handles navigation to the selected revision's containers
func (p *RevisionsPage) handleNavigation(rev models.Revision) tea.Cmd {
	if p.navigateToContainersFunc != nil {
		return p.navigateToContainersFunc(rev)
	}
	return nil
}

// Event handling methods

// HandleKeyMsg handles key messages for the revisions page
func (p *RevisionsPage) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// First, try base actionable page key handling
	if cmd, handled := p.ActionablePage.HandleKeyMsg(msg); handled {
		return cmd, handled
	}

	// Handle revisions-specific keys
	switch msg.String() {
	case "esc":
		if p.backToAppsFunc != nil {
			return p.backToAppsFunc(), true
		}
		return nil, true
	case "?":
		// Help toggle - let the parent handle this
		return nil, false
	}

	return nil, false
}

// GetHelpKeys returns the help keys for the revisions page
func (p *RevisionsPage) GetHelpKeys() []key.Binding {
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

	// Add action keys (includes restart, logs, and exec)
	actionKeys := p.GetActionKeys()
	return append(baseKeys, actionKeys...)
}

// View rendering methods

// View renders the revisions page
func (p *RevisionsPage) View() string {
	// Use default help context (ShowAll = false)
	return p.ViewWithHelpContext(layouts.HelpContext{
		Mode: "Revisions",
	})
}

// ViewWithHelpContext renders the revisions page with help context
func (p *RevisionsPage) ViewWithHelpContext(helpContext layouts.HelpContext) string {
	// Ensure the mode is set correctly
	helpContext.Mode = "Revisions"

	// Handle loading state
	if p.IsLoading() {
		return p.layoutSystem.CreateLoadingLayout(
			"Loading revisions...",
			layouts.StatusContext{
				Mode:        "Revisions",
				ContextInfo: map[string]string{"app": p.appName},
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
				Mode:        "Revisions",
				Error:       err,
				ContextInfo: map[string]string{"app": p.appName},
			},
			helpContext,
		)
	}

	// Render the table view
	tableView := p.GetTable().View()
	return p.layoutSystem.CreateTableLayout(
		tableView,
		layouts.StatusContext{
			Mode:        "Revisions",
			ContextInfo: map[string]string{"app": p.appName},
			Counters:    map[string]int{"count": len(p.GetData())},
		},
		helpContext,
	)
}

// Helper functions

// getStatusColor returns a color for the given status
func getStatusColor(status string) lipgloss.Color {
	switch status {
	case "Running", "Succeeded", "Healthy":
		return lipgloss.Color("#32CD32") // Success green
	case "Failed", "Error", "Unhealthy":
		return lipgloss.Color("#FF6B6B") // Error red
	case "Pending", "Creating", "Updating":
		return lipgloss.Color("#FFB347") // Warning yellow
	default:
		return lipgloss.Color("#808080") // Gray
	}
}
