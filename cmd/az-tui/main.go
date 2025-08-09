// az-tui (Bubble Tea): K9s-like TUI prototype for Azure Container Apps
// MVP features:
// - Left: list of Container Apps (optional RG filter via ACA_RG)
// - Right (tabs): Details JSON, Revisions table
// - Keybindings: q quit | r refresh | tab switch pane | R reload revisions | l tail logs | s exec shell
// - Uses Azure CLI under the hood for fast iteration. Auth/ctx via your existing `az` session.
//
// Build & Run
//   go mod init az-tui
//   go get github.com/charmbracelet/bubbletea \
//          github.com/charmbracelet/bubbles \
//          github.com/charmbracelet/lipgloss
//   # requires: Azure CLI (`az`) + containerapp extension installed
//   go run .
//

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	azure "az-tui/internal/azure"
	ui "az-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// Guard: if Azure CLI extension is missing, surface a clean error
func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	_, err := azure.RunAz(ctx, "extension", "show", "-n", "containerapp", "-o", "none")
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.StyleError.Render("Azure CLI 'containerapp' extension not found. Install with: az extension add -n containerapp"))
		// not fatal, but commands will fail later; show friendly note.
	}
}
