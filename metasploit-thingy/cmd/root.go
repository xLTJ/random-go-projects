package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "metasploit-thingy",
	Short: "It does stuff",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in ~/.config directory with name "metasploit-thingy"
	viper.AddConfigPath(home + "/.config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("metasploit-thingy")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in. Otherwise create a new default config
	if err := viper.ReadInConfig(); err != nil {
		createConfig()
	}
}

// createConfig creates a new config file if none exists
func createConfig() {
	// viper.SetDefault() // in case any defaults are needed, use this

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsExist(err) {
			fmt.Println("config file already exists") // this really shouldnt happen but just in case lol
		} else {
			fmt.Printf("error creating config file: %v\n", err)
		}
		return
	}

	fmt.Println("config file created")
}
