package main

import (
	"log"
	"main/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Error running tool: %v", err)
	}
}
