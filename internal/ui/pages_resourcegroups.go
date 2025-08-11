package ui

import (
	models "az-tui/internal/models"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

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

			// Format tags
			tags := "-"
			if len(rg.Tags) > 0 {
				var tagPairs []string
				for k, v := range rg.Tags {
					tagPairs = append(tagPairs, fmt.Sprintf("%s=%s", k, v))
				}
				tags = strings.Join(tagPairs, ", ")
				// Truncate if too long for the wider column
				if len(tags) > 120 {
					tags = tags[:117] + "..."
				}
			}

			rows[i] = table.NewRow(table.RowData{
				columnKeyRGName:     rg.Name,
				columnKeyRGLocation: rg.Location,
				columnKeyRGState:    table.NewStyledCell(state, lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor(state)))),
				columnKeyRGTags:     tags,
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
		WithFilterInput(m.resourceGroupsFilterInput).
		Focused(true).
		SortByAsc(columnKeyRGName)

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
		return m, LoadAppsCmd(m.rg), true

	case "r":
		// Refresh resource groups list - clear data and show loading state
		m.loading = true
		m.resourceGroups = nil
		m.resourceGroupsTable = m.createResourceGroupsTable()
		return m, LoadResourceGroupsCmd(), true
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
