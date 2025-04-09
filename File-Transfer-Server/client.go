package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type ServerConnection struct {
	connection net.Conn
}

func StartClient() {
	connection, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer func() { _ = connection.Close() }()

	fmt.Println("Connected to server")

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		serverConnection := ServerConnection{connection: connection}
		for scanner.Scan() {
			command, err := parseCommand(scanner.Text())
			if err != nil {
				log.Printf("Error parsing command: %v\n", err)
			}

			err = serverConnection.handleCommand(command)
			if err != nil {
				log.Printf("Error executing '%s' command: %v\n", command.command, err)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("Shutting down")
	_ = connection.Close()

	fmt.Println("Shutdown complete")
}

func (c ServerConnection) handleCommand(command Command) error {
	switch command.command {
	case CmdUpload:
		if len(command.arguments) != 1 {
			return fmt.Errorf("%s command can have exactly 1 argument", CmdUpload)
		}

		err := c.uploadFile(command.arguments[0])
		if err != nil {
			return err
		}

	case CmdDownload:
		if len(command.arguments) != 1 {
			return fmt.Errorf("%s command can have exactly 1 argument", CmdDownload)
		}

		err := c.downloadFile(command.arguments[0])
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown command: %s\n", command.command)
	}

	return nil
}

func (c ServerConnection) uploadFile(filename string) error {
	file, err := os.Open(filepath.Join(fileDirectory, filename))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(c.connection, "%s %s %d\n", CmdUpload, filepath.Base(filename), fileInfo.Size())
	if err != nil {
		return err
	}

	buffer := make([]byte, 4096)
	bytesTotal := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		_, err = c.connection.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}

		bytesTotal += bytesRead
		fmt.Printf("Uploading file... %d/%d bytes (%.2f%%)\n", bytesTotal, fileInfo.Size(), float64(bytesTotal)/float64(fileInfo.Size())*100)
	}

	fmt.Println("Upload completed !!")
	return nil
}

func (c ServerConnection) downloadFile(filename string) error {
	fmt.Printf("Starting download of: %s...\n", filename)
	_, err := fmt.Fprintf(c.connection, "%s %s\n", CmdDownload, filename)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	bytesRead, err := c.connection.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	statusMessage := string(buffer[:bytesRead])
	if strings.HasPrefix(statusMessage, "Error") {
		return fmt.Errorf(statusMessage)
	}

	fileSize, err := strconv.Atoi(string(buffer[:bytesRead-len("\n")]))
	if err != nil {
		return err
	}
	fmt.Printf("File size: %d\n", fileSize)
	bytesTotal := 0

	file, err := os.Create(filepath.Join(fileDirectory, filename))
	defer func() { _ = file.Close() }()
	for bytesTotal < fileSize {
		bytesRead, err = c.connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		_, err = file.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}

		bytesTotal += bytesRead
		fmt.Printf("Downloading %s... %d/%d bytes (%.2f%%)\n", filename, bytesTotal, fileSize, float64(bytesTotal)/float64(fileSize)*100)
	}

	fmt.Printf("Finished downloading %s !!", filename)
	return nil
}
