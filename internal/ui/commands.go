package ui

import (
	"fmt"
	"os"
	"os/exec"

	models "github.com/IAL32/az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// AzureCommands centralizes Azure CLI command execution
type AzureCommands struct{}

// NewAzureCommands creates a new Azure commands handler
func NewAzureCommands() *AzureCommands {
	return &AzureCommands{}
}

// ExecIntoApp executes a shell into the app's latest revision
func (az *AzureCommands) ExecIntoApp(app models.ContainerApp) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup, "--command", "/bin/sh")
}

// ExecIntoRevision executes a shell into a specific revision
func (az *AzureCommands) ExecIntoRevision(app models.ContainerApp, revision string) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--command", "/bin/sh")
}

// ExecIntoContainer executes a shell into a specific container within a revision
func (az *AzureCommands) ExecIntoContainer(app models.ContainerApp, revision, container string) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--container", container, "--command", "/bin/sh")
}

// ShowAppLogs shows logs for the entire app
func (az *AzureCommands) ShowAppLogs(app models.ContainerApp) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup, "--follow")
}

// ShowRevisionLogs shows logs for a specific revision
func (az *AzureCommands) ShowRevisionLogs(app models.ContainerApp, revision string) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--follow")
}

// ShowContainerLogs shows logs for a specific container within a revision
func (az *AzureCommands) ShowContainerLogs(app models.ContainerApp, revision, container string) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--container", container, "--follow")
}

// execCommand creates a tea.Cmd that executes the given command with proper I/O setup
func (az *AzureCommands) execCommand(name string, args ...string) tea.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} })
}
