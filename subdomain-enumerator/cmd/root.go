package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var domain string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "subdomain-enumerator",
	Short: "Enumerates the subdomains",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(searchCmd)
}

// initConfig reads in config file
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home + "/.config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("subdomain-enumerator")

	if err = viper.ReadInConfig(); err != nil {
		createConfig()
	}
}

// createConfig creates a new config file
func createConfig() {
	viper.Set("workers", 10)
	viper.Set("servers", []string{"8.8.8.8:53"})

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
