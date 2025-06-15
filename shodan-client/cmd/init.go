package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [api-key]",
	Short: "Initialize the client with your Shodan api key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("SHODAN_API_KEY", args[0])
		err := viper.WriteConfig()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Succesfully initialized with api key: %s\n", args[0])
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// empty to prevent config checking
	},
}
