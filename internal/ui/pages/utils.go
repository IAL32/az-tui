package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Common key bindings that can be reused across pages
var (
	// Navigation keys
	EnterKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	)
	BackKey = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	)
	QuitKey = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)

	// Action keys
	RefreshKey = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	)
	FilterKey = key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	)
	HelpKey = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	)

	// Table navigation keys
	ScrollLeftKey = key.NewBinding(
		key.WithKeys("shift+left"),
		key.WithHelp("shift+←", "scroll left"),
	)
	ScrollRightKey = key.NewBinding(
		key.WithKeys("shift+right"),
		key.WithHelp("shift+→", "scroll right"),
	)

	// Action-specific keys
	LogsKey = key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "logs"),
	)
	ExecKey = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "exec"),
	)
	EnvVarsKey = key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "env vars"),
	)
	RestartKey = key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "restart"),
	)
)

// BaseKeyMap provides a comprehensive set of common key bindings
type BaseKeyMap struct {
	// Navigation
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding

	// Actions
	Refresh key.Binding
	Filter  key.Binding
	Help    key.Binding

	// Table navigation
	ScrollLeft  key.Binding
	ScrollRight key.Binding

	// Specific actions
	Logs    key.Binding
	Exec    key.Binding
	EnvVars key.Binding
	Restart key.Binding
}

// NewBaseKeyMap creates a new BaseKeyMap with default bindings
func NewBaseKeyMap() BaseKeyMap {
	return BaseKeyMap{
		Enter:       EnterKey,
		Back:        BackKey,
		Quit:        QuitKey,
		Refresh:     RefreshKey,
		Filter:      FilterKey,
		Help:        HelpKey,
		ScrollLeft:  ScrollLeftKey,
		ScrollRight: ScrollRightKey,
		Logs:        LogsKey,
		Exec:        ExecKey,
		EnvVars:     EnvVarsKey,
		Restart:     RestartKey,
	}
}

// GetCommonKeys returns the common key bindings used across all pages
func (b BaseKeyMap) GetCommonKeys() []key.Binding {
	return []key.Binding{
		b.Refresh,
		b.Filter,
		b.ScrollLeft,
		b.ScrollRight,
		b.Help,
		b.Back,
		b.Quit,
	}
}

// GetNavigationKeys returns key bindings for navigable pages
func (b BaseKeyMap) GetNavigationKeys() []key.Binding {
	return append([]key.Binding{b.Enter}, b.GetCommonKeys()...)
}

// GetActionKeys returns key bindings for actionable pages
func (b BaseKeyMap) GetActionKeys() []key.Binding {
	return []key.Binding{
		b.Logs,
		b.Exec,
		b.EnvVars,
		b.Restart,
	}
}

// Legacy functions for backward compatibility
func GetCommonKeys() []key.Binding {
	base := NewBaseKeyMap()
	return base.GetCommonKeys()
}

func GetNavigationKeys() []key.Binding {
	base := NewBaseKeyMap()
	return base.GetNavigationKeys()
}

func GetActionKeys() []key.Binding {
	base := NewBaseKeyMap()
	return base.GetActionKeys()
}

// KeyMatcher provides utility functions for matching key events
type KeyMatcher struct{}

// NewKeyMatcher creates a new KeyMatcher
func NewKeyMatcher() *KeyMatcher {
	return &KeyMatcher{}
}

// MatchesAny checks if a key message matches any of the provided key bindings
func (km *KeyMatcher) MatchesAny(msg tea.KeyMsg, bindings ...key.Binding) (key.Binding, bool) {
	msgKey := msg.String()
	for _, binding := range bindings {
		for _, bindingKey := range binding.Keys() {
			if msgKey == bindingKey {
				return binding, true
			}
		}
	}
	return key.Binding{}, false
}

// Matches checks if a key message matches a specific key binding
func (km *KeyMatcher) Matches(msg tea.KeyMsg, binding key.Binding) bool {
	msgKey := msg.String()
	for _, bindingKey := range binding.Keys() {
		if msgKey == bindingKey {
			return true
		}
	}
	return false
}

// PageState represents the current state of a page
type PageState int

const (
	StateLoading PageState = iota
	StateLoaded
	StateError
	StateEmpty
)

// String returns the string representation of the page state
func (ps PageState) String() string {
	switch ps {
	case StateLoading:
		return "loading"
	case StateLoaded:
		return "loaded"
	case StateError:
		return "error"
	case StateEmpty:
		return "empty"
	default:
		return "unknown"
	}
}

// PageStateManager helps manage page states
type PageStateManager struct {
	currentState PageState
}

// NewPageStateManager creates a new PageStateManager
func NewPageStateManager() *PageStateManager {
	return &PageStateManager{
		currentState: StateEmpty,
	}
}

// GetState returns the current page state
func (psm *PageStateManager) GetState() PageState {
	return psm.currentState
}

// SetState sets the current page state
func (psm *PageStateManager) SetState(state PageState) {
	psm.currentState = state
}

// IsLoading returns true if the page is in loading state
func (psm *PageStateManager) IsLoading() bool {
	return psm.currentState == StateLoading
}

// IsLoaded returns true if the page is in loaded state
func (psm *PageStateManager) IsLoaded() bool {
	return psm.currentState == StateLoaded
}

// HasError returns true if the page is in error state
func (psm *PageStateManager) HasError() bool {
	return psm.currentState == StateError
}

// IsEmpty returns true if the page is in empty state
func (psm *PageStateManager) IsEmpty() bool {
	return psm.currentState == StateEmpty
}

// UpdateState updates the state based on data and error conditions
func (psm *PageStateManager) UpdateState(hasData bool, isLoading bool, hasError bool) {
	if isLoading {
		psm.currentState = StateLoading
	} else if hasError {
		psm.currentState = StateError
	} else if hasData {
		psm.currentState = StateLoaded
	} else {
		psm.currentState = StateEmpty
	}
}

// Common error messages
const (
	ErrNoDataSelected         = "No data selected"
	ErrInvalidSelection       = "Invalid selection"
	ErrActionNotAvailable     = "Action not available"
	ErrNavigationNotSupported = "Navigation not supported"
)

// Common loading messages
const (
	LoadingDefault        = "Loading..."
	LoadingResourceGroups = "Loading resource groups..."
	LoadingApps           = "Loading container apps..."
	LoadingRevisions      = "Loading revisions..."
	LoadingContainers     = "Loading containers..."
	LoadingEnvVars        = "Loading environment variables..."
)

// Common filter placeholders
const (
	FilterResourceGroups = "Filter resource groups..."
	FilterApps           = "Filter apps..."
	FilterRevisions      = "Filter revisions..."
	FilterContainers     = "Filter containers..."
	FilterEnvVars        = "Filter environment variables..."
)

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	IsValid bool
	Message string
}

// NewValidationResult creates a new ValidationResult
func NewValidationResult(isValid bool, message string) ValidationResult {
	return ValidationResult{
		IsValid: isValid,
		Message: message,
	}
}

// Success creates a successful validation result
func Success() ValidationResult {
	return ValidationResult{IsValid: true}
}

// Error creates an error validation result
func Error(message string) ValidationResult {
	return ValidationResult{IsValid: false, Message: message}
}

// Validator provides common validation functions for pages
type Validator struct{}

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateSelection validates that an item is selected
func (v *Validator) ValidateSelection(hasSelection bool) ValidationResult {
	if !hasSelection {
		return Error(ErrNoDataSelected)
	}
	return Success()
}

// ValidateData validates that data is available
func (v *Validator) ValidateData(hasData bool) ValidationResult {
	if !hasData {
		return Error("No data available")
	}
	return Success()
}

// ValidateAction validates that an action is available
func (v *Validator) ValidateAction(actionName string, availableActions []string) ValidationResult {
	for _, action := range availableActions {
		if action == actionName {
			return Success()
		}
	}
	return Error(ErrActionNotAvailable)
}

// GetStatusColor returns a color for the given status string.
// This function handles all status types used across different pages:
// - Success states: Running, Succeeded, Healthy, Active
// - Error states: Failed, Error, Unhealthy
// - Warning states: Pending, Creating, Updating
// - Default: Gray for unknown states
func GetStatusColor(status string) lipgloss.Color {
	// Normalize status to lowercase for consistent matching
	normalizedStatus := strings.ToLower(status)

	switch normalizedStatus {
	case "running", "succeeded", "healthy", "active":
		return lipgloss.Color("#32CD32") // Success green
	case "failed", "error", "unhealthy":
		return lipgloss.Color("#FF6B6B") // Error red
	case "pending", "creating", "updating":
		return lipgloss.Color("#FFB347") // Warning yellow
	default:
		return lipgloss.Color("#808080") // Gray
	}
}
