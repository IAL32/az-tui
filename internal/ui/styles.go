package ui

import "github.com/charmbracelet/lipgloss"

var (
	styleTitle  = lipgloss.NewStyle().Bold(true)
	StyleError  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleAccent = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)
