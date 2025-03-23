package main

import "fmt"

// still got more stuff to do, but it works at the moment (provided u dont do any weird stuff)
func main() {
	input := 0
	for input != 1 && input != 2 {
		fmt.Println("Chose option:\n1: Start chat server\n2: Start chat client")
		_, _ = fmt.Scanf("%d", &input)
	}

	if input == 1 {
		startServer()
	} else {
		startChatClient()
	}
}
