package session

import (
	"fmt"
	"github.com/spf13/cobra"
)

var SessionCmd = &cobra.Command{
	Use:   "session <command> [arguments]",
	Short: "Manage individual Metasploit sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("lol")
		return nil
	},
}

func init() {
	SessionCmd.AddCommand(infoCmd)
	SessionCmd.AddCommand(killCmd)
}
