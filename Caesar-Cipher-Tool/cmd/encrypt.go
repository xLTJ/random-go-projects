package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/lib"
)

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().IntVarP(&shift, "shift", "s", 3, "Shift value for encryption (-26 <= [shift] <= 26)")
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt [flags] [arguments]",
	Short: "Encrypt the input using the Caesar Cipher. Requires either an input file or an argument",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := lib.ReadInputs(inputFile, args)
		if err != nil {
			return err
		}

		if input == "" {
			return fmt.Errorf("You must provide either an input file or an argument\n")
		}

		encryptedInput, err := lib.ShiftString(input, shift)
		if err != nil {
			return err
		}

		err = lib.HandleOutput(encryptedInput, outputFile)
		if err != nil {
			return err
		}

		return nil
	},
}
