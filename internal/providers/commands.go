package providers

import (
	"github.com/IAL32/az-tui/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// CommandProvider defines the interface for executing Azure CLI operations
type CommandProvider interface {
	// ExecIntoApp executes a shell into the app's latest revision
	ExecIntoApp(app models.ContainerApp) tea.Cmd

	// ExecIntoRevision executes a shell into a specific revision
	ExecIntoRevision(app models.ContainerApp, revision string) tea.Cmd

	// ExecIntoContainer executes a shell into a specific container within a revision
	ExecIntoContainer(app models.ContainerApp, revision, container string) tea.Cmd

	// ShowAppLogs shows logs for the entire app
	ShowAppLogs(app models.ContainerApp) tea.Cmd

	// ShowRevisionLogs shows logs for a specific revision
	ShowRevisionLogs(app models.ContainerApp, revision string) tea.Cmd

	// ShowContainerLogs shows logs for a specific container within a revision
	ShowContainerLogs(app models.ContainerApp, revision, container string) tea.Cmd

	// RestartRevision restarts a specific revision
	RestartRevision(app models.ContainerApp, revision string) tea.Cmd
}
