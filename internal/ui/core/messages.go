package core

import (
	"context"
	"time"

	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/providers"
	tea "github.com/charmbracelet/bubbletea"
)

// Message types for core coordination

// LoadResourceGroupsMsg represents a command to load resource groups
type LoadResourceGroupsMsg struct {
	Provider providers.DataProvider
}

// LoadAppsMsg represents a command to load apps
type LoadAppsMsg struct {
	Provider      providers.DataProvider
	ResourceGroup string
}

// LoadRevisionsMsg represents a command to load revisions
type LoadRevisionsMsg struct {
	Provider providers.DataProvider
	App      models.ContainerApp
}

// LoadContainersMsg represents a command to load containers
type LoadContainersMsg struct {
	Provider providers.DataProvider
	App      models.ContainerApp
	RevName  string
}

// Data loaded messages

// LoadedResourceGroupsMsg represents loaded resource groups data
type LoadedResourceGroupsMsg struct {
	ResourceGroups []models.ResourceGroup
	Error          error
}

// LoadedAppsMsg represents loaded apps data
type LoadedAppsMsg struct {
	Apps  []models.ContainerApp
	Error error
}

// LoadedRevisionsMsg represents loaded revisions data
type LoadedRevisionsMsg struct {
	Revisions []models.Revision
	Error     error
}

// LoadedContainersMsg represents loaded containers data
type LoadedContainersMsg struct {
	AppID      string
	RevName    string
	Containers []models.Container
	Error      error
}

// RevisionRestartedMsg represents a revision restart result
type RevisionRestartedMsg struct {
	AppID   string
	RevName string
	Error   error
}

// LeaveEnvVarsMsg represents leaving environment variables mode
type LeaveEnvVarsMsg struct{}

// Command creators that return tea.Cmd

// CreateLoadResourceGroupsCmd creates a command to load resource groups
func CreateLoadResourceGroupsCmd(provider providers.DataProvider) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		resourceGroups, err := provider.ListResourceGroups(ctx)
		return LoadedResourceGroupsMsg{ResourceGroups: resourceGroups, Error: err}
	}
}

// CreateLoadAppsCmd creates a command to load apps
func CreateLoadAppsCmd(provider providers.DataProvider, resourceGroup string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		apps, err := provider.ListContainerApps(ctx, resourceGroup)
		return LoadedAppsMsg{Apps: apps, Error: err}
	}
}

// CreateLoadRevisionsCmd creates a command to load revisions
func CreateLoadRevisionsCmd(provider providers.DataProvider, app models.ContainerApp) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		revisions, err := provider.ListRevisions(ctx, app.Name, app.ResourceGroup)
		return LoadedRevisionsMsg{Revisions: revisions, Error: err}
	}
}

// CreateLoadContainersCmd creates a command to load containers
func CreateLoadContainersCmd(provider providers.DataProvider, app models.ContainerApp, revName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		containers, err := provider.ListContainers(ctx, app, revName)
		appID := app.ResourceGroup + "/" + app.Name
		return LoadedContainersMsg{AppID: appID, RevName: revName, Containers: containers, Error: err}
	}
}
