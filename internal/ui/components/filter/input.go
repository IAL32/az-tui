// Package filter provides reusable filter input components for the UI.
package filter

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// FilterInput wraps a textinput.Model with additional filter-specific functionality
type FilterInput struct {
	textinput.Model
}

// NewFilterInput creates a new filter input with the given placeholder
func NewFilterInput(placeholder string) FilterInput {
	input := textinput.New()
	input.Placeholder = placeholder
	return FilterInput{Model: input}
}

// HandleFilterKey handles common filter input key events
// Returns (updated FilterInput, command, handled bool)
func (f FilterInput) HandleFilterKey(msg tea.KeyMsg, table table.Model) (FilterInput, table.Model, tea.Cmd, bool) {
	switch msg.String() {
	case "enter":
		f.Model.Blur()
		// Sync the filter with the table after applying
		table = table.WithFilterInput(f.Model)
		return f, table, nil, true
	case "esc":
		f.Model.SetValue("")
		f.Model.Blur()
		table = table.WithFilterInput(f.Model)
		return f, table, nil, true
	default:
		var cmd tea.Cmd
		f.Model, cmd = f.Model.Update(msg)
		table = table.WithFilterInput(f.Model)
		return f, table, cmd, true
	}
}

// ActivateFilter activates the filter input by clearing any existing value and focusing
func (f FilterInput) ActivateFilter(table table.Model) (FilterInput, table.Model) {
	f.Model.SetValue("") // Clear any existing value
	f.Model.Focus()
	table = table.WithFilterInput(f.Model)
	return f, table
}

// IsActive returns true if the filter input is currently focused
func (f FilterInput) IsActive() bool {
	return f.Model.Focused()
}

// FilterState represents the state of filter inputs across different pages
type FilterState struct {
	ResourceGroups FilterInput
	Apps           FilterInput
	Revisions      FilterInput
	Containers     FilterInput
	EnvVars        FilterInput
}

// NewFilterState creates a new filter state with all inputs initialized
func NewFilterState() FilterState {
	return FilterState{
		ResourceGroups: NewFilterInput("Filter resource groups..."),
		Apps:           NewFilterInput("Filter apps..."),
		Revisions:      NewFilterInput("Filter revisions..."),
		Containers:     NewFilterInput("Filter containers..."),
		EnvVars:        NewFilterInput("Filter environment variables..."),
	}
}

// IsAnyActive returns true if any filter input is currently active
func (fs FilterState) IsAnyActive() bool {
	return fs.ResourceGroups.IsActive() ||
		fs.Apps.IsActive() ||
		fs.Revisions.IsActive() ||
		fs.Containers.IsActive() ||
		fs.EnvVars.IsActive()
}
