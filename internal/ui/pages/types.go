package pages

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// NavigablePage extends TablePage with navigation capabilities.
// This is useful for pages that are part of a hierarchical navigation structure.
type NavigablePage[T any] struct {
	*BaseTablePage[T]

	// Navigation state
	breadcrumb string
	canGoBack  bool
	backFunc   func() tea.Cmd
	parentPage Page
}

// NewNavigablePage creates a new NavigablePage with navigation capabilities.
func NewNavigablePage[T any](filterPlaceholder string) *NavigablePage[T] {
	return &NavigablePage[T]{
		BaseTablePage: NewBaseTablePage[T](filterPlaceholder),
		canGoBack:     false,
	}
}

// Navigation methods (implementing Navigable interface)

func (np *NavigablePage[T]) CanGoBack() bool {
	return np.canGoBack && np.backFunc != nil
}

func (np *NavigablePage[T]) GoBack() tea.Cmd {
	if np.backFunc != nil {
		return np.backFunc()
	}
	return nil
}

func (np *NavigablePage[T]) GetBreadcrumb() string {
	return np.breadcrumb
}

// Configuration methods

func (np *NavigablePage[T]) SetBreadcrumb(breadcrumb string) {
	np.breadcrumb = breadcrumb
}

func (np *NavigablePage[T]) SetBackFunc(fn func() tea.Cmd) {
	np.backFunc = fn
	np.canGoBack = true
}

func (np *NavigablePage[T]) SetParentPage(parent Page) {
	np.parentPage = parent
}

// Enhanced key handling with back navigation

func (np *NavigablePage[T]) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// First, try base table page key handling
	if cmd, handled := np.BaseTablePage.HandleKeyMsg(msg); handled {
		return cmd, handled
	}

	// Handle navigation-specific keys
	switch msg.String() {
	case "esc":
		if np.CanGoBack() {
			return np.GoBack(), true
		}
		return nil, true
	}

	return nil, false
}

// ActionablePage extends NavigablePage with action handling capabilities.
// This is useful for pages that support multiple actions on items.
type ActionablePage[T any] struct {
	*NavigablePage[T]

	// Action handling
	actions    map[string]func(T) tea.Cmd
	actionKeys map[string]key.Binding
}

// NewActionablePage creates a new ActionablePage with action capabilities.
func NewActionablePage[T any](filterPlaceholder string) *ActionablePage[T] {
	return &ActionablePage[T]{
		NavigablePage: NewNavigablePage[T](filterPlaceholder),
		actions:       make(map[string]func(T) tea.Cmd),
		actionKeys:    make(map[string]key.Binding),
	}
}

// Action methods (implementing ActionHandler interface)

func (ap *ActionablePage[T]) HandleAction(action string, item T) tea.Cmd {
	if actionFunc, exists := ap.actions[action]; exists {
		return actionFunc(item)
	}
	return nil
}

func (ap *ActionablePage[T]) GetAvailableActions(item T) []string {
	actions := make([]string, 0, len(ap.actions))
	for action := range ap.actions {
		actions = append(actions, action)
	}
	return actions
}

// Configuration methods

func (ap *ActionablePage[T]) AddAction(name string, keyBinding key.Binding, actionFunc func(T) tea.Cmd) {
	ap.actions[name] = actionFunc
	ap.actionKeys[name] = keyBinding
}

func (ap *ActionablePage[T]) RemoveAction(name string) {
	delete(ap.actions, name)
	delete(ap.actionKeys, name)
}

func (ap *ActionablePage[T]) GetActionKeys() []key.Binding {
	// Get action names and sort them for consistent ordering
	actionNames := make([]string, 0, len(ap.actionKeys))
	for actionName := range ap.actionKeys {
		actionNames = append(actionNames, actionName)
	}

	// Sort action names to ensure consistent order
	for i := 0; i < len(actionNames); i++ {
		for j := i + 1; j < len(actionNames); j++ {
			if actionNames[i] > actionNames[j] {
				actionNames[i], actionNames[j] = actionNames[j], actionNames[i]
			}
		}
	}

	// Build keys in sorted order
	keys := make([]key.Binding, 0, len(ap.actionKeys))
	for _, actionName := range actionNames {
		keys = append(keys, ap.actionKeys[actionName])
	}
	return keys
}

// Enhanced key handling with actions

func (ap *ActionablePage[T]) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// First, try base navigable page key handling
	if cmd, handled := ap.NavigablePage.HandleKeyMsg(msg); handled {
		return cmd, handled
	}

	// Handle action keys by comparing key strings
	msgKey := msg.String()
	for actionName, keyBinding := range ap.actionKeys {
		// Check if the key matches any of the binding's keys
		for _, bindingKey := range keyBinding.Keys() {
			if msgKey == bindingKey {
				if item, ok := ap.GetSelectedItem(); ok {
					return ap.HandleAction(actionName, item), true
				}
				return nil, true
			}
		}
	}

	return nil, false
}

// ReadOnlyPage is a specialized page for displaying read-only data.
// It extends BaseTablePage but removes navigation and actions.
type ReadOnlyPage[T any] struct {
	*BaseTablePage[T]
}

// NewReadOnlyPage creates a new ReadOnlyPage for read-only data display.
func NewReadOnlyPage[T any](filterPlaceholder string) *ReadOnlyPage[T] {
	page := &ReadOnlyPage[T]{
		BaseTablePage: NewBaseTablePage[T](filterPlaceholder),
	}
	// Disable navigation for read-only pages
	page.canNavigate = false
	return page
}

// Override navigation methods to disable them

func (rop *ReadOnlyPage[T]) CanNavigateToItem() bool {
	return false
}

func (rop *ReadOnlyPage[T]) NavigateToItem(item T) tea.Cmd {
	return nil
}

// Enhanced key handling for read-only pages (no enter navigation)

func (rop *ReadOnlyPage[T]) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	// Handle filter input when focused
	if rop.filterInput.Focused() {
		if cmd, handled := rop.BasePage.HandleKeyMsg(msg); handled {
			return cmd, handled
		}
	}

	// Handle common keys but skip enter navigation
	switch msg.String() {
	case "/":
		return rop.startFiltering(), true
	case "r":
		return rop.Refresh(), true
	case "ctrl+c", "q":
		return tea.Quit, true
	case "esc":
		// For read-only pages, esc just clears filter or does nothing
		if rop.filterInput.Focused() {
			rop.ClearFilter()
			return nil, true
		}
		return nil, true
	}

	return nil, false
}

// Helper functions for creating configured pages

// CreateConfiguredTablePage creates a basic table page with table creation function.
func CreateConfiguredTablePage[T any](filterPlaceholder string, createTableFunc func([]T) table.Model) *BaseTablePage[T] {
	page := NewBaseTablePage[T](filterPlaceholder)
	page.SetCreateTableFunc(createTableFunc)
	return page
}

// CreateConfiguredNavigablePage creates a navigable table page with table creation function.
func CreateConfiguredNavigablePage[T any](filterPlaceholder string, createTableFunc func([]T) table.Model) *NavigablePage[T] {
	page := NewNavigablePage[T](filterPlaceholder)
	page.SetCreateTableFunc(createTableFunc)
	return page
}

// CreateConfiguredActionablePage creates an actionable table page with table creation function.
func CreateConfiguredActionablePage[T any](filterPlaceholder string, createTableFunc func([]T) table.Model) *ActionablePage[T] {
	page := NewActionablePage[T](filterPlaceholder)
	page.SetCreateTableFunc(createTableFunc)
	return page
}

// CreateConfiguredReadOnlyPage creates a read-only table page with table creation function.
func CreateConfiguredReadOnlyPage[T any](filterPlaceholder string, createTableFunc func([]T) table.Model) *ReadOnlyPage[T] {
	page := NewReadOnlyPage[T](filterPlaceholder)
	page.SetCreateTableFunc(createTableFunc)
	return page
}
