package cmd

import "github.com/spf13/cobra"

var (
	inputFile  string
	outputFile string
	shift      int

	rootCmd = &cobra.Command{
		Use:   "caesar-moment",
		Short: "It literally just encrypts and decrypts using the Caesar Cipher",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Input file (optional). If both input file and argument is supplied, the input file will be read first. Only encrypts letters.")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Output file (optional).")
}
