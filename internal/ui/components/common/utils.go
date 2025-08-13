// Package common provides shared utilities and components for the UI.
package common

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// GetStatusColor returns the appropriate color for a status value
func GetStatusColor(status string) string {
	switch strings.ToLower(status) {
	case "running", "succeeded", "healthy", "ready", "active":
		return "#8c8" // Green
	case "failed", "error", "unhealthy", "critical":
		return "#c88" // Red
	case "pending", "provisioning", "starting", "updating":
		return "#cc8" // Yellow
	case "stopped", "inactive", "unknown", "-":
		return "#888" // Gray
	default:
		return "#888" // Default gray
	}
}

// ConfirmationDialog represents a confirmation dialog state
type ConfirmationDialog struct {
	Visible bool
	Text    string
}

// RenderConfirmationBox renders a confirmation dialog box
func RenderConfirmationBox(dialog ConfirmationDialog) string {
	if !dialog.Visible {
		return ""
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center).
		Render(dialog.Text + "\n\n[y] Yes  [n] No")

	return box
}
