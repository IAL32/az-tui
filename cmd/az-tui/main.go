package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/IAL32/az-tui/internal/build"
	"github.com/IAL32/az-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	mockMode := flag.Bool("mock", false, "use mock data instead of Azure CLI")
	mockModeShort := flag.Bool("m", false, "use mock data instead of Azure CLI (shorthand)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("az-tui %s (%s, %s)\n", build.Version, build.Commit, build.Date)
		os.Exit(0)
	}

	// Use mock mode if either flag is set
	useMockMode := *mockMode || *mockModeShort

	model := ui.InitialModel(useMockMode)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
