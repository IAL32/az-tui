package envvars

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"

	"github.com/IAL32/az-tui/internal/models"
	tablebuilder "github.com/IAL32/az-tui/internal/ui/components/table"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/IAL32/az-tui/internal/ui/pages"
)

// EnvVarsPage represents the environment variables page using the new page interface system.
// It displays environment variables for a specific container in a read-only table format.
type EnvVarsPage struct {
	*pages.ReadOnlyPage[map[string]string]

	// Navigation context
	containerName string
	containers    []models.Container

	// Layout system
	layoutSystem *layouts.LayoutSystem

	// Key bindings
	keys EnvVarsKeyMap

	// Back navigation function
	backFunc func() tea.Cmd
}

// EnvVarsKeyMap defines the key bindings for the environment variables page
type EnvVarsKeyMap struct {
	Filter      key.Binding
	ScrollLeft  key.Binding
	ScrollRight key.Binding
	Help        key.Binding
	Back        key.Binding
	Quit        key.Binding
}

// NewEnvVarsPage creates a new environment variables page
func NewEnvVarsPage(layoutSystem *layouts.LayoutSystem) *EnvVarsPage {
	// Create the base read-only page
	basePage := pages.NewReadOnlyPage[map[string]string]("Filter environment variables...")

	// Create the envvars page
	page := &EnvVarsPage{
		ReadOnlyPage: basePage,
		layoutSystem: layoutSystem,
		keys:         defaultEnvVarsKeyMap(),
	}

	// Set the table creation function
	page.SetCreateTableFunc(page.createEnvVarsTable)

	return page
}

// defaultEnvVarsKeyMap returns the default key bindings for environment variables
func defaultEnvVarsKeyMap() EnvVarsKeyMap {
	return EnvVarsKeyMap{
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

// SetContainerContext sets the container context for the environment variables page
func (p *EnvVarsPage) SetContainerContext(containerName string, containers []models.Container) {
	p.containerName = containerName
	p.containers = containers
	p.loadEnvVarsData()
}

// SetBackFunc sets the function to call when navigating back
func (p *EnvVarsPage) SetBackFunc(fn func() tea.Cmd) {
	p.backFunc = fn
}

// Data loading methods

// loadEnvVarsData loads environment variables data from the current container
func (p *EnvVarsPage) loadEnvVarsData() {
	if p.containerName == "" {
		p.SetData([]map[string]string{})
		return
	}

	// Find the current container and get its environment variables
	for _, ctr := range p.containers {
		if ctr.Name == p.containerName {
			if len(ctr.Env) > 0 {
				// Convert the environment variables map to a slice of maps for the table
				envVars := []map[string]string{ctr.Env}
				p.SetData(envVars)
			} else {
				// Set empty data if no environment variables
				p.SetData([]map[string]string{})
			}
			return
		}
	}

	// Container not found, set empty data
	p.SetData([]map[string]string{})
}

// Table creation methods

// createEnvVarsTable creates a table for displaying environment variables
func (p *EnvVarsPage) createEnvVarsTable(data []map[string]string) table.Model {
	// Create dynamic column builder
	builder := tablebuilder.NewDynamicColumnBuilder().
		AddColumn("name", "Name", 15, true).  // Min width 15, with filter
		AddColumn("value", "Value", 20, true) // Min width 20, with filter

	var rows []table.Row

	// Process environment variables data
	if len(data) > 0 && len(data[0]) > 0 {
		envVars := data[0] // Get the first (and only) map of environment variables

		// Update column widths based on actual content
		for name, value := range envVars {
			builder.UpdateWidthFromString("name", name)
			builder.UpdateWidthFromString("value", value)
		}

		// Create rows from environment variables
		rows = make([]table.Row, 0, len(envVars))
		for name, value := range envVars {
			rows = append(rows, table.NewRow(table.RowData{
				"name":  name,
				"value": value,
			}))
		}
	}

	// If no environment variables, create a placeholder row
	if len(rows) == 0 {
		rows = []table.Row{
			table.NewRow(table.RowData{
				"name":  "No environment variables",
				"value": "",
			}),
		}
	}

	// Build columns with calculated widths
	columns := builder.Build()

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

// HandleKeyMsg handles key messages for the environment variables page
func (p *EnvVarsPage) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// Handle environment variables specific keys FIRST (before base page)
	// This prevents ReadOnlyPage from intercepting our navigation keys
	switch msg.String() {
	case "esc":
		if p.backFunc != nil {
			return p.backFunc(), true
		}
		return nil, true
	case "?":
		// Help toggle - let the parent handle this
		return nil, false
	case "j", "k", "up", "down":
		// Navigation keys - let the parent handle these
		return nil, false
	}

	// Now try base page key handling for other keys (like filtering)
	if cmd, handled := p.ReadOnlyPage.HandleKeyMsg(msg); handled {
		return cmd, handled
	}

	return nil, false
}

// GetHelpKeys returns the help keys for the environment variables page
func (p *EnvVarsPage) GetHelpKeys() []key.Binding {
	return []key.Binding{
		p.keys.Filter,
		p.keys.ScrollLeft,
		p.keys.ScrollRight,
		p.keys.Help,
		p.keys.Back,
		p.keys.Quit,
	}
}

// View rendering methods

// View renders the environment variables page
func (p *EnvVarsPage) View() string {
	// Use default help context (ShowAll = false)
	return p.ViewWithHelpContext(layouts.HelpContext{
		Mode: "Environment Variables",
	})
}

// ViewWithHelpContext renders the environment variables page with help context
func (p *EnvVarsPage) ViewWithHelpContext(helpContext layouts.HelpContext) string {
	// Ensure the mode is set correctly
	helpContext.Mode = "Environment Variables"

	// Handle loading state
	if p.IsLoading() {
		return p.layoutSystem.CreateLoadingLayout(
			"Loading environment variables...",
			layouts.StatusContext{
				Mode:        "Environment Variables",
				ContextInfo: map[string]string{"container": p.containerName},
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
				Mode:        "Environment Variables",
				Error:       err,
				ContextInfo: map[string]string{"container": p.containerName},
			},
			helpContext,
		)
	}

	// Render the table view
	tableView := p.GetTable().View()
	return p.layoutSystem.CreateTableLayout(
		tableView,
		layouts.StatusContext{
			Mode:        "Environment Variables",
			ContextInfo: map[string]string{"container": p.containerName},
			Counters:    map[string]int{"count": len(p.GetData())},
		},
		helpContext,
	)
}

// Refresh refreshes the environment variables data
func (p *EnvVarsPage) Refresh() tea.Cmd {
	p.SetLoading(true)
	p.SetError(nil)
	p.loadEnvVarsData()
	p.SetLoading(false)
	return nil
}

// Reset resets the page state
func (p *EnvVarsPage) Reset() {
	p.ReadOnlyPage.Reset()
	p.containerName = ""
	p.containers = nil
}
