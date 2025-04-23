package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:443")
	if err != nil {
		log.Fatalf("Failed to connect to shell: %v", err)
	}
	fmt.Printf("Connected to shell\n> ")

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		reader := bufio.NewReader(connection)
		for scanner.Scan() {
			// get input and write to bind shell
			input := scanner.Text()
			_, err = connection.Write([]byte(input + "\n"))
			if err != nil {
				log.Printf("Error sending to shell: %v\n", err)
				fmt.Print("> ")
				continue
			}

			// get header with content length (or "Error" if an error happened)
			header, err := reader.ReadString('\n')
			if strings.HasPrefix(header, "Error") {
				fmt.Printf("%s\n", header)
				fmt.Print("> ")
				continue
			}

			var contentLength = 0
			_, err = fmt.Sscanf(header, "Content-Length: %d", &contentLength)
			if err != nil {
				log.Printf("Error reading header: %v\n", err)
				continue
			}

			// read output
			buffer := make([]byte, contentLength)
			_, err = reader.Read(buffer)
			if err != nil {
				log.Printf("Error reading output: %v\n", err)
				continue
			}

			fmt.Println(string(buffer))
			fmt.Print("> ")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	_ = connection.Close()
}
