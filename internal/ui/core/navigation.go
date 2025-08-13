package core

import (
	"github.com/IAL32/az-tui/internal/models"
)

// NavigationManager handles navigation state and flow management
type NavigationManager struct {
	state       NavigationState
	currentMode Mode
	history     []NavigationStep
}

// NavigationStep represents a step in navigation history
type NavigationStep struct {
	Mode  Mode
	State NavigationState
}

// NewNavigationManager creates a new navigation manager
func NewNavigationManager() *NavigationManager {
	return &NavigationManager{
		state:       NavigationState{},
		currentMode: ModeResourceGroups,
		history:     make([]NavigationStep, 0),
	}
}

// GetCurrentMode returns the current mode
func (nm *NavigationManager) GetCurrentMode() Mode {
	return nm.currentMode
}

// GetNavigationState returns the current navigation state
func (nm *NavigationManager) GetNavigationState() NavigationState {
	return nm.state
}

// SetCurrentRG sets the current resource group
func (nm *NavigationManager) SetCurrentRG(rg string) {
	nm.state.CurrentRG = rg
}

// SetCurrentAppID sets the current app ID
func (nm *NavigationManager) SetCurrentAppID(appID string) {
	nm.state.CurrentAppID = appID
}

// SetCurrentRevName sets the current revision name
func (nm *NavigationManager) SetCurrentRevName(revName string) {
	nm.state.CurrentRevName = revName
}

// SetCurrentContainerName sets the current container name
func (nm *NavigationManager) SetCurrentContainerName(containerName string) {
	nm.state.CurrentContainerName = containerName
}

// NavigateToMode navigates to a specific mode and updates state accordingly
func (nm *NavigationManager) NavigateToMode(mode Mode) {
	// Save current step to history
	nm.pushToHistory()

	// Reset navigation state from the target mode level
	nm.state.ResetFrom(mode)

	// Update current mode
	nm.currentMode = mode
}

// NavigateToResourceGroups navigates to resource groups mode
func (nm *NavigationManager) NavigateToResourceGroups() {
	nm.NavigateToMode(ModeResourceGroups)
}

// NavigateToApps navigates to apps mode with resource group context
func (nm *NavigationManager) NavigateToApps(rg models.ResourceGroup) {
	nm.pushToHistory()
	nm.currentMode = ModeApps
	nm.state.CurrentRG = rg.Name
	nm.state.ResetFrom(ModeApps)
}

// NavigateToRevisions navigates to revisions mode with app context
func (nm *NavigationManager) NavigateToRevisions(app models.ContainerApp) {
	nm.pushToHistory()
	nm.currentMode = ModeRevisions
	nm.state.CurrentAppID = nm.formatAppID(app)
	nm.state.ResetFrom(ModeRevisions)
}

// NavigateToContainers navigates to containers mode with revision context
func (nm *NavigationManager) NavigateToContainers(rev models.Revision) {
	nm.pushToHistory()
	nm.currentMode = ModeContainers
	nm.state.CurrentRevName = rev.Name
	nm.state.ResetFrom(ModeContainers)
}

// NavigateToEnvVars navigates to environment variables mode with container context
func (nm *NavigationManager) NavigateToEnvVars(container models.Container) {
	nm.pushToHistory()
	nm.currentMode = ModeEnvVars
	nm.state.CurrentContainerName = container.Name
}

// GoBack navigates back to the previous mode
func (nm *NavigationManager) GoBack() bool {
	if len(nm.history) == 0 {
		return false
	}

	// Pop the last step from history
	lastStep := nm.history[len(nm.history)-1]
	nm.history = nm.history[:len(nm.history)-1]

	// Restore the previous state
	nm.currentMode = lastStep.Mode
	nm.state = lastStep.State

	return true
}

// CanGoBack returns true if there's navigation history to go back to
func (nm *NavigationManager) CanGoBack() bool {
	return len(nm.history) > 0
}

// GetBreadcrumb returns the current breadcrumb
func (nm *NavigationManager) GetBreadcrumb() string {
	return nm.state.GetBreadcrumb()
}

// GetParentMode returns the parent mode for the current mode
func (nm *NavigationManager) GetParentMode() (Mode, bool) {
	switch nm.currentMode {
	case ModeApps:
		return ModeResourceGroups, true
	case ModeRevisions:
		return ModeApps, true
	case ModeContainers:
		return ModeRevisions, true
	case ModeEnvVars:
		return ModeContainers, true
	default:
		return ModeResourceGroups, false
	}
}

// ValidateNavigation validates if navigation to a mode is possible
func (nm *NavigationManager) ValidateNavigation(mode Mode) bool {
	switch mode {
	case ModeResourceGroups:
		return true // Always can go to resource groups
	case ModeApps:
		return nm.state.CurrentRG != "" // Need resource group
	case ModeRevisions:
		return nm.state.CurrentRG != "" && nm.state.CurrentAppID != "" // Need RG and app
	case ModeContainers:
		return nm.state.CurrentRG != "" && nm.state.CurrentAppID != "" && nm.state.CurrentRevName != "" // Need RG, app, and revision
	case ModeEnvVars:
		return nm.state.CurrentRG != "" && nm.state.CurrentAppID != "" && nm.state.CurrentRevName != "" && nm.state.CurrentContainerName != "" // Need all
	default:
		return false
	}
}

// GetNavigationFlow returns the expected navigation flow for the current state
func (nm *NavigationManager) GetNavigationFlow() []Mode {
	flow := []Mode{ModeResourceGroups}

	if nm.state.CurrentRG != "" {
		flow = append(flow, ModeApps)
	}
	if nm.state.CurrentAppID != "" {
		flow = append(flow, ModeRevisions)
	}
	if nm.state.CurrentRevName != "" {
		flow = append(flow, ModeContainers)
	}
	if nm.state.CurrentContainerName != "" {
		flow = append(flow, ModeEnvVars)
	}

	return flow
}

// Reset resets the navigation manager to initial state
func (nm *NavigationManager) Reset() {
	nm.state.Reset()
	nm.currentMode = ModeResourceGroups
	nm.history = nm.history[:0] // Clear history
}

// pushToHistory saves the current state to history
func (nm *NavigationManager) pushToHistory() {
	step := NavigationStep{
		Mode:  nm.currentMode,
		State: nm.state,
	}
	nm.history = append(nm.history, step)

	// Limit history size to prevent memory issues
	const maxHistorySize = 50
	if len(nm.history) > maxHistorySize {
		nm.history = nm.history[1:]
	}
}

// formatAppID formats an app into an ID string
func (nm *NavigationManager) formatAppID(app models.ContainerApp) string {
	return app.ResourceGroup + "/" + app.Name
}

// CreateNavigationEvent creates a navigation event
func (nm *NavigationManager) CreateNavigationEvent(fromMode, toMode Mode, data interface{}) NavigationEvent {
	return NavigationEvent{
		FromMode: fromMode,
		ToMode:   toMode,
		Data:     data,
	}
}

// CreateModeChangeEvent creates a mode change event
func (nm *NavigationManager) CreateModeChangeEvent(oldMode, newMode Mode) ModeChangeEvent {
	return ModeChangeEvent{
		OldMode: oldMode,
		NewMode: newMode,
	}
}
