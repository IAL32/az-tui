package ui

import (
	"context"
	"time"

	azure "az-tui/internal/azure"
	m "az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// -------------------------- Commands ----------------------------

func LoadAppsCmd(rg string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		apps, err := azure.ListContainerApps(ctx, rg)
		return loadedAppsMsg{apps, err}
	}
}

func LoadDetailsCmd(a m.ContainerApp) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		js, err := azure.GetAppDetails(ctx, a.Name, a.ResourceGroup)
		return loadedDetailsMsg{js, err}
	}
}

func LoadRevsCmd(a m.ContainerApp) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		revs, err := azure.ListRevisions(ctx, a.Name, a.ResourceGroup)
		return loadedRevsMsg{revs, err}
	}
}
