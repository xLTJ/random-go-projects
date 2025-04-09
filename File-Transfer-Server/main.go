package main

import "fmt"

func main() {
	input := 0

	for input != 1 && input != 2 {
		fmt.Println("Chose option:\n1: Start file server\n2: Start client")
		_, _ = fmt.Scanf("%d", &input)
	}

	if input == 1 {
		StartFileServer()
	} else {
		StartClient()
	}
}
