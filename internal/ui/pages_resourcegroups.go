package ui

import (
	"sort"
	"strings"
	"unicode/utf8"

	models "github.com/IAL32/az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	maxTagsDisplayLength = 120
	truncationSuffix     = "..."
	noTagsPlaceholder    = "-"
	tagSeparator         = ", "
	tagKeyValueSeparator = "="
)

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

// createResourceGroupsTable creates a table component for resource groups
func (m model) createResourceGroupsTable() table.Model {
	// Create dynamic column builder
	builder := NewDynamicColumnBuilder().
		AddColumn(columnKeyRGName, "Name", 25, true).         // Dynamic width, min 25
		AddColumn(columnKeyRGLocation, "Location", 20, true). // Fixed width
		AddColumn(columnKeyRGState, "State", 15, true).       // Fixed width
		AddColumn(columnKeyRGTags, "Tags", 80, false)         // Fixed width - much larger for tags

	// Update dynamic column widths based on actual content
	for _, rg := range m.resourceGroups {
		builder.UpdateWidthFromString(columnKeyRGName, rg.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(m.resourceGroups) > 0 {
		rows = make([]table.Row, len(m.resourceGroups))
		for i, rg := range m.resourceGroups {
			state := rg.State
			if state == "" {
				state = "Unknown"
			}

			tags := formatResourceGroupTags(rg.Tags)

			rows[i] = table.NewRow(table.RowData{
				columnKeyRGName:     rg.Name,
				columnKeyRGLocation: rg.Location,
				columnKeyRGState:    table.NewStyledCell(state, lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor(state)))),
				columnKeyRGTags:     tags,
			})
		}
	}
	// Don't show any placeholder rows - empty table is fine

	return m.
		createUnifiedTable(columns, rows, m.resourceGroupsFilterInput).
		SortByAsc(columnKeyRGName)
}

// getResourceGroupsHelpKeys returns the key bindings for resource groups mode
func (m model) getResourceGroupsHelpKeys() keyMap {
	return keyMap{
		Enter:       m.keys.Enter,
		Refresh:     m.keys.Refresh,
		Filter:      m.keys.Filter,
		ScrollLeft:  m.keys.ScrollLeft,
		ScrollRight: m.keys.ScrollRight,
		Help:        m.keys.Help,
		Quit:        m.keys.Quit,
	}
}

// Message handlers
func (m model) handleLoadedResourceGroupsMsg(msg loadedResourceGroupsMsg) (model, tea.Cmd) {
	m.loading = false
	m.err = msg.err
	if msg.err != nil {
		return m, nil
	}
	m.resourceGroups = msg.resourceGroups

	if len(m.resourceGroups) == 0 {
		return m, nil
	}

	// Create the resource groups table with the loaded data
	m.resourceGroupsTable = m.createResourceGroupsTable()

	return m, nil
}

// Key handlers
func (m model) handleResourceGroupsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	// Handle filter input when focused
	if m.resourceGroupsFilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.resourceGroupsFilterInput.Blur()
			// Sync the filter with the table after applying
			m.resourceGroupsTable = m.resourceGroupsTable.WithFilterInput(m.resourceGroupsFilterInput)
			return m, nil, true
		case "esc":
			m.resourceGroupsFilterInput.SetValue("")
			m.resourceGroupsFilterInput.Blur()
			m.resourceGroupsTable = m.resourceGroupsTable.WithFilterInput(m.resourceGroupsFilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.resourceGroupsFilterInput, cmd = m.resourceGroupsFilterInput.Update(msg)
			m.resourceGroupsTable = m.resourceGroupsTable.WithFilterInput(m.resourceGroupsFilterInput)
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
		m.resourceGroupsFilterInput.SetValue("") // Clear any existing value
		m.resourceGroupsFilterInput.Focus()
		m.resourceGroupsTable = m.resourceGroupsTable.WithFilterInput(m.resourceGroupsFilterInput)
		return m, nil, true
	case "enter":
		if len(m.resourceGroups) == 0 {
			return m, nil, true
		}

		// Get selected resource group from table
		selectedRow := m.resourceGroupsTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		rgName, ok := selectedRow.Data[columnKeyRGName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the resource group by name
		var selectedRG models.ResourceGroup
		found := false
		for _, rg := range m.resourceGroups {
			if rg.Name == rgName {
				selectedRG = rg
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		// Navigate to apps for this resource group
		m.mode = modeApps
		m.rg = selectedRG.Name
		m.loading = true
		m.apps = nil
		m.appsTable = m.createAppsTable()
		return m, LoadAppsCmd(m.dataProvider, m.rg), true

	case "r":
		// Refresh resource groups list - clear data and show loading state
		m.loading = true
		m.resourceGroups = nil
		m.resourceGroupsTable = m.createResourceGroupsTable()
		return m, LoadResourceGroupsCmd(m.dataProvider), true
	}

	return m, nil, false
}

// View functions
func (m model) viewResourceGroups() string {
	if m.loading && len(m.resourceGroups) == 0 {
		// Show loading state using generalized layout
		return m.createLoadingLayout("Loading resource groups...")
	}

	if m.err != nil && len(m.resourceGroups) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(m.err.Error(), "Press r to retry or q to quit.")
	}

	// Show table view using generalized layout
	tableView := m.resourceGroupsTable.View()
	return m.createTableLayout(tableView)
}
