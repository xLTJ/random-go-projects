package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// ip and port to listen on
	addr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 1234,
	}

	// sets up listener
	listener, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}
	fmt.Printf("UDP server listening on %s\n", listener.LocalAddr().String())

	go func() {
		// buffer to store messages
		buffer := make([]byte, 1024)

		for {
			// wait for a packet, copies the content to the buffer. n is the length os the message
			n, remoteAddr, err := listener.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error reading from UDP: %v\n", err)
				return
			}

			// grabs the actual message and prints it.
			message := string(buffer[:n])
			fmt.Printf("Message from %s: %s\n", remoteAddr.String(), message)
		}
	}()

	// this part is just the same as with the TCP receiver lol
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
