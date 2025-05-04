package app

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strings"
)

func (m Model) View() string {
	n := len(m.ports)
	width := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		sort.Ints(m.openPorts)

		result := fmt.Sprintf("Finished scanning, found %d open ports on %s\n\n", len(m.openPorts), m.targetHost)

		if len(m.openPorts) > 0 {
			result += "Open ports:\n"
			for _, port := range m.openPorts {
				result += fmt.Sprintf(" - %d\n", port)
			}
		}

		return DoneStyle.Render(result)
	}

	portCount := fmt.Sprintf(" %*d/%*d (%d/%d open)", width, m.scannedCount, width, n, len(m.openPorts), m.scannedCount)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvailable := max(0, m.width-lipgloss.Width(spin+prog+portCount))

	info := fmt.Sprintf("Scanning %s", m.targetHost)
	info = lipgloss.NewStyle().MaxWidth(cellsAvailable).Render(info)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+portCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + portCount
}
