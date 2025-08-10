package ui

import "github.com/charmbracelet/lipgloss"

var (
	styleTitle  = lipgloss.NewStyle().Bold(true)
	StyleError  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleAccent = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

func (m model) confirmBox() string {
	box := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(m.confirm.Text + "\n\n[y] Yes   [n] No")
	return box
}
