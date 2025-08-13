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

func (az *AzureCommandProvider) ExecIntoApp(app models.ContainerApp) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup, "--command", "/bin/sh")
}

func (az *AzureCommandProvider) ExecIntoRevision(app models.ContainerApp, revision string) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--command", "/bin/sh")
}

func (az *AzureCommandProvider) ExecIntoContainer(app models.ContainerApp, revision, container string) tea.Cmd {
	return az.execCommand("az", "containerapp", "exec",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--container", container, "--command", "/bin/sh")
}

func (az *AzureCommandProvider) ShowAppLogs(app models.ContainerApp) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup, "--follow")
}

func (az *AzureCommandProvider) ShowRevisionLogs(app models.ContainerApp, revision string) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--follow")
}

func (az *AzureCommandProvider) ShowContainerLogs(app models.ContainerApp, revision, container string) tea.Cmd {
	fmt.Println("--- Ctrl+C to stop logs ---")
	return az.execCommand("az", "containerapp", "logs", "show",
		"-n", app.Name, "-g", app.ResourceGroup,
		"--revision", revision, "--container", container, "--follow")
}

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
