package core

import (
	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// CoreInterface defines the clean API that the main UI model uses to interact with the core system
type CoreInterface interface {
	// Navigation
	GetCurrentMode() Mode
	GetNavigationState() NavigationState
	GoBack() tea.Cmd

	// State management
	GetStatusLine() string
	SetStatusLine(status string)

	// Context management
	IsShowingContextList() bool
	SetShowContextList(show bool)
	GetContextList() list.Model
	SetContextList(contextList list.Model)

	// Event handling
	HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool)
	UpdateTable(msg tea.KeyMsg) tea.Cmd
	IsAnyFilterActive() bool

	// View rendering
	View() string
	ViewWithHelpContext(helpContext layouts.HelpContext) string

	// Page access
	GetCurrentPage() interface{}

	// Layout management
	UpdateDimensions(width, height int)
	GetLayoutSystem() *layouts.LayoutSystem

	// Data access
	GetCurrentApp() models.ContainerApp
	GetCurrentRevision() models.Revision
	GetCurrentContainer() models.Container

	// Message handling
	HandleMessage(msg tea.Msg) tea.Cmd

	// Additional methods needed by main model
	LoadResourceGroups() tea.Cmd
	GetModeString() string
	GetError() error
	GetBreadcrumb() string
}

// Ensure CoreModel implements CoreInterface
var _ CoreInterface = (*CoreModel)(nil)

// HandleMessage handles various message types and delegates to appropriate handlers
func (cm *CoreModel) HandleMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case LoadedResourceGroupsMsg:
		return cm.handleLoadedResourceGroups(msg)
	case LoadedAppsMsg:
		return cm.handleLoadedApps(msg)
	case LoadedRevisionsMsg:
		return cm.handleLoadedRevisions(msg)
	case LoadedContainersMsg:
		return cm.handleLoadedContainers(msg)
	case RevisionRestartedMsg:
		return cm.handleRevisionRestarted(msg)
	case LeaveEnvVarsMsg:
		return cm.handleLeaveEnvVars(msg)
	default:
		return nil
	}
}

// Message handlers

func (cm *CoreModel) handleLoadedResourceGroups(msg LoadedResourceGroupsMsg) tea.Cmd {
	page := cm.pageManager.GetResourceGroupsPage()
	page.SetLoading(false)

	if msg.Error != nil {
		page.SetError(msg.Error)
		page.ClearData()
	} else {
		page.SetError(nil)
		page.SetData(msg.ResourceGroups)
	}

	return nil
}

func (cm *CoreModel) handleLoadedApps(msg LoadedAppsMsg) tea.Cmd {
	page := cm.pageManager.GetAppsPage()
	page.SetLoading(false)

	if msg.Error != nil {
		page.SetError(msg.Error)
		page.ClearData()
	} else {
		page.SetError(nil)
		page.SetData(msg.Apps)
	}

	return nil
}

func (cm *CoreModel) handleLoadedRevisions(msg LoadedRevisionsMsg) tea.Cmd {
	page := cm.pageManager.GetRevisionsPage()
	page.SetLoading(false)

	if msg.Error != nil {
		page.SetError(msg.Error)
		page.ClearData()
	} else {
		page.SetError(nil)
		page.SetData(msg.Revisions)
	}

	return nil
}

func (cm *CoreModel) handleLoadedContainers(msg LoadedContainersMsg) tea.Cmd {
	page := cm.pageManager.GetContainersPage()
	page.SetLoading(false)

	if msg.Error != nil {
		page.SetError(msg.Error)
		page.ClearData()
	} else {
		page.SetError(nil)
		// Cache containers
		cm.SetContainersCache(msg.AppID, msg.RevName, msg.Containers)
		page.SetData(msg.Containers)
	}

	return nil
}

func (cm *CoreModel) handleRevisionRestarted(msg RevisionRestartedMsg) tea.Cmd {
	if msg.Error != nil {
		cm.SetStatusLine("Restart failed: " + msg.Error.Error())
	} else {
		cm.SetStatusLine("Revision restart triggered.")
		// Reload revisions to reflect status changes after restart
		if app := cm.GetCurrentApp(); app.Name != "" {
			navState := cm.GetNavigationState()
			appID := cm.formatAppID(app)
			if appID == msg.AppID && navState.CurrentRevName == msg.RevName {
				return cm.LoadRevisions(app)
			}
		}
	}

	return nil
}

func (cm *CoreModel) handleLeaveEnvVars(msg LeaveEnvVarsMsg) tea.Cmd {
	return cm.GoBack()
}

// Convenience methods for common operations

// NavigateToModeByName navigates to a mode by its string name
func (cm *CoreModel) NavigateToModeByName(modeName string) tea.Cmd {
	switch modeName {
	case "resource-groups":
		return cm.NavigateToResourceGroups()
	default:
		return nil
	}
}

// GetBreadcrumb returns the current navigation breadcrumb
func (cm *CoreModel) GetBreadcrumb() string {
	return cm.navigationManager.GetBreadcrumb()
}

// CanGoBack returns whether back navigation is possible
func (cm *CoreModel) CanGoBack() bool {
	return cm.navigationManager.CanGoBack()
}

// GetModeString returns the current mode as a string
func (cm *CoreModel) GetModeString() string {
	return cm.GetCurrentMode().String()
}

// IsLoading returns whether the current page is loading
func (cm *CoreModel) IsLoading() bool {
	switch cm.GetCurrentMode() {
	case ModeResourceGroups:
		return cm.pageManager.GetResourceGroupsPage().IsLoading()
	case ModeApps:
		return cm.pageManager.GetAppsPage().IsLoading()
	case ModeRevisions:
		return cm.pageManager.GetRevisionsPage().IsLoading()
	case ModeContainers:
		return cm.pageManager.GetContainersPage().IsLoading()
	case ModeEnvVars:
		return cm.pageManager.GetEnvVarsPage().IsLoading()
	default:
		return false
	}
}

// GetError returns the current page's error
func (cm *CoreModel) GetError() error {
	switch cm.GetCurrentMode() {
	case ModeResourceGroups:
		return cm.pageManager.GetResourceGroupsPage().GetError()
	case ModeApps:
		return cm.pageManager.GetAppsPage().GetError()
	case ModeRevisions:
		return cm.pageManager.GetRevisionsPage().GetError()
	case ModeContainers:
		return cm.pageManager.GetContainersPage().GetError()
	case ModeEnvVars:
		return cm.pageManager.GetEnvVarsPage().GetError()
	default:
		return nil
	}
}

// RefreshCurrentPage refreshes the current page's data
func (cm *CoreModel) RefreshCurrentPage() tea.Cmd {
	switch cm.GetCurrentMode() {
	case ModeResourceGroups:
		return cm.LoadResourceGroups()
	case ModeApps:
		navState := cm.GetNavigationState()
		return cm.LoadApps(navState.CurrentRG)
	case ModeRevisions:
		if app := cm.GetCurrentApp(); app.Name != "" {
			return cm.LoadRevisions(app)
		}
	case ModeContainers:
		if app := cm.GetCurrentApp(); app.Name != "" {
			navState := cm.GetNavigationState()
			return cm.LoadContainers(app, navState.CurrentRevName)
		}
	}
	return nil
}

// GetCacheStats returns cache statistics
func (cm *CoreModel) GetCacheStats() map[string]int {
	return cm.stateManager.GetCacheStats()
}
