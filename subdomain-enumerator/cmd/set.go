package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set default values to use",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// all flags except for help are used for setting a value in the config
		cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
			if flag.Name == "help" || !flag.Changed {
				return
			}
			fmt.Printf("%s - updated to %v\n", flag.Name, flag.Value)
		})

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to write to config: %v", err)
		}
		fmt.Println("\nConfig successfully updated")
		return nil
	},
}

func init() {
	setCmd.Flags().String("wordlist", "", "set wordlist")
	setCmd.Flags().StringSlice("servers", nil, "set dns server to use")
	setCmd.Flags().Int("workers", 1, "set how many workers to use")
	setCmd.MarkFlagsOneRequired("wordlist", "servers", "workers")

	cobra.CheckErr(viper.BindPFlags(setCmd.Flags()))
}
