package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func startChatClient() {
	fmt.Println("Enter your username: ")
	var username string
	_, err := fmt.Scanf("%s", &username)
	if err != nil {
		log.Fatalf("What the fuck did you do, literally just put in a normal string its not that hard")
		return
	}

	serverAddr := net.UDPAddr{
		IP:   net.ParseIP("localhost"),
		Port: 1234,
	}

	var connection *net.UDPConn
	connection, err = net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to server")
	}
	fmt.Println("Connected to server")

	_, err = connection.Write([]byte("!SetUsername " + username))
	if err != nil {
		log.Fatalf("Error writing to server: %v", err)
	}

	if err != nil {
		log.Fatalf("Error writing to server: %v", err)
	}

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, _ := connection.ReadFromUDP(buffer)
			fmt.Println(string(buffer[:n]))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err = connection.Write([]byte(scanner.Text()))
			if err != nil {
				log.Fatalf("Error writing to server: %v", err)
			}
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
	_, err = connection.Write([]byte("!Disconnect"))
	if err != nil {
		log.Fatalf("Error writing to server: %v", err)
	}
	_ = connection.Close()

	fmt.Println("Shutdown complete")
}
