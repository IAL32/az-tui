package pages

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// BaseTablePage provides a generic implementation of TablePage interface.
// It embeds BasePage and adds table-specific functionality with type safety.
type BaseTablePage[T any] struct {
	*BasePage

	// Data storage
	data []T

	// Table configuration
	createTableFunc func([]T) table.Model
	helpKeys        []key.Binding

	// Navigation
	canNavigate    bool
	navigationFunc func(T) tea.Cmd

	// Data loading
	loadDataFunc func() tea.Cmd
	refreshFunc  func() tea.Cmd
}

// NewBaseTablePage creates a new BaseTablePage with the given configuration.
func NewBaseTablePage[T any](filterPlaceholder string) *BaseTablePage[T] {
	return &BaseTablePage[T]{
		BasePage:    NewBasePage(filterPlaceholder),
		data:        make([]T, 0),
		canNavigate: false,
	}
}

// Configuration methods

// SetCreateTableFunc sets the function used to create the table with data
func (btp *BaseTablePage[T]) SetCreateTableFunc(fn func([]T) table.Model) {
	btp.createTableFunc = fn
}

// SetHelpKeys sets the help keys for this page
func (btp *BaseTablePage[T]) SetHelpKeys(keys []key.Binding) {
	btp.helpKeys = keys
}

// SetNavigationFunc enables navigation and sets the navigation function
func (btp *BaseTablePage[T]) SetNavigationFunc(fn func(T) tea.Cmd) {
	btp.navigationFunc = fn
	btp.canNavigate = true
}

// SetLoadDataFunc sets the function used to load data
func (btp *BaseTablePage[T]) SetLoadDataFunc(fn func() tea.Cmd) {
	btp.loadDataFunc = fn
}

// SetRefreshFunc sets the function used to refresh data
func (btp *BaseTablePage[T]) SetRefreshFunc(fn func() tea.Cmd) {
	btp.refreshFunc = fn
}

// Data Operations (implementing TablePage interface)

func (btp *BaseTablePage[T]) GetData() []T {
	return btp.data
}

func (btp *BaseTablePage[T]) SetData(data []T) {
	btp.data = data
	btp.UpdateTableWithData()
}

func (btp *BaseTablePage[T]) GetSelectedItem() (T, bool) {
	var zero T

	if len(btp.data) == 0 {
		return zero, false
	}

	selectedRow := btp.table.HighlightedRow()
	if selectedRow.Data == nil {
		return zero, false
	}

	// Get the selected index from the table
	selectedIndex := btp.GetSelectedIndex()
	if selectedIndex < 0 || selectedIndex >= len(btp.data) {
		return zero, false
	}

	return btp.data[selectedIndex], true
}

func (btp *BaseTablePage[T]) GetSelectedIndex() int {
	if len(btp.data) == 0 {
		return -1
	}

	// Get the currently highlighted row index
	// Note: This is a simplified implementation - in practice, you might need
	// to track the selected index more carefully depending on filtering/sorting
	selectedRow := btp.table.HighlightedRow()
	if selectedRow.Data == nil {
		return -1
	}

	// For now, return 0 as a placeholder - this would need to be implemented
	// based on the specific table implementation and how row selection works
	return 0
}

// Table-specific Operations

func (btp *BaseTablePage[T]) UpdateTableWithData() {
	if btp.createTableFunc != nil {
		btp.table = btp.createTableFunc(btp.data)
		// Reapply filter if it was active
		if btp.filterInput.Value() != "" {
			btp.ApplyFilter()
		}
	}
}

func (btp *BaseTablePage[T]) GetHelpKeys() []key.Binding {
	return btp.helpKeys
}

// Navigation

func (btp *BaseTablePage[T]) CanNavigateToItem() bool {
	return btp.canNavigate && btp.navigationFunc != nil
}

func (btp *BaseTablePage[T]) NavigateToItem(item T) tea.Cmd {
	if btp.navigationFunc != nil {
		return btp.navigationFunc(item)
	}
	return nil
}

// Override base methods with table-specific behavior

func (btp *BaseTablePage[T]) ClearData() {
	btp.BasePage.ClearData()
	btp.data = make([]T, 0)
	btp.UpdateTableWithData()
}

func (btp *BaseTablePage[T]) CreateTable() table.Model {
	if btp.createTableFunc != nil {
		return btp.createTableFunc(btp.data)
	}
	return btp.BasePage.CreateTable()
}

func (btp *BaseTablePage[T]) HasData() bool {
	return len(btp.data) > 0 && btp.BasePage.HasData()
}

func (btp *BaseTablePage[T]) Refresh() tea.Cmd {
	btp.BasePage.Refresh()
	if btp.refreshFunc != nil {
		return btp.refreshFunc()
	}
	if btp.loadDataFunc != nil {
		return btp.loadDataFunc()
	}
	return nil
}

// Enhanced key handling for table pages

func (btp *BaseTablePage[T]) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// First, try base page key handling
	if cmd, handled := btp.BasePage.HandleKeyMsg(msg); handled {
		return cmd, handled
	}

	// Handle table-specific keys
	switch msg.String() {
	case "enter":
		if btp.CanNavigateToItem() {
			if item, ok := btp.GetSelectedItem(); ok {
				return btp.NavigateToItem(item), true
			}
		}
		return nil, true

	case "r":
		return btp.Refresh(), true
	}

	return nil, false
}

// Helper methods for common table operations

// FindItemByPredicate finds an item in the data using a predicate function
func (btp *BaseTablePage[T]) FindItemByPredicate(predicate func(T) bool) (T, bool) {
	var zero T
	for _, item := range btp.data {
		if predicate(item) {
			return item, true
		}
	}
	return zero, false
}

// FilterData returns filtered data based on a predicate
func (btp *BaseTablePage[T]) FilterData(predicate func(T) bool) []T {
	filtered := make([]T, 0)
	for _, item := range btp.data {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// GetDataCount returns the number of items in the data
func (btp *BaseTablePage[T]) GetDataCount() int {
	return len(btp.data)
}

// IsEmpty returns true if the page has no data
func (btp *BaseTablePage[T]) IsEmpty() bool {
	return len(btp.data) == 0
}
