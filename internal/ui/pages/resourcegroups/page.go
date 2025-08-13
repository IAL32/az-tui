package resourcegroups

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/IAL32/az-tui/internal/models"
	tablebuilder "github.com/IAL32/az-tui/internal/ui/components/table"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/IAL32/az-tui/internal/ui/pages"
)

const (
	maxTagsDisplayLength = 120
	truncationSuffix     = "..."
	noTagsPlaceholder    = "-"
	tagSeparator         = ", "
	tagKeyValueSeparator = "="
)

// ResourceGroupsPage represents the resource groups page using the new page interface system.
// It displays resource groups in a navigable table format and serves as the entry point.
type ResourceGroupsPage struct {
	*pages.NavigablePage[models.ResourceGroup]

	// Layout system
	layoutSystem *layouts.LayoutSystem

	// Key bindings
	keys ResourceGroupsKeyMap

	// Navigation function
	navigateToAppsFunc func(models.ResourceGroup) tea.Cmd
}

// ResourceGroupsKeyMap defines the key bindings for the resource groups page
type ResourceGroupsKeyMap struct {
	Enter       key.Binding
	Refresh     key.Binding
	Filter      key.Binding
	ScrollLeft  key.Binding
	ScrollRight key.Binding
	Help        key.Binding
	Quit        key.Binding
}

// NewResourceGroupsPage creates a new resource groups page
func NewResourceGroupsPage(layoutSystem *layouts.LayoutSystem) *ResourceGroupsPage {
	// Create the base navigable page
	basePage := pages.NewNavigablePage[models.ResourceGroup]("Filter resource groups...")

	// Create the resource groups page
	page := &ResourceGroupsPage{
		NavigablePage: basePage,
		layoutSystem:  layoutSystem,
		keys:          defaultResourceGroupsKeyMap(),
	}

	// Set the table creation function
	page.SetCreateTableFunc(page.createResourceGroupsTable)

	// Enable navigation
	page.SetNavigationFunc(page.handleNavigation)

	return page
}

// defaultResourceGroupsKeyMap returns the default key bindings for resource groups
func defaultResourceGroupsKeyMap() ResourceGroupsKeyMap {
	return ResourceGroupsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
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
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// Configuration methods

// SetNavigateToAppsFunc sets the function to call when navigating to apps
func (p *ResourceGroupsPage) SetNavigateToAppsFunc(fn func(models.ResourceGroup) tea.Cmd) {
	p.navigateToAppsFunc = fn
}

// Table creation methods

// createResourceGroupsTable creates a table for displaying resource groups
func (p *ResourceGroupsPage) createResourceGroupsTable(data []models.ResourceGroup) table.Model {
	// Create dynamic column builder
	builder := tablebuilder.NewDynamicColumnBuilder().
		AddColumn("name", "Name", 25, true).         // Dynamic width, min 25
		AddColumn("location", "Location", 20, true). // Fixed width
		AddColumn("state", "State", 15, true).       // Fixed width
		AddColumn("tags", "Tags", 80, false)         // Fixed width - much larger for tags

	// Update dynamic column widths based on actual content
	for _, rg := range data {
		builder.UpdateWidthFromString("name", rg.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(data) > 0 {
		rows = make([]table.Row, len(data))
		for i, rg := range data {
			state := rg.State
			if state == "" {
				state = "Unknown"
			}

			tags := formatResourceGroupTags(rg.Tags)

			rows[i] = table.NewRow(table.RowData{
				"name":     rg.Name,
				"location": rg.Location,
				"state":    table.NewStyledCell(state, lipgloss.NewStyle().Foreground(pages.GetStatusColor(state))),
				"tags":     tags,
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

	return tablebuilder.CreateUnifiedTable(config).SortByAsc("name")
}

// Navigation methods

// handleNavigation handles navigation to the selected resource group
func (p *ResourceGroupsPage) handleNavigation(rg models.ResourceGroup) tea.Cmd {
	if p.navigateToAppsFunc != nil {
		return p.navigateToAppsFunc(rg)
	}
	return nil
}

// Event handling methods

// GetHelpKeys returns the help keys for the resource groups page
func (p *ResourceGroupsPage) GetHelpKeys() []key.Binding {
	return []key.Binding{
		p.keys.Enter,
		p.keys.Refresh,
		p.keys.Filter,
		p.keys.ScrollLeft,
		p.keys.ScrollRight,
		p.keys.Help,
		p.keys.Quit,
	}
}

// View rendering methods

// View renders the resource groups page
func (p *ResourceGroupsPage) View() string {
	// Use default help context (ShowAll = false)
	return p.ViewWithHelpContext(layouts.HelpContext{
		Mode: "Resource Groups",
	})
}

// ViewWithHelpContext renders the resource groups page with help context
func (p *ResourceGroupsPage) ViewWithHelpContext(helpContext layouts.HelpContext) string {
	// Ensure the mode is set correctly
	helpContext.Mode = "Resource Groups"

	// Handle loading state
	if p.IsLoading() {
		return p.layoutSystem.CreateLoadingLayout(
			"Loading resource groups...",
			layouts.StatusContext{
				Mode: "Resource Groups",
			},
			helpContext,
		)
	}

	// Handle error state
	if err := p.GetError(); err != nil {
		return p.layoutSystem.CreateErrorLayout(
			err.Error(),
			"Press 'r' to retry or 'q' to quit",
			layouts.StatusContext{
				Mode:  "Resource Groups",
				Error: err,
			},
			helpContext,
		)
	}

	// Render the table view
	tableView := p.GetTable().View()
	return p.layoutSystem.CreateTableLayout(
		tableView,
		layouts.StatusContext{
			Mode:     "Resource Groups",
			Counters: map[string]int{"count": len(p.GetData())},
		},
		helpContext,
	)
}

// Helper functions for resource group formatting

// formatResourceGroupTags formats resource group tags for display
func formatResourceGroupTags(tags map[string]string) string {
	if len(tags) == 0 {
		return noTagsPlaceholder
	}

	tagPairs := make([]string, 0, len(tags))
	var builder strings.Builder

	for key, value := range tags {
		if key == "" {
			continue
		}

		displayValue := value
		if displayValue == "" {
			displayValue = `""`
		}

		builder.Reset()
		builder.WriteString(key)
		builder.WriteString(tagKeyValueSeparator)
		builder.WriteString(displayValue)

		tagPairs = append(tagPairs, builder.String())
	}

	if len(tagPairs) == 0 {
		return noTagsPlaceholder
	}

	sort.Strings(tagPairs)
	result := strings.Join(tagPairs, tagSeparator)

	return truncateStringUnicodeSafe(result, maxTagsDisplayLength, truncationSuffix)
}

// truncateStringUnicodeSafe truncates a string safely handling Unicode characters
func truncateStringUnicodeSafe(s string, maxLength int, suffix string) string {
	if utf8.RuneCountInString(s) <= maxLength {
		return s
	}

	suffixLen := utf8.RuneCountInString(suffix)
	if maxLength <= suffixLen {
		return string([]rune(suffix)[:maxLength])
	}

	targetLength := maxLength - suffixLen
	runes := []rune(s)

	if len(runes) <= targetLength {
		return s
	}

	return string(runes[:targetLength]) + suffix
}
