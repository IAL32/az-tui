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
	flag.Parse()
	if *showVersion {
		fmt.Printf("az-tui %s (%s, %s)\n", build.Version, build.Commit, build.Date)
		os.Exit(0)
	}

	model := ui.InitialModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
