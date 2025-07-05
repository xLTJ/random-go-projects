package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd is used to set config values
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set variables used by the CLI",
	Long:  "Example: metasploit-thingy set --username msf --password Abscission616F33 --host 10.10.1.6 --port 55552",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		flagToConfigMap := map[string]string{
			"password": "msgrpc.password",
			"username": "msgrpc.username",
			"host":     "msgrpc.host",
			"port":     "msgrpc.port",
		}

		for flagName, configKey := range flagToConfigMap {
			if cmd.Flags().Changed(flagName) {
				flag := cmd.Flags().Lookup(flagName)
				viper.Set(configKey, flag.Value.String())
				fmt.Printf("%s - updated\n", flagName)
			}
		}

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to write to config: %v", err)
		}
		fmt.Println("\nConfig successfully updated")
		return nil
	},
}

func init() {
	setCmd.Flags().String("password", "", "Set password")
	setCmd.Flags().String("username", "", "Set username")
	setCmd.Flags().String("host", "", "Set host")
	setCmd.Flags().String("port", "", "Set port")
	setCmd.MarkFlagsOneRequired("password", "username", "host", "port")
}
