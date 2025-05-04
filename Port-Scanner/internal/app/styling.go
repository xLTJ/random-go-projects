package app

import "github.com/charmbracelet/lipgloss"

var (
	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a684ff"))

	DoneStyle = lipgloss.NewStyle().Margin(1, 2)

	CurrentPortStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed6bff"))

	PortOpen = lipgloss.NewStyle().Foreground(lipgloss.Color("#7bf1a8")).SetString("✓")

	PortClosed = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6467")).SetString("x")
)
