package main

import (
	"fmt"
	"strings"
)

type Command struct {
	command   string
	arguments []string
}

const (
	CmdUpload   = "UPLOAD"
	CmdDownload = "GET"
	CmdHelp     = "HELP"
)

const fileDirectory = "files"

// parseCommand takes an input string and splits it into command and arguments, then returns a Command with those values
func parseCommand(input string) (Command, error) {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) == 0 {
		return Command{}, fmt.Errorf("empty command")
	}

	command := parts[0]
	arguments := parts[1:]

	return Command{command: command, arguments: arguments}, nil
}
