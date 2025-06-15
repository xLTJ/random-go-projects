package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"main/shodan"
	"os"
)

var (
	client     shodan.Client
	configPath = "/.config/shodan-client"
	rootCmd    = &cobra.Command{
		Use:   "shodan-client",
		Short: "its a client... for Shodan",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(viper.GetString("SHODAN_API_KEY")) == 0 {
				log.Fatalln("Run init before any other command")
			}

			client = shodan.NewClient(viper.GetString("SHODAN_API_KEY"))
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home, err := os.UserHomeDir()
	if _, err := os.Stat(home + configPath); err != nil {
		err = os.Mkdir(home+configPath, 0755)
	}
	cobra.CheckErr(err)

	viper.AddConfigPath(home + "/.config/shodan-client")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				log.Fatalf("Error creating config file: %v", err)
			}
		} else {
			log.Fatalln("Error reading config: ", err)
		}
	}
}
