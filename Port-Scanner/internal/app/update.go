package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"main/internal/scanner"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}

	case scanResultMsg:
		m.scannedCount++

		progressCmd := m.progress.SetPercent(float64(m.scannedCount) / float64(len(m.ports)))
		commands := []tea.Cmd{progressCmd}

		// if the port is open, print that
		if msg.scanResult.IsOpen {
			m.openPorts = append(m.openPorts, msg.scanResult.Port)
			status := fmt.Sprintf("%s Port %d: Open", PortOpen, msg.scanResult.Port)
			commands = append(commands, tea.Printf(status))
		}

		if m.scannedCount >= len(m.ports) {
			m.done = true
			commands = append(commands, tea.Quit)
			return m, tea.Batch(commands...)
		}
		
		commands = append(commands, awaitScanResult(m.resultChan))

		return m, tea.Batch(commands...)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}

	return m, nil
}

func awaitScanResult(resultChan chan scanner.ScanResult) tea.Cmd {
	return func() tea.Msg {
		result, ok := <-resultChan
		// if chan is closed, all ports has been scanned
		if !ok {
			return scanResultMsg{done: true}
		}

		return scanResultMsg{scanResult: result, done: false}
	}
}
