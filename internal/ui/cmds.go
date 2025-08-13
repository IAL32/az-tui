package ui

import (
	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/providers"
	"github.com/IAL32/az-tui/internal/ui/core"

	tea "github.com/charmbracelet/bubbletea"
)

// -------------------------- Commands ----------------------------
// These functions now use the new core message system

func LoadAppsCmd(provider providers.DataProvider, rg string) tea.Cmd {
	return core.CreateLoadAppsCmd(provider, rg)
}

func LoadRevsCmd(provider providers.DataProvider, a models.ContainerApp) tea.Cmd {
	return core.CreateLoadRevisionsCmd(provider, a)
}

func LoadResourceGroupsCmd(provider providers.DataProvider) tea.Cmd {
	return core.CreateLoadResourceGroupsCmd(provider)
}

func LoadContainersCmd(provider providers.DataProvider, a models.ContainerApp, revName string) tea.Cmd {
	return core.CreateLoadContainersCmd(provider, a, revName)
}
