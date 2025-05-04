package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"main/internal/app"
	"main/internal/config"
	"os"
)

var (
	programConfig = config.ProgramConfig{}
	rootCmd       = &cobra.Command{
		Use:   "TNNmap",
		Short: "Totally Not Nmap: One of the port scanners of all time",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			targetHost := args[0]

			model, err := app.NewModel(programConfig, targetHost)
			if err != nil {
				return err
			}

			if _, err := tea.NewProgram(model).Run(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}

			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&programConfig.PortFlag, "ports", "p", "", "ports to scan. Can be a single port (eg 143), multiple ports (eg 143,443,8000) or a range (eg 1-100)")
	rootCmd.Flags().BoolVarP(&programConfig.ListAll, "listall", "l", false, "List all ports (also closed ones)")
}
