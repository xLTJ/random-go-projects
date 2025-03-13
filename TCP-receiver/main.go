package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// sets up a listener to listen on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:")

	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}
	fmt.Printf("TCP Server listening on: %s\n", listener.Addr().String())

	go func() {
		for {
			// wait for a connection to the listener
			connection, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				return
			}

			// handle the connection in a goroutine, so we can have multiple connections at the same time
			go handleConnection(connection)
		}
	}()

	// creates channel to hold one signal
	sigChan := make(chan os.Signal, 1)

	// when a SIGINT or a SIGTERM signal is sent, send it to the sigChan channel.
	// these signal are sent when the program terminates
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// wait for a signal to show up. when it does, it means the program should shut down, and so we close the listener.
	<-sigChan
	fmt.Println("Shutting down")
	_ = listener.Close()

	fmt.Println("Shutdown complete")
}

func handleConnection(connection net.Conn) {
	defer func() { _ = connection.Close() }()
	remoteAddress := connection.RemoteAddr().String()

	fmt.Printf("New connection from: %s\n", remoteAddress)

	// create buffer to store the message
	buffer := make([]byte, 1024)

	// keep listening for messages until the connection is stopped (usually due to an EOF when the client disconnects)
	for {
		// n is the length of the message
		n, err := connection.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			fmt.Printf("Connection from %s is closed", remoteAddress)
			return
		}

		// only grab the actual message
		message := string(buffer[:n])
		fmt.Printf("Message recieved from %s: %s\n", remoteAddress, message)
	}
}
