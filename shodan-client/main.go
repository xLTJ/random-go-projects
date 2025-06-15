package main

import (
	"log"
	"main/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalln("Error running tool: ", err)
	}
}
