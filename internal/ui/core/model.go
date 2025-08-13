package core

import (
	"fmt"
	"os"

	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/providers"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// CoreModel is the central coordination model that manages all pages, navigation, and shared state
type CoreModel struct {
	// Core managers
	navigationManager *NavigationManager
	pageManager       *PageManager
	stateManager      *StateManager

	// Layout system
	layoutSystem *layouts.LayoutSystem

	// Providers
	dataProvider    providers.DataProvider
	commandProvider providers.CommandProvider

	// Context list for mode switching
	contextList list.Model

	// Terminal dimensions
	termW, termH int
}

// NewCoreModel creates a new core model
func NewCoreModel(dataProvider providers.DataProvider, commandProvider providers.CommandProvider, termW, termH int) *CoreModel {
	// Create managers
	navigationManager := NewNavigationManager()
	stateManager := NewStateManager()
	layoutSystem := layouts.NewLayoutSystem(termW, termH)
	pageManager := NewPageManager(layoutSystem, navigationManager)

	// Create core model
	coreModel := &CoreModel{
		navigationManager: navigationManager,
		pageManager:       pageManager,
		stateManager:      stateManager,
		layoutSystem:      layoutSystem,
		dataProvider:      dataProvider,
		commandProvider:   commandProvider,
		termW:             termW,
		termH:             termH,
	}

	// Setup page navigation and actions
	pageManager.SetupPageNavigation(coreModel)
	pageManager.SetupPageActions(coreModel)

	// Initialize with environment variable if available
	if rg := os.Getenv("ACA_RG"); rg != "" {
		navigationManager.SetCurrentRG(rg)
	}

	return coreModel
}

// Navigation methods

// NavigateToResourceGroups navigates to resource groups mode
func (cm *CoreModel) NavigateToResourceGroups() tea.Cmd {
	cm.navigationManager.NavigateToResourceGroups()
	cm.stateManager.ValidateState(cm.navigationManager.GetNavigationState())

	// Set up the resource groups page
	page := cm.pageManager.GetResourceGroupsPage()
	page.SetLoading(true)
	page.SetError(nil)
	page.ClearData()

	return cm.LoadResourceGroups()
}

// NavigateToApps navigates to apps mode with resource group context
func (cm *CoreModel) NavigateToApps(rg models.ResourceGroup) tea.Cmd {
	cm.navigationManager.NavigateToApps(rg)
	cm.stateManager.ValidateState(cm.navigationManager.GetNavigationState())

	// Set up the apps page
	page := cm.pageManager.GetAppsPage()
	page.SetResourceGroupContext(rg.Name)
	page.SetLoading(true)
	page.SetError(nil)
	page.ClearData()

	return cm.LoadApps(rg.Name)
}

// NavigateToRevisions navigates to revisions mode with app context
func (cm *CoreModel) NavigateToRevisions(app models.ContainerApp) tea.Cmd {
	cm.navigationManager.NavigateToRevisions(app)
	cm.stateManager.SetCurrentApp(app)
	cm.stateManager.ValidateState(cm.navigationManager.GetNavigationState())

	// Set up the revisions page
	page := cm.pageManager.GetRevisionsPage()
	appID := cm.formatAppID(app)
	page.SetAppContext(app.Name, appID)
	page.SetLoading(true)
	page.SetError(nil)
	page.ClearData()

	return cm.LoadRevisions(app)
}

// NavigateToContainers navigates to containers mode with revision context
func (cm *CoreModel) NavigateToContainers(rev models.Revision) tea.Cmd {
	cm.navigationManager.NavigateToContainers(rev)
	cm.stateManager.SetCurrentRevision(rev)
	cm.stateManager.ValidateState(cm.navigationManager.GetNavigationState())

	// Set up the containers page
	page := cm.pageManager.GetContainersPage()
	app := cm.GetCurrentApp()
	appID := cm.formatAppID(app)
	page.SetRevisionContext(app.Name, appID, rev.Name)
	page.SetLoading(true)
	page.SetError(nil)
	page.ClearData()

	return cm.LoadContainers(app, rev.Name)
}

// NavigateToEnvVars navigates to environment variables mode with container context
func (cm *CoreModel) NavigateToEnvVars(container models.Container) tea.Cmd {
	cm.navigationManager.NavigateToEnvVars(container)
	cm.stateManager.SetCurrentContainer(container)
	cm.stateManager.ValidateState(cm.navigationManager.GetNavigationState())

	// Set up the env vars page
	page := cm.pageManager.GetEnvVarsPage()
	containers := cm.pageManager.GetContainersPage().GetData()
	page.SetContainerContext(container.Name, containers)

	return nil
}

// GoBack navigates back to the previous mode
func (cm *CoreModel) GoBack() tea.Cmd {
	if !cm.navigationManager.GoBack() {
		return nil
	}

	cm.stateManager.ValidateState(cm.navigationManager.GetNavigationState())

	// Refresh the current page if needed
	switch cm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		return cm.LoadResourceGroups()
	case ModeApps:
		navState := cm.navigationManager.GetNavigationState()
		return cm.LoadApps(navState.CurrentRG)
	case ModeRevisions:
		if app, ok := cm.stateManager.GetCurrentApp(); ok {
			return cm.LoadRevisions(app)
		}
	case ModeContainers:
		if app, ok := cm.stateManager.GetCurrentApp(); ok {
			navState := cm.navigationManager.GetNavigationState()
			return cm.LoadContainers(app, navState.CurrentRevName)
		}
	}

	return nil
}

// Data loading methods

// LoadResourceGroups loads resource groups data
func (cm *CoreModel) LoadResourceGroups() tea.Cmd {
	return CreateLoadResourceGroupsCmd(cm.dataProvider)
}

// LoadApps loads apps data for a resource group
func (cm *CoreModel) LoadApps(resourceGroup string) tea.Cmd {
	return CreateLoadAppsCmd(cm.dataProvider, resourceGroup)
}

// LoadRevisions loads revisions data for an app
func (cm *CoreModel) LoadRevisions(app models.ContainerApp) tea.Cmd {
	return CreateLoadRevisionsCmd(cm.dataProvider, app)
}

// LoadContainers loads containers data for a revision
func (cm *CoreModel) LoadContainers(app models.ContainerApp, revName string) tea.Cmd {
	return CreateLoadContainersCmd(cm.dataProvider, app, revName)
}

// Action methods

// ShowAppLogs shows logs for an app
func (cm *CoreModel) ShowAppLogs(app models.ContainerApp) tea.Cmd {
	return cm.commandProvider.ShowAppLogs(app)
}

// ExecIntoApp executes into an app
func (cm *CoreModel) ExecIntoApp(app models.ContainerApp) tea.Cmd {
	return cm.commandProvider.ExecIntoApp(app)
}

// RestartRevision restarts a revision
func (cm *CoreModel) RestartRevision(rev models.Revision) tea.Cmd {
	app := cm.GetCurrentApp()
	return cm.commandProvider.RestartRevision(app, rev.Name)
}

// ShowRevisionLogs shows logs for a revision
func (cm *CoreModel) ShowRevisionLogs(rev models.Revision) tea.Cmd {
	app := cm.GetCurrentApp()
	return cm.commandProvider.ShowRevisionLogs(app, rev.Name)
}

// ExecIntoRevision executes into a revision
func (cm *CoreModel) ExecIntoRevision(rev models.Revision) tea.Cmd {
	app := cm.GetCurrentApp()
	return cm.commandProvider.ExecIntoRevision(app, rev.Name)
}

// ShowContainerLogs shows logs for a container
func (cm *CoreModel) ShowContainerLogs(container models.Container) tea.Cmd {
	app := cm.GetCurrentApp()
	navState := cm.navigationManager.GetNavigationState()
	return cm.commandProvider.ShowContainerLogs(app, navState.CurrentRevName, container.Name)
}

// ExecIntoContainer executes into a container
func (cm *CoreModel) ExecIntoContainer(container models.Container) tea.Cmd {
	app := cm.GetCurrentApp()
	navState := cm.navigationManager.GetNavigationState()
	return cm.commandProvider.ExecIntoContainer(app, navState.CurrentRevName, container.Name)
}

// State access methods

// GetCurrentMode returns the current mode
func (cm *CoreModel) GetCurrentMode() Mode {
	return cm.navigationManager.GetCurrentMode()
}

// GetNavigationState returns the current navigation state
func (cm *CoreModel) GetNavigationState() NavigationState {
	return cm.navigationManager.GetNavigationState()
}

// GetCurrentApp returns the current app
func (cm *CoreModel) GetCurrentApp() models.ContainerApp {
	if app, ok := cm.stateManager.GetCurrentApp(); ok {
		return app
	}
	// Fallback: try to get from apps page
	if app, ok := cm.pageManager.GetAppsPage().GetSelectedItem(); ok {
		return app
	}
	return models.ContainerApp{}
}

// GetCurrentRevision returns the current revision
func (cm *CoreModel) GetCurrentRevision() models.Revision {
	if rev, ok := cm.stateManager.GetCurrentRevision(); ok {
		return rev
	}
	// Fallback: try to get from revisions page
	if rev, ok := cm.pageManager.GetRevisionsPage().GetSelectedItem(); ok {
		return rev
	}
	return models.Revision{}
}

// GetCurrentContainer returns the current container
func (cm *CoreModel) GetCurrentContainer() models.Container {
	if container, ok := cm.stateManager.GetCurrentContainer(); ok {
		return container
	}
	// Fallback: try to get from containers page
	if container, ok := cm.pageManager.GetContainersPage().GetSelectedItem(); ok {
		return container
	}
	return models.Container{}
}

// SetStatusLine sets the global status line
func (cm *CoreModel) SetStatusLine(status string) {
	cm.stateManager.SetStatusLine(status)
}

// GetStatusLine returns the global status line
func (cm *CoreModel) GetStatusLine() string {
	return cm.stateManager.GetStatusLine()
}

// Context management

// SetShowContextList sets whether to show the context list
func (cm *CoreModel) SetShowContextList(show bool) {
	cm.stateManager.SetShowContextList(show)
}

// IsShowingContextList returns whether the context list is being shown
func (cm *CoreModel) IsShowingContextList() bool {
	return cm.stateManager.IsShowingContextList()
}

// SetContextList sets the context list model
func (cm *CoreModel) SetContextList(contextList list.Model) {
	cm.contextList = contextList
}

// GetContextList returns the context list model
func (cm *CoreModel) GetContextList() list.Model {
	return cm.contextList
}

// Page delegation methods

// HandleKeyMsg delegates key handling to the page manager
func (cm *CoreModel) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	return cm.pageManager.HandleKeyMsg(msg)
}

// UpdateTable delegates table updates to the page manager
func (cm *CoreModel) UpdateTable(msg tea.KeyMsg) tea.Cmd {
	return cm.pageManager.UpdateTable(msg)
}

// View delegates view rendering to the page manager
func (cm *CoreModel) View() string {
	return cm.pageManager.View()
}

// ViewWithHelpContext delegates view rendering to the page manager with help context
func (cm *CoreModel) ViewWithHelpContext(helpContext layouts.HelpContext) string {
	return cm.pageManager.ViewWithHelpContext(helpContext)
}

// GetCurrentPage returns the current page instance
func (cm *CoreModel) GetCurrentPage() interface{} {
	return cm.pageManager.GetCurrentPage()
}

// IsAnyFilterActive checks if any page has an active filter
func (cm *CoreModel) IsAnyFilterActive() bool {
	return cm.pageManager.IsAnyFilterActive()
}

// Layout management

// UpdateDimensions updates the terminal dimensions
func (cm *CoreModel) UpdateDimensions(width, height int) {
	cm.termW = width
	cm.termH = height
	cm.layoutSystem.SetDimensions(width, height)
}

// GetLayoutSystem returns the layout system
func (cm *CoreModel) GetLayoutSystem() *layouts.LayoutSystem {
	return cm.layoutSystem
}

// Helper methods

// formatAppID formats an app into an ID string
func (cm *CoreModel) formatAppID(app models.ContainerApp) string {
	return fmt.Sprintf("%s/%s", app.ResourceGroup, app.Name)
}

// Cache management

// SetContainersCache caches containers for a revision
func (cm *CoreModel) SetContainersCache(appID, revName string, containers []models.Container) {
	cm.stateManager.SetContainersForRevision(appID, revName, containers)
}

// GetContainersCache retrieves cached containers for a revision
func (cm *CoreModel) GetContainersCache(appID, revName string) ([]models.Container, bool) {
	return cm.stateManager.GetContainersForRevision(appID, revName)
}
