package app

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"main/internal/config"
	"main/internal/scanner"
	"main/internal/utils"
)

type scanResultMsg struct {
	scanResult scanner.ScanResult
	done       bool // if its done with all results
}

type Model struct {
	ports      []int
	openPorts  []int
	targetHost string
	resultChan chan scanner.ScanResult

	scannedCount int
	spinner      spinner.Model
	progress     progress.Model
	done         bool

	width  int
	height int
}

func NewModel(config config.ProgramConfig, targetHost string) (Model, error) {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	s := spinner.New()
	s.Style = SpinnerStyle

	ports, err := utils.ParsePorts(config.PortFlag)
	if err != nil {
		return Model{}, err
	}

	return Model{
		ports:      ports,
		targetHost: targetHost,
		resultChan: make(chan scanner.ScanResult, len(ports)),

		scannedCount: 0,
		spinner:      s,
		progress:     p,
		done:         false,
	}, nil
}

func (m Model) Init() tea.Cmd {
	go scanner.RunScan(m.ports, m.targetHost, m.resultChan)

	return tea.Batch(
		m.spinner.Tick,
		awaitScanResult(m.resultChan),
	)
}
