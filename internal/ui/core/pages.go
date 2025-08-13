package core

import (
	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/IAL32/az-tui/internal/ui/pages/apps"
	"github.com/IAL32/az-tui/internal/ui/pages/containers"
	"github.com/IAL32/az-tui/internal/ui/pages/envvars"
	"github.com/IAL32/az-tui/internal/ui/pages/resourcegroups"
	"github.com/IAL32/az-tui/internal/ui/pages/revisions"
	tea "github.com/charmbracelet/bubbletea"
)

// PageManager manages all page instances and their lifecycle
type PageManager struct {
	// Page instances
	resourceGroupsPage *resourcegroups.ResourceGroupsPage
	appsPage           *apps.AppsPage
	revisionsPage      *revisions.RevisionsPage
	containersPage     *containers.ContainersPage
	envVarsPage        *envvars.EnvVarsPage

	// Layout system
	layoutSystem *layouts.LayoutSystem

	// Navigation manager reference
	navigationManager *NavigationManager
}

// NewPageManager creates a new page manager
func NewPageManager(layoutSystem *layouts.LayoutSystem, navigationManager *NavigationManager) *PageManager {
	pm := &PageManager{
		layoutSystem:      layoutSystem,
		navigationManager: navigationManager,
	}

	// Initialize all pages
	pm.initializePages()

	return pm
}

// initializePages creates and configures all page instances
func (pm *PageManager) initializePages() {
	// Create page instances
	pm.resourceGroupsPage = resourcegroups.NewResourceGroupsPage(pm.layoutSystem)
	pm.appsPage = apps.NewAppsPage(pm.layoutSystem)
	pm.revisionsPage = revisions.NewRevisionsPage(pm.layoutSystem)
	pm.containersPage = containers.NewContainersPage(pm.layoutSystem)
	pm.envVarsPage = envvars.NewEnvVarsPage(pm.layoutSystem)
}

// SetupPageNavigation configures navigation functions between pages
func (pm *PageManager) SetupPageNavigation(coreModel *CoreModel) {
	// ResourceGroups -> Apps navigation
	pm.resourceGroupsPage.SetNavigateToAppsFunc(func(rg models.ResourceGroup) tea.Cmd {
		return coreModel.NavigateToApps(rg)
	})

	// Apps -> Revisions navigation
	pm.appsPage.SetNavigateToRevisionsFunc(func(app models.ContainerApp) tea.Cmd {
		return coreModel.NavigateToRevisions(app)
	})

	// Apps -> ResourceGroups back navigation
	pm.appsPage.SetBackToResourceGroupsFunc(func() tea.Cmd {
		return coreModel.NavigateToResourceGroups()
	})

	// Revisions -> Containers navigation
	pm.revisionsPage.SetNavigateToContainersFunc(func(rev models.Revision) tea.Cmd {
		return coreModel.NavigateToContainers(rev)
	})

	// Revisions -> Apps back navigation
	pm.revisionsPage.SetBackToAppsFunc(func() tea.Cmd {
		return coreModel.GoBack()
	})

	// Containers -> EnvVars navigation
	pm.containersPage.SetNavigateToEnvVarsFunc(func(container models.Container) tea.Cmd {
		return coreModel.NavigateToEnvVars(container)
	})

	// Containers -> Revisions back navigation
	pm.containersPage.SetBackToRevisionsFunc(func() tea.Cmd {
		return coreModel.GoBack()
	})

	// EnvVars -> Containers back navigation
	pm.envVarsPage.SetBackFunc(func() tea.Cmd {
		return coreModel.GoBack()
	})
}

// SetupPageActions configures action functions for pages
func (pm *PageManager) SetupPageActions(coreModel *CoreModel) {
	// Apps page actions
	pm.appsPage.SetShowLogsFunc(func(app models.ContainerApp) tea.Cmd {
		return coreModel.ShowAppLogs(app)
	})
	pm.appsPage.SetExecIntoAppFunc(func(app models.ContainerApp) tea.Cmd {
		return coreModel.ExecIntoApp(app)
	})

	// Revisions page actions
	pm.revisionsPage.SetRestartRevisionFunc(func(rev models.Revision) tea.Cmd {
		return coreModel.RestartRevision(rev)
	})
	pm.revisionsPage.SetShowLogsFunc(func(rev models.Revision) tea.Cmd {
		return coreModel.ShowRevisionLogs(rev)
	})
	pm.revisionsPage.SetExecIntoRevisionFunc(func(rev models.Revision) tea.Cmd {
		return coreModel.ExecIntoRevision(rev)
	})

	// Containers page actions
	pm.containersPage.SetShowLogsFunc(func(container models.Container) tea.Cmd {
		return coreModel.ShowContainerLogs(container)
	})
	pm.containersPage.SetExecIntoContainerFunc(func(container models.Container) tea.Cmd {
		return coreModel.ExecIntoContainer(container)
	})
}

// GetCurrentPage returns the page instance for the current mode
func (pm *PageManager) GetCurrentPage() interface{} {
	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		return pm.resourceGroupsPage
	case ModeApps:
		return pm.appsPage
	case ModeRevisions:
		return pm.revisionsPage
	case ModeContainers:
		return pm.containersPage
	case ModeEnvVars:
		return pm.envVarsPage
	default:
		return pm.resourceGroupsPage
	}
}

// GetResourceGroupsPage returns the resource groups page
func (pm *PageManager) GetResourceGroupsPage() *resourcegroups.ResourceGroupsPage {
	return pm.resourceGroupsPage
}

// GetAppsPage returns the apps page
func (pm *PageManager) GetAppsPage() *apps.AppsPage {
	return pm.appsPage
}

// GetRevisionsPage returns the revisions page
func (pm *PageManager) GetRevisionsPage() *revisions.RevisionsPage {
	return pm.revisionsPage
}

// GetContainersPage returns the containers page
func (pm *PageManager) GetContainersPage() *containers.ContainersPage {
	return pm.containersPage
}

// GetEnvVarsPage returns the environment variables page
func (pm *PageManager) GetEnvVarsPage() *envvars.EnvVarsPage {
	return pm.envVarsPage
}

// HandleKeyMsg delegates key handling to the current page
func (pm *PageManager) HandleKeyMsg(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		return pm.resourceGroupsPage.HandleKeyMsg(msg)
	case ModeApps:
		return pm.appsPage.HandleKeyMsg(msg)
	case ModeRevisions:
		return pm.revisionsPage.HandleKeyMsg(msg)
	case ModeContainers:
		return pm.containersPage.HandleKeyMsg(msg)
	case ModeEnvVars:
		return pm.envVarsPage.HandleKeyMsg(msg)
	default:
		return nil, false
	}
}

// UpdateTable updates the table for the current page
func (pm *PageManager) UpdateTable(msg tea.KeyMsg) tea.Cmd {
	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		table := pm.resourceGroupsPage.GetTable()
		table, cmd := table.Update(msg)
		pm.resourceGroupsPage.SetTable(table)
		return cmd
	case ModeApps:
		table := pm.appsPage.GetTable()
		table, cmd := table.Update(msg)
		pm.appsPage.SetTable(table)
		return cmd
	case ModeRevisions:
		table := pm.revisionsPage.GetTable()
		table, cmd := table.Update(msg)
		pm.revisionsPage.SetTable(table)
		return cmd
	case ModeContainers:
		table := pm.containersPage.GetTable()
		table, cmd := table.Update(msg)
		pm.containersPage.SetTable(table)
		return cmd
	case ModeEnvVars:
		// Handle table updates for envvars page like other pages
		table := pm.envVarsPage.GetTable()
		table, cmd := table.Update(msg)
		pm.envVarsPage.SetTable(table)
		return cmd
	default:
		return nil
	}
}

// View renders the current page
func (pm *PageManager) View() string {
	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		return pm.resourceGroupsPage.View()
	case ModeApps:
		return pm.appsPage.View()
	case ModeRevisions:
		return pm.revisionsPage.View()
	case ModeContainers:
		return pm.containersPage.View()
	case ModeEnvVars:
		return pm.envVarsPage.View()
	default:
		return pm.resourceGroupsPage.View()
	}
}

// ViewWithHelpContext renders the current page with help context
func (pm *PageManager) ViewWithHelpContext(helpContext layouts.HelpContext) string {
	// Set the mode in the help context based on current mode
	helpContext.Mode = pm.navigationManager.GetCurrentMode().String()

	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		return pm.resourceGroupsPage.ViewWithHelpContext(helpContext)
	case ModeApps:
		return pm.appsPage.ViewWithHelpContext(helpContext)
	case ModeRevisions:
		return pm.revisionsPage.ViewWithHelpContext(helpContext)
	case ModeContainers:
		return pm.containersPage.ViewWithHelpContext(helpContext)
	case ModeEnvVars:
		return pm.envVarsPage.ViewWithHelpContext(helpContext)
	default:
		return pm.resourceGroupsPage.ViewWithHelpContext(helpContext)
	}
}

// SetLoading sets loading state for the current page
func (pm *PageManager) SetLoading(loading bool) {
	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		pm.resourceGroupsPage.SetLoading(loading)
	case ModeApps:
		pm.appsPage.SetLoading(loading)
	case ModeRevisions:
		pm.revisionsPage.SetLoading(loading)
	case ModeContainers:
		pm.containersPage.SetLoading(loading)
	case ModeEnvVars:
		pm.envVarsPage.SetLoading(loading)
	}
}

// SetError sets error state for the current page
func (pm *PageManager) SetError(err error) {
	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		pm.resourceGroupsPage.SetError(err)
	case ModeApps:
		pm.appsPage.SetError(err)
	case ModeRevisions:
		pm.revisionsPage.SetError(err)
	case ModeContainers:
		pm.containersPage.SetError(err)
	case ModeEnvVars:
		pm.envVarsPage.SetError(err)
	}
}

// ClearData clears data for the current page
func (pm *PageManager) ClearData() {
	switch pm.navigationManager.GetCurrentMode() {
	case ModeResourceGroups:
		pm.resourceGroupsPage.ClearData()
	case ModeApps:
		pm.appsPage.ClearData()
	case ModeRevisions:
		pm.revisionsPage.ClearData()
	case ModeContainers:
		pm.containersPage.ClearData()
	case ModeEnvVars:
		pm.envVarsPage.ClearData()
	}
}

// IsAnyFilterActive checks if any page has an active filter
func (pm *PageManager) IsAnyFilterActive() bool {
	return pm.resourceGroupsPage.GetFilterInput().Focused() ||
		pm.appsPage.GetFilterInput().Focused() ||
		pm.revisionsPage.GetFilterInput().Focused() ||
		pm.containersPage.GetFilterInput().Focused() ||
		pm.envVarsPage.GetFilterInput().Focused()
}

// UpdateLayoutSystem updates the layout system for all pages
func (pm *PageManager) UpdateLayoutSystem(layoutSystem *layouts.LayoutSystem) {
	pm.layoutSystem = layoutSystem
	// Note: Pages hold references to the layout system, so they'll automatically use the updated one
}
