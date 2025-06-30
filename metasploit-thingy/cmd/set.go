package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// setCmd is used to set config values
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set variables used by the CLI",
	Long:  "Example: metasploit-thingy set --password Abscission616F33 --host 10.10.1.6 --port 55552",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		flagToConfigMap := map[string]string{
			"password": "msgrpc.password",
			"host":     "msgrpc.host",
			"port":     "msgrpc.port",
		}

		changedConfig := false

		for flagName, configKey := range flagToConfigMap {
			if cmd.Flags().Changed(flagName) {
				flag := cmd.Flags().Lookup(flagName)
				viper.Set(configKey, flag.Value.String())
				fmt.Printf("%s - updated\n", flagName)
				changedConfig = true
			}
		}

		if !changedConfig {
			fmt.Println("No config values provided")
			return
		}

		if err := viper.WriteConfig(); err != nil {
			log.Panic("\nFailed to write to config :c")
		}
		fmt.Println("\nConfig successfully updated")
	},
}

func init() {
	setCmd.Flags().String("password", "", "Set password")
	setCmd.Flags().String("host", "", "Set host")
	setCmd.Flags().String("port", "", "Set port")
	rootCmd.AddCommand(setCmd)
}
