package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/lib"
)

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().IntVarP(&shift, "shift", "s", 3, "Shift value for decryption (-26 <= [shift] <= 26)")
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt [flags] [arguments]",
	Short: "Decrypts the input using the Caesar Cipher. Requires either an input file or an argument",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := lib.ReadInputs(inputFile, args)
		if err != nil {
			return err
		}

		if input == "" {
			return fmt.Errorf("You must provide either an input file or an argument\n")
		}

		decryptedInput, err := lib.ShiftString(input, shift*-1)
		if err != nil {
			return err
		}

		err = lib.HandleOutput(decryptedInput, outputFile)
		if err != nil {
			return err
		}

		return nil
	},
}
