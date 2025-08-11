package ui

import (
	"context"
	"time"

	m "github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/providers"

	tea "github.com/charmbracelet/bubbletea"
)

// -------------------------- Commands ----------------------------

func LoadAppsCmd(provider providers.DataProvider, rg string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		apps, err := provider.ListContainerApps(ctx, rg)
		return loadedAppsMsg{apps, err}
	}
}

func LoadRevsCmd(provider providers.DataProvider, a m.ContainerApp) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		revs, err := provider.ListRevisions(ctx, a.Name, a.ResourceGroup)
		return loadedRevsMsg{revs, err}
	}
}

func LoadResourceGroupsCmd(provider providers.DataProvider) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		resourceGroups, err := provider.ListResourceGroups(ctx)
		return loadedResourceGroupsMsg{resourceGroups, err}
	}
}

func LoadContainersCmd(provider providers.DataProvider, a m.ContainerApp, revName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		containers, err := provider.ListContainers(ctx, a, revName)
		return loadedContainersMsg{appID: a.Name, revName: revName, ctrs: containers, err: err}
	}
}
