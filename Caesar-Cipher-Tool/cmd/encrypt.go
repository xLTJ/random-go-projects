package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"unicode"
)

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().IntVar(&shift, "shift", 3, "Shift value for encryption")
}

var (
	shift int

	encryptCmd = &cobra.Command{
		Use:   "encrypt [OPTIONS] [ARGUMENTS]",
		Short: "Encrypt the input using the Caesar Cipher. Requires either an input file or an argument",
		RunE: func(cmd *cobra.Command, args []string) error {
			var builder strings.Builder

			if inputFile != "" {
				file, err := os.Open(inputFile)
				if err != nil {
					return err
				}

				defer func() { _ = file.Close() }()

				_, err = io.Copy(&builder, file)
				if err != nil {
					return err
				}
			}

			if len(args) > 0 {
				for _, input := range args {
					builder.WriteString(input)
				}
			}

			input := builder.String()
			if input == "" {
				return fmt.Errorf("You must provide either an input file or an argument\n")
			}

			encryptedInput := encrypt(input)

			if outputFile == "" {
				fmt.Print(encryptedInput)
				return nil
			}

			file, err := os.Create(outputFile)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, strings.NewReader(encryptedInput))
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func encrypt(input string) string {
	var builder strings.Builder

	for _, character := range input {
		if !unicode.IsLetter(character) {
			builder.WriteRune(character)
			continue
		}

		isLowerCase := unicode.IsLower(character)
		character += int32(shift)

		if isLowerCase {
			for character > 'z' {
				character -= 26
			}
		} else {
			for character > 'Z' {
				character -= 26
			}
		}

		builder.WriteRune(character)
	}

	return builder.String()
}
