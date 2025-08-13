package pages

import (
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// Page represents the common interface that all pages must implement.
// This interface defines the core contracts for data management, UI rendering,
// event handling, and state management across all pages in the UI system.
type Page interface {
	// Data Management
	IsLoading() bool
	SetLoading(loading bool)
	GetError() error
	SetError(err error)
	ClearData()

	// Table Management
	GetTable() table.Model
	SetTable(t table.Model)
	CreateTable() table.Model

	// Filter Management
	GetFilterInput() textinput.Model
	SetFilterInput(input textinput.Model)
	ApplyFilter()
	ClearFilter()

	// Event Handling
	HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool)
	HandleMessage(msg tea.Msg) tea.Cmd

	// View Rendering
	View() string
	ViewWithHelpContext(helpContext layouts.HelpContext) string

	// State Management
	Reset()
	Refresh() tea.Cmd
}

// TablePage extends the base Page interface for pages that display tables.
// This interface provides additional contracts specific to table-based pages
// including row selection, data operations, and table-specific navigation.
type TablePage[T any] interface {
	Page

	// Data Operations
	GetData() []T
	SetData(data []T)
	GetSelectedItem() (T, bool)
	GetSelectedIndex() int

	// Table-specific Operations
	UpdateTableWithData()
	GetHelpKeys() []key.Binding

	// Navigation
	CanNavigateToItem() bool
	NavigateToItem(item T) tea.Cmd
}

// DataLoader defines the interface for loading data asynchronously.
// Pages can implement this interface to provide custom data loading logic.
type DataLoader[T any] interface {
	LoadData() tea.Cmd
	HandleLoadedData(data []T, err error) tea.Cmd
}

// ActionHandler defines the interface for handling page-specific actions.
// Pages can implement this interface to provide custom action handling.
type ActionHandler[T any] interface {
	HandleAction(action string, item T) tea.Cmd
	GetAvailableActions(item T) []string
}

// Navigable defines the interface for pages that support navigation.
// This allows pages to define their navigation behavior and hierarchy.
type Navigable interface {
	CanGoBack() bool
	GoBack() tea.Cmd
	GetBreadcrumb() string
}
