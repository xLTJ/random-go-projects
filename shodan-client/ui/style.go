package ui

import "github.com/charmbracelet/lipgloss"

var (
	Pink       = lipgloss.Color("212")
	Purple     = lipgloss.Color("99")
	Gray       = lipgloss.Color("248")
	lightGray  = lipgloss.Color("244")
	SubtleGray = lipgloss.Color("236")
	Red        = lipgloss.Color("196")
	Orange     = lipgloss.Color("208")
	Blue       = lipgloss.Color("26")

	TableHeaderStyle = lipgloss.NewStyle().
		Foreground(Pink).
		Bold(true)

	TableCellStyle = lipgloss.NewStyle().
		Padding(0, 2, 0, 0)

	TableOddRowStyle = TableCellStyle.
		Foreground(Gray)

	TableEvenRowStyle = TableCellStyle.
		Foreground(lightGray)

	TitleStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(SubtleGray).
		Padding(0, 4, 0, 1)

	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())

	Base = lipgloss.NewStyle()
)
