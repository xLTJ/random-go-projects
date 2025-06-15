package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for something",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostSearch, err := client.HostSearch(args[0])
		cobra.CheckErr(err)

		for _, host := range hostSearch.Matches {
			fmt.Printf("%18s%8d\n", host.IPString, host.Port)
		}
	},
}
