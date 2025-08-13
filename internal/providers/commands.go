package providers

import (
	"github.com/IAL32/az-tui/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// CommandProvider defines the interface for executing Azure CLI operations
type CommandProvider interface {
	ExecIntoApp(app models.ContainerApp) tea.Cmd
	ExecIntoRevision(app models.ContainerApp, revision string) tea.Cmd
	ExecIntoContainer(app models.ContainerApp, revision, container string) tea.Cmd
	ShowAppLogs(app models.ContainerApp) tea.Cmd
	ShowRevisionLogs(app models.ContainerApp, revision string) tea.Cmd
	ShowContainerLogs(app models.ContainerApp, revision, container string) tea.Cmd
	RestartRevision(app models.ContainerApp, revision string) tea.Cmd
}
