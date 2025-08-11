package main

import (
	"fmt"
	"os"

	"github.com/IAL32/az-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := ui.InitialModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
