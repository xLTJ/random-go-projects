package cmd

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"metasploit-thingy/msfrpc"
	"os"
)

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Get a list of current sessions",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		msfrpcClient, err := msfrpc.NewClient()
		if err != nil {
			return err
		}

		defer func() {
			err := msfrpcClient.Logout()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to logout: %v", err)
			}
		}()

		sessions, err := msfrpcClient.SessionList()
		if err != nil {
			return err
		}

		table := tabby.New()
		table.AddHeader("Id", "Type", "Information", "Description", "Connection")
		for _, session := range sessions {
			table.AddLine(
				session.ID,
				session.Type,
				session.Info,
				session.Description,
				fmt.Sprintf("%s -> %s", session.TunnelLocal, session.TunnelPeer),
			)
		}

		table.Print()
		return nil
	},
}
