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
	listener, err := net.Listen("tcp", ":7263")
	if err != nil {
		log.Fatalln("Error setting up server: ", err)
	}
	fmt.Printf("Listening on port: %s\n", listener.Addr().String())

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Unable to accept connection: ", err)
		}

		go handleClient(connection)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down")
	_ = listener.Close()
	log.Println("Shutdown complete")
}

func handleClient(connection net.Conn) {
	_, err := io.Copy(os.Stdout, connection)
	if err != nil {
		log.Println("Error writing to Stdout: ", err)
		return
	}
}
