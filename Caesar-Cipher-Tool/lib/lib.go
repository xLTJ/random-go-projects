package lib

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"unicode"
)

func ShiftString(input string, shiftAmount int) (string, error) {
	if math.Abs(float64(shiftAmount)) > 26 {
		return "", fmt.Errorf("Shift must be equal or between -26 and 26\n")
	}

	var builder strings.Builder

	for _, character := range input {
		if !unicode.IsLetter(character) {
			builder.WriteRune(character)
			continue
		}

		isLowerCase := unicode.IsLower(character)
		character += int32(shiftAmount)

		if isLowerCase {
			if character > 'z' {
				character -= 26
			}
			if character < 'a' {
				character += 26
			}
		} else {
			if character > 'Z' {
				character -= 26
			}
			if character < 'A' {
				character += 26
			}
		}

		builder.WriteRune(character)
	}

	return builder.String(), nil
}

func ReadInputs(inputFile string, args []string) (string, error) {
	var builder strings.Builder

	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			return "", err
		}

		defer func() { _ = file.Close() }()

		_, err = io.Copy(&builder, file)
		if err != nil {
			return "", err
		}
	}

	if len(args) > 0 {
		for _, input := range args {
			builder.WriteString(input)
		}
	}

	input := builder.String()
	return input, nil
}

func HandleOutput(output string, outputFile string) error {
	if outputFile == "" {
		fmt.Print(output)
		return nil
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, strings.NewReader(output))
	if err != nil {
		return err
	}

	return nil
}
