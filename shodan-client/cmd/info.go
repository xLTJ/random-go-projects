package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(apiInfoCmd)
}

var apiInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Shows info about the API plan",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := client.ApiInfo()
		cobra.CheckErr(err)

		fmt.Printf(
			"Plan: %s\nQuery credits: %d\nScan credits: %d\n",
			info.Plan,
			info.QueryCredits,
			info.ScanCredits,
		)
	},
}
