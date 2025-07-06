package session

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"metasploit-thingy/msfrpc"
	"os"
	"strconv"
)

var infoCmd = &cobra.Command{
	Use:   "info <session_id>",
	Short: "Show detailed session information",
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

		targetSessions, ok := sessions[sessionId]
		if !ok {
			return fmt.Errorf("session with id %d not found", sessionId)
		}

		table := tabby.New()
		table.AddHeader(fmt.Sprintf("Info for session: %d", sessionId))

		for _, line := range parseSessionData(targetSessions) {
			table.AddLine(line[0], line[1])
		}

		table.Print()
		return nil
	},
}

func parseSessionData(sessionData msfrpc.SessionListResp) [][]string {
	return [][]string{
		{"Id", strconv.Itoa(sessionData.ID)},
		{"Type", sessionData.Type},
		{"Information", sessionData.Info},
		{"Description", sessionData.Description},
		{"Connection", fmt.Sprintf("%s -> %s", sessionData.TunnelLocal, sessionData.TunnelPeer)},
		{"Exploit", sessionData.ViaExploit},
		{"Payload", sessionData.ViaPayload},
		{"Username", sessionData.Username},
	}
}
