package providers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/IAL32/az-tui/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// AzureCommandProvider implements the CommandProvider interface using Azure CLI
type AzureCommandProvider struct{}

// NewAzureCommandProvider creates a new Azure CLI command provider
func NewAzureCommandProvider() *AzureCommandProvider {
	return &AzureCommandProvider{}
}

// ExecIntoApp executes a shell into the app's latest revision
func (az *AzureCommandProvider) ExecIntoApp(app models.ContainerApp) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup, "--command", "/bin/sh")
}

// ExecIntoRevision executes a shell into a specific revision
func (az *AzureCommandProvider) ExecIntoRevision(app models.ContainerApp, revision string) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--command", "/bin/sh")
}

// ExecIntoContainer executes a shell into a specific container within a revision
func (az *AzureCommandProvider) ExecIntoContainer(app models.ContainerApp, revision, container string) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--container", container, "--command", "/bin/sh")
}

// ShowAppLogs shows logs for the entire app
func (az *AzureCommandProvider) ShowAppLogs(app models.ContainerApp) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup, "--follow")
}

// ShowRevisionLogs shows logs for a specific revision
func (az *AzureCommandProvider) ShowRevisionLogs(app models.ContainerApp, revision string) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--follow")
}

// ShowContainerLogs shows logs for a specific container within a revision
func (az *AzureCommandProvider) ShowContainerLogs(app models.ContainerApp, revision, container string) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--container", container, "--follow")
}

// RestartRevision restarts a specific revision
func (az *AzureCommandProvider) RestartRevision(app models.ContainerApp, revision string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("az", "containerapp", "revision", "restart",
			"-n", app.Name, "-g", app.ResourceGroup, "--revision", revision)
		b, err := cmd.CombinedOutput()
		return revisionRestartedMsg{
			appID:   fmt.Sprintf("%s/%s", app.ResourceGroup, app.Name),
			revName: revision,
			err:     err,
			out:     string(b),
		}
	}
}

// execCommand creates a tea.Cmd that executes the given command with proper I/O setup
func (az *AzureCommandProvider) execCommand(name string, args ...string) tea.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} })
}

// Message types
type revisionRestartedMsg struct {
	appID   string
	revName string
	err     error
	out     string
}

type noop struct{}
