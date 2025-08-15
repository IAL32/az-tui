package pages

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// BasePage provides a base implementation of common page functionality.
// It implements the Page interface and can be embedded in specific page types
// to reduce code duplication and provide consistent behavior.
type BasePage struct {
	// State fields
	isLoading bool
	error     error

	// UI components
	table       table.Model
	filterInput textinput.Model

	// Configuration
	loadingMessage string
	errorMessage   string
}

// NewBasePage creates a new BasePage with the given filter placeholder.
func NewBasePage(filterPlaceholder string) *BasePage {
	filterInput := textinput.New()
	filterInput.Placeholder = filterPlaceholder

	return &BasePage{
		isLoading:   false,
		error:       nil,
		filterInput: filterInput,
	}
}

// Data Management methods

func (b *BasePage) IsLoading() bool {
	return b.isLoading
}

func (b *BasePage) SetLoading(loading bool) {
	b.isLoading = loading
}

func (b *BasePage) GetError() error {
	return b.error
}

func (b *BasePage) SetError(err error) {
	b.error = err
}

func (b *BasePage) ClearData() {
	// Base implementation - specific pages should override this
	b.error = nil
}

// Table Management methods

func (b *BasePage) GetTable() table.Model {
	return b.table
}

func (b *BasePage) SetTable(t table.Model) {
	b.table = t
}

func (b *BasePage) CreateTable() table.Model {
	// Base implementation - specific pages should override this
	return table.New([]table.Column{})
}

// Filter Management methods

func (b *BasePage) GetFilterInput() textinput.Model {
	return b.filterInput
}

func (b *BasePage) SetFilterInput(input textinput.Model) {
	b.filterInput = input
}

func (b *BasePage) ApplyFilter() {
	b.table = b.table.WithFilterInput(b.filterInput)
}

func (b *BasePage) ClearFilter() {
	b.filterInput.SetValue("")
	b.filterInput.Blur()
	b.table = b.table.WithFilterInput(b.filterInput)
}

// Event Handling methods

func (b *BasePage) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// Handle filter input when focused
	if b.filterInput.Focused() {
		return b.handleFilterInput(msg)
	}

	// Handle common keys
	switch msg.String() {
	case "/":
		return b.startFiltering(), true
	case "ctrl+c", "q":
		return tea.Quit, true
	}

	return nil, false
}

func (b *BasePage) HandleMessage(msg tea.Msg) tea.Cmd {
	// Base implementation - specific pages should override this
	return nil
}

// handleFilterInput handles key input when the filter is focused
func (b *BasePage) handleFilterInput(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "enter":
		b.filterInput.Blur()
		b.ApplyFilter()
		return nil, true
	case "esc":
		b.ClearFilter()
		return nil, true
	default:
		var cmd tea.Cmd
		b.filterInput, cmd = b.filterInput.Update(msg)
		b.ApplyFilter()
		return cmd, true
	}
}

// startFiltering initiates filter mode
func (b *BasePage) startFiltering() tea.Cmd {
	b.filterInput.SetValue("")
	b.filterInput.Focus()
	b.ApplyFilter()
	return nil
}

// View Rendering methods

func (b *BasePage) View() string {
	// Base implementation - specific pages should override this
	return b.table.View()
}

// State Management methods

func (b *BasePage) Reset() {
	b.isLoading = false
	b.error = nil
	b.ClearFilter()
}

func (b *BasePage) Refresh() tea.Cmd {
	// Base implementation - specific pages should override this
	b.SetLoading(true)
	b.SetError(nil)
	return nil
}

// Helper methods for common operations

// SetLoadingMessage sets a custom loading message
func (b *BasePage) SetLoadingMessage(message string) {
	b.loadingMessage = message
}

// GetLoadingMessage returns the loading message
func (b *BasePage) GetLoadingMessage() string {
	if b.loadingMessage == "" {
		return "Loading..."
	}
	return b.loadingMessage
}

// SetErrorMessage sets a custom error message
func (b *BasePage) SetErrorMessage(message string) {
	b.errorMessage = message
}

// GetErrorMessage returns the error message
func (b *BasePage) GetErrorMessage() string {
	if b.errorMessage == "" && b.error != nil {
		return b.error.Error()
	}
	return b.errorMessage
}

// HasData returns true if the page has data to display
func (b *BasePage) HasData() bool {
	// Base implementation - specific pages should override this
	return !b.isLoading && b.error == nil
}

// ShouldShowLoading returns true if loading state should be displayed
func (b *BasePage) ShouldShowLoading() bool {
	return b.isLoading && !b.HasData()
}

// ShouldShowError returns true if error state should be displayed
func (b *BasePage) ShouldShowError() bool {
	return b.error != nil && !b.HasData()
}
