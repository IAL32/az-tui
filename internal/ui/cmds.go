package ui

import (
	"context"
	"os/exec"
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

func LoadRevsCmd(a m.ContainerApp) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		revs, err := azure.ListRevisions(ctx, a.Name, a.ResourceGroup)
		return loadedRevsMsg{revs, err}
	}
}

func LoadResourceGroupsCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		resourceGroups, err := azure.ListResourceGroups(ctx)
		return loadedResourceGroupsMsg{resourceGroups, err}
	}
}

func LoadContainersCmd(a m.ContainerApp, revName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		containers, err := azure.ListContainersCmd(ctx, a, revName)
		return loadedContainersMsg{appID: a.Name, revName: revName, ctrs: containers, err: err}
	}
}

func RestartRevisionCmd(app m.ContainerApp, rev string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("az", "containerapp", "revision", "restart",
			"-n", app.Name, "-g", app.ResourceGroup, "--revision", rev)
		b, err := cmd.CombinedOutput()
		return revisionRestartedMsg{
			appID: appID(app), revName: rev, err: err, out: string(b),
		}
	}
}
