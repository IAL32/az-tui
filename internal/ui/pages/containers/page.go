package containers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/IAL32/az-tui/internal/models"
	tablebuilder "github.com/IAL32/az-tui/internal/ui/components/table"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/IAL32/az-tui/internal/ui/pages"
)

// ContainersPage represents the containers page using the new page interface system.
// It displays containers in an actionable table format with logs, exec, and envvars actions.
type ContainersPage struct {
	*pages.ActionablePage[models.Container]

	// Navigation context
	appName      string
	appID        string
	revisionName string

	// Layout system
	layoutSystem *layouts.LayoutSystem

	// Key bindings
	keys ContainersKeyMap

	// Action functions
	showLogsFunc          func(models.Container) tea.Cmd
	execIntoContainerFunc func(models.Container) tea.Cmd

	// Navigation functions
	navigateToEnvVarsFunc func(models.Container) tea.Cmd
	backToRevisionsFunc   func() tea.Cmd
}

// ContainersKeyMap defines the key bindings for the containers page
type ContainersKeyMap struct {
	Logs        key.Binding
	Exec        key.Binding
	EnvVars     key.Binding
	Refresh     key.Binding
	Filter      key.Binding
	ScrollLeft  key.Binding
	ScrollRight key.Binding
	Help        key.Binding
	Back        key.Binding
	Quit        key.Binding
}

// NewContainersPage creates a new containers page
func NewContainersPage(layoutSystem *layouts.LayoutSystem) *ContainersPage {
	// Create the base actionable page
	basePage := pages.NewActionablePage[models.Container]("Filter containers...")

	// Create the containers page
	page := &ContainersPage{
		ActionablePage: basePage,
		layoutSystem:   layoutSystem,
		keys:           defaultContainersKeyMap(),
	}

	// Set the table creation function
	page.SetCreateTableFunc(page.createContainersTable)

	// Disable navigation (containers don't navigate on enter, they use specific actions)
	// page.SetNavigationFunc is not called

	// Set up actions
	page.setupActions()

	return page
}

// defaultContainersKeyMap returns the default key bindings for containers
func defaultContainersKeyMap() ContainersKeyMap {
	return ContainersKeyMap{
		Logs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "logs"),
		),
		Exec: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "exec"),
		),
		EnvVars: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "env vars"),
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

// SetRevisionContext sets the revision context for the containers page
func (p *ContainersPage) SetRevisionContext(appName, appID, revisionName string) {
	p.appName = appName
	p.appID = appID
	p.revisionName = revisionName
}

// SetShowLogsFunc sets the function to call for showing logs
func (p *ContainersPage) SetShowLogsFunc(fn func(models.Container) tea.Cmd) {
	p.showLogsFunc = fn
}

// SetExecIntoContainerFunc sets the function to call for exec into container
func (p *ContainersPage) SetExecIntoContainerFunc(fn func(models.Container) tea.Cmd) {
	p.execIntoContainerFunc = fn
}

// SetNavigateToEnvVarsFunc sets the function to call when navigating to env vars
func (p *ContainersPage) SetNavigateToEnvVarsFunc(fn func(models.Container) tea.Cmd) {
	p.navigateToEnvVarsFunc = fn
}

// SetBackToRevisionsFunc sets the function to call when going back to revisions
func (p *ContainersPage) SetBackToRevisionsFunc(fn func() tea.Cmd) {
	p.backToRevisionsFunc = fn
	p.SetBackFunc(fn)
}

// Action setup

// setupActions configures the available actions for the containers page
func (p *ContainersPage) setupActions() {
	// Add logs action
	p.AddAction("logs", p.keys.Logs, func(container models.Container) tea.Cmd {
		if p.showLogsFunc != nil {
			return p.showLogsFunc(container)
		}
		return nil
	})

	// Add exec action
	p.AddAction("exec", p.keys.Exec, func(container models.Container) tea.Cmd {
		if p.execIntoContainerFunc != nil {
			return p.execIntoContainerFunc(container)
		}
		return nil
	})

	// Add env vars action
	p.AddAction("envvars", p.keys.EnvVars, func(container models.Container) tea.Cmd {
		if p.navigateToEnvVarsFunc != nil {
			return p.navigateToEnvVarsFunc(container)
		}
		return nil
	})
}

// Table creation methods

// createContainersTable creates a table for displaying containers
func (p *ContainersPage) createContainersTable(data []models.Container) table.Model {
	// Create dynamic column builder
	builder := tablebuilder.NewDynamicColumnBuilder().
		AddColumn("name", "Container", 12, true).       // Dynamic width, min 12
		AddColumn("status", "Status", 10, true).        // Fixed width - moved to second position
		AddColumn("image", "Image", 50, true).          // Fixed width (longest content)
		AddColumn("command", "Command", 25, false).     // Fixed width
		AddColumn("args", "Args", 25, false).           // Fixed width
		AddColumn("resources", "Resources", 15, false). // Fixed width
		AddColumn("envcount", "Env", 8, false).         // Fixed width
		AddColumn("probes", "Probes", 12, false).       // Fixed width
		AddColumn("volumes", "Volumes", 10, false)      // Fixed width

	// Update dynamic column widths based on actual content
	for _, ctr := range data {
		builder.UpdateWidthFromString("name", ctr.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(data) > 0 {
		rows = make([]table.Row, len(data))
		for i, ctr := range data {
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
				"name":      ctr.Name,
				"image":     ctr.Image,
				"command":   command,
				"args":      args,
				"resources": resources,
				"envcount":  envCount,
				"probes":    probes,
				"volumes":   volumes,
				"status":    table.NewStyledCell("Running", lipgloss.NewStyle().Foreground(pages.GetStatusColor("Running"))),
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

// Event handling methods

// HandleKeyMsg handles key messages for the containers page
func (p *ContainersPage) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// First, try base actionable page key handling
	if cmd, handled := p.ActionablePage.HandleKeyMsg(msg); handled {
		return cmd, handled
	}

	// Handle containers-specific keys
	switch msg.String() {
	case "esc":
		if p.backToRevisionsFunc != nil {
			return p.backToRevisionsFunc(), true
		}
		return nil, true
	case "?":
		// Help toggle - let the parent handle this
		return nil, false
	}

	return nil, false
}

// GetHelpKeys returns the help keys for the containers page
func (p *ContainersPage) GetHelpKeys() []key.Binding {
	baseKeys := []key.Binding{
		p.keys.Refresh,
		p.keys.Filter,
		p.keys.ScrollLeft,
		p.keys.ScrollRight,
		p.keys.Help,
		p.keys.Back,
		p.keys.Quit,
	}

	// Add action keys (includes logs, exec, and envvars)
	actionKeys := p.GetActionKeys()
	return append(baseKeys, actionKeys...)
}

// View rendering methods

// View renders the containers page
func (p *ContainersPage) View() string {
	// Use default help context (ShowAll = false)
	return p.ViewWithHelpContext(layouts.HelpContext{
		Mode: "Containers",
	})
}

// ViewWithHelpContext renders the containers page with help context
func (p *ContainersPage) ViewWithHelpContext(helpContext layouts.HelpContext) string {
	// Ensure the mode is set correctly
	helpContext.Mode = "Containers"

	// Handle loading state
	if p.IsLoading() {
		return p.layoutSystem.CreateLoadingLayout(
			"Loading containers...",
			layouts.StatusContext{
				Mode: "Containers",
				ContextInfo: map[string]string{
					"app":      p.appName,
					"revision": p.revisionName,
				},
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
				Mode:  "Containers",
				Error: err,
				ContextInfo: map[string]string{
					"app":      p.appName,
					"revision": p.revisionName,
				},
			},
			helpContext,
		)
	}

	// Render the table view
	tableView := p.GetTable().View()
	return p.layoutSystem.CreateTableLayout(
		tableView,
		layouts.StatusContext{
			Mode: "Containers",
			ContextInfo: map[string]string{
				"app":      p.appName,
				"revision": p.revisionName,
			},
			Counters: map[string]int{"count": len(p.GetData())},
		},
		helpContext,
	)
}

// Helper functions
