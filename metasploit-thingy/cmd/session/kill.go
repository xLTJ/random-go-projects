package session

import (
	"fmt"
	"github.com/spf13/cobra"
	"metasploit-thingy/msfrpc"
	"strconv"
)

var killCmd = &cobra.Command{
	Use:   "kill <session_id>",
	Short: "Stop the specified session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionId, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid input: %v", err)
		}

		msfrpcClient, err := msfrpc.NewClient()
		if err != nil {
			return err
		}

		err = msfrpcClient.KillSession(sessionId)
		if err != nil {
			return err
		}

		fmt.Printf("Killed session: %d\n", sessionId)
		return nil
	},
}
