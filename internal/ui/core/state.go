package core

import (
	"github.com/IAL32/az-tui/internal/models"
)

// StateManager manages global state and shared data between pages
type StateManager struct {
	// Performance cache (shared across pages)
	containersByRev map[string][]models.Container // key: revKey(appID, revName)

	// Current data cache for quick access
	currentApp       *models.ContainerApp
	currentRevision  *models.Revision
	currentContainer *models.Container

	// Global status
	statusLine string

	// Context management
	showContextList bool
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		containersByRev: make(map[string][]models.Container),
		showContextList: false,
	}
}

// Container cache management

// SetContainersForRevision caches containers for a specific revision
func (sm *StateManager) SetContainersForRevision(appID, revName string, containers []models.Container) {
	key := sm.revKey(appID, revName)
	sm.containersByRev[key] = containers
}

// GetContainersForRevision retrieves cached containers for a specific revision
func (sm *StateManager) GetContainersForRevision(appID, revName string) ([]models.Container, bool) {
	key := sm.revKey(appID, revName)
	containers, exists := sm.containersByRev[key]
	return containers, exists
}

// ClearContainerCache clears the container cache
func (sm *StateManager) ClearContainerCache() {
	sm.containersByRev = make(map[string][]models.Container)
}

// Current data management

// SetCurrentApp sets the current app
func (sm *StateManager) SetCurrentApp(app models.ContainerApp) {
	sm.currentApp = &app
}

// GetCurrentApp returns the current app
func (sm *StateManager) GetCurrentApp() (models.ContainerApp, bool) {
	if sm.currentApp == nil {
		return models.ContainerApp{}, false
	}
	return *sm.currentApp, true
}

// ClearCurrentApp clears the current app
func (sm *StateManager) ClearCurrentApp() {
	sm.currentApp = nil
}

// SetCurrentRevision sets the current revision
func (sm *StateManager) SetCurrentRevision(revision models.Revision) {
	sm.currentRevision = &revision
}

// GetCurrentRevision returns the current revision
func (sm *StateManager) GetCurrentRevision() (models.Revision, bool) {
	if sm.currentRevision == nil {
		return models.Revision{}, false
	}
	return *sm.currentRevision, true
}

// ClearCurrentRevision clears the current revision
func (sm *StateManager) ClearCurrentRevision() {
	sm.currentRevision = nil
}

// SetCurrentContainer sets the current container
func (sm *StateManager) SetCurrentContainer(container models.Container) {
	sm.currentContainer = &container
}

// GetCurrentContainer returns the current container
func (sm *StateManager) GetCurrentContainer() (models.Container, bool) {
	if sm.currentContainer == nil {
		return models.Container{}, false
	}
	return *sm.currentContainer, true
}

// ClearCurrentContainer clears the current container
func (sm *StateManager) ClearCurrentContainer() {
	sm.currentContainer = nil
}

// Status management

// SetStatusLine sets the global status line
func (sm *StateManager) SetStatusLine(status string) {
	sm.statusLine = status
}

// GetStatusLine returns the global status line
func (sm *StateManager) GetStatusLine() string {
	return sm.statusLine
}

// ClearStatusLine clears the global status line
func (sm *StateManager) ClearStatusLine() {
	sm.statusLine = ""
}

// Context list management

// SetShowContextList sets whether the context list should be shown
func (sm *StateManager) SetShowContextList(show bool) {
	sm.showContextList = show
}

// IsShowingContextList returns whether the context list is being shown
func (sm *StateManager) IsShowingContextList() bool {
	return sm.showContextList
}

// State validation and cleanup

// ValidateState ensures the state is consistent with the current navigation
func (sm *StateManager) ValidateState(navState NavigationState) {
	// Clear current app if we're not in app-related modes
	if navState.CurrentAppID == "" {
		sm.ClearCurrentApp()
	}

	// Clear current revision if we're not in revision-related modes
	if navState.CurrentRevName == "" {
		sm.ClearCurrentRevision()
	}

	// Clear current container if we're not in container-related modes
	if navState.CurrentContainerName == "" {
		sm.ClearCurrentContainer()
	}
}

// Reset resets all state to initial values
func (sm *StateManager) Reset() {
	sm.ClearContainerCache()
	sm.ClearCurrentApp()
	sm.ClearCurrentRevision()
	sm.ClearCurrentContainer()
	sm.ClearStatusLine()
	sm.showContextList = false
}

// Persistence and restoration

// StateSnapshot represents a snapshot of the current state
type StateSnapshot struct {
	ContainersByRev  map[string][]models.Container
	CurrentApp       *models.ContainerApp
	CurrentRevision  *models.Revision
	CurrentContainer *models.Container
	StatusLine       string
	ShowContextList  bool
}

// CreateSnapshot creates a snapshot of the current state
func (sm *StateManager) CreateSnapshot() StateSnapshot {
	// Deep copy the containers cache
	containersCopy := make(map[string][]models.Container)
	for k, v := range sm.containersByRev {
		containersCopy[k] = make([]models.Container, len(v))
		copy(containersCopy[k], v)
	}

	return StateSnapshot{
		ContainersByRev:  containersCopy,
		CurrentApp:       sm.currentApp,
		CurrentRevision:  sm.currentRevision,
		CurrentContainer: sm.currentContainer,
		StatusLine:       sm.statusLine,
		ShowContextList:  sm.showContextList,
	}
}

// RestoreSnapshot restores state from a snapshot
func (sm *StateManager) RestoreSnapshot(snapshot StateSnapshot) {
	sm.containersByRev = snapshot.ContainersByRev
	sm.currentApp = snapshot.CurrentApp
	sm.currentRevision = snapshot.CurrentRevision
	sm.currentContainer = snapshot.CurrentContainer
	sm.statusLine = snapshot.StatusLine
	sm.showContextList = snapshot.ShowContextList
}

// Helper methods

// revKey creates a key for the containers cache
func (sm *StateManager) revKey(appID, revName string) string {
	return appID + "@" + revName
}

// formatAppID formats an app into an ID string
func (sm *StateManager) formatAppID(app models.ContainerApp) string {
	return app.ResourceGroup + "/" + app.Name
}

// GetCacheStats returns statistics about the cache
func (sm *StateManager) GetCacheStats() map[string]int {
	stats := make(map[string]int)
	stats["containers_cache_entries"] = len(sm.containersByRev)

	totalContainers := 0
	for _, containers := range sm.containersByRev {
		totalContainers += len(containers)
	}
	stats["total_cached_containers"] = totalContainers

	return stats
}
