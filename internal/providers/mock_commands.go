package providers

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/IAL32/az-tui/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

// MockCommandProvider implements the CommandProvider interface using mock operations
type MockCommandProvider struct{}

// NewMockCommandProvider creates a new mock command provider
func NewMockCommandProvider() *MockCommandProvider {
	return &MockCommandProvider{}
}

func (m *MockCommandProvider) ExecIntoApp(app models.ContainerApp) tea.Cmd {
	return m.mockExecCommand(fmt.Sprintf("Executing shell into app '%s' (latest revision)", app.Name))
}

func (m *MockCommandProvider) ExecIntoRevision(app models.ContainerApp, revision string) tea.Cmd {
	return m.mockExecCommand(fmt.Sprintf("Executing shell into app '%s', revision '%s'", app.Name, revision))
}

func (m *MockCommandProvider) ExecIntoContainer(app models.ContainerApp, revision, container string) tea.Cmd {
	return m.mockExecCommand(fmt.Sprintf("Executing shell into container '%s' in app '%s', revision '%s'", container, app.Name, revision))
}

func (m *MockCommandProvider) ShowAppLogs(app models.ContainerApp) tea.Cmd {
	return m.mockLogsCommand(fmt.Sprintf("Showing logs for app '%s'", app.Name))
}

func (m *MockCommandProvider) ShowRevisionLogs(app models.ContainerApp, revision string) tea.Cmd {
	return m.mockLogsCommand(fmt.Sprintf("Showing logs for app '%s', revision '%s'", app.Name, revision))
}

func (m *MockCommandProvider) ShowContainerLogs(app models.ContainerApp, revision, container string) tea.Cmd {
	return m.mockLogsCommand(fmt.Sprintf("Showing logs for container '%s' in app '%s', revision '%s'", container, app.Name, revision))
}

func (m *MockCommandProvider) RestartRevision(app models.ContainerApp, revision string) tea.Cmd {
	return func() tea.Msg {
		// Simulate the restart operation
		time.Sleep(2 * time.Second)

		return revisionRestartedMsg{
			appID:   fmt.Sprintf("%s/%s", app.ResourceGroup, app.Name),
			revName: revision,
			err:     nil,
			out:     fmt.Sprintf("Mock: Successfully restarted revision '%s' for app '%s'", revision, app.Name),
		}
	}
}

// mockExecCommand creates a mock exec command that shows a message and simulates a shell
func (m *MockCommandProvider) mockExecCommand(message string) tea.Cmd {
	// Create a mock shell session
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`
echo "=== MOCK MODE ==="
echo "%s"
echo "This is a simulated shell session."
echo "In real mode, this would connect to the actual container."
echo "Press Ctrl+C to exit this mock session."
echo "=== MOCK MODE ==="
echo ""
echo "Mock shell session started..."
while true; do
    echo "Mock shell session running..."
    sleep 1
done
echo "Mock shell session ended."
`, message))

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} })
}

// mockLogsCommand creates a mock logs command that shows sample log output
func (m *MockCommandProvider) mockLogsCommand(message string) tea.Cmd {
	// Create a mock logs session with sample log output
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`
echo "=== MOCK MODE ==="
echo "%s"
echo "This is simulated log output."
echo "In real mode, this would show actual container logs."
echo "Press Ctrl+C to stop following logs."
echo "=== MOCK MODE ==="
echo ""

# Simulate streaming logs
i=0
while true; do
    timestamp=$(date '+%%Y-%%m-%%d %%H:%%M:%%S')
    echo "[$timestamp] INFO: Mock log entry $i - Application is running normally"
    sleep 1
    
    if [ $((i %% 3)) -eq 0 ]; then
        echo "[$timestamp] DEBUG: Mock debug message - Processing request batch $i"
    fi
    
    if [ $((i %% 5)) -eq 0 ]; then
        echo "[$timestamp] WARN: Mock warning - High memory usage detected: $((60 + i * 2))%%"
    fi
done
`, message))

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return tea.ExecProcess(cmd, func(error) tea.Msg { return noop{} })
}
