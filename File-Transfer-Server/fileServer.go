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

type Client struct {
	connection net.Conn
}

func StartFileServer() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}
	fmt.Printf("File server listening on: %s\n", listener.Addr().String())

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v\n", err)
				continue
			}

			client := Client{connection: connection}

			go client.handleClient()
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("Shutting down")
	_ = listener.Close()

	fmt.Println("Shutdown complete")
}

func (c Client) handleClient() {
	defer func() { _ = c.connection.Close() }()

	fmt.Printf("New connection from: %s\n", c.connection.RemoteAddr().String())

	reader := bufio.NewReader(c.connection)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		command, err := parseCommand(input)
		if err != nil {
			break
		}

		err = c.handleCommand(command)
		if err != nil {
			log.Printf("Error executing '%s' command: %v\n", command.command, err)
		}
	}
}

func (c Client) handleCommand(command Command) error {
	switch command.command {
	case CmdUpload:
		if len(command.arguments) != 2 {
			return fmt.Errorf("%s command can have exactly 2 arguments", CmdUpload)
		}

		fmt.Printf("Executing %s with arguments: %s %s\n", CmdUpload, command.arguments[0], command.arguments[1])

		fileSize, err := strconv.Atoi(command.arguments[1])
		if err != nil {
			return err
		}

		err = c.handleFileUpload(command.arguments[0], fileSize)
		if err != nil {
			return err
		}

	case CmdDownload:
		if len(command.arguments) != 1 {
			return fmt.Errorf("%s command can have exactly 1 arguments", CmdDownload)
		}

		fmt.Printf("Executing %s with argument: %s\n", CmdDownload, command.arguments[0])
		err := c.handleFileDownload(command.arguments[0])
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown command: %s\n", command.command)
	}

	return nil
}

func (c Client) handleFileUpload(filename string, fileSize int) error {
	file, err := os.Create(filepath.Join(fileDirectory, filename))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	buffer := make([]byte, 1024)
	bytesTotal := 0

	fmt.Printf("Started writing to: %s...\n", filename)
	for bytesTotal < fileSize {
		bytesRead, err := c.connection.Read(buffer)
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
		fmt.Printf("Writing to file... %d/%d bytes (%.2f%%)\n", bytesTotal, fileSize, float64(bytesTotal)/float64(fileSize)*100)
	}

	fmt.Printf("Finished writing to: %s", filename)
	return nil
}

func (c Client) handleFileDownload(filename string) error {
	if !isPathSafe(filename) {
		pathError := fmt.Errorf("requested file outside of allowed path")
		c.sendError(pathError)
		return pathError
	}

	file, err := os.Open(filepath.Join(fileDirectory, filename))
	if err != nil {
		c.sendError(err)
		return err
	}
	defer func() { _ = file.Close() }()

	fileInfo, err := file.Stat()
	if err != nil {
		c.sendError(err)
		return err
	}

	_, err = fmt.Fprintf(c.connection, "%d\n", fileInfo.Size())
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
		fmt.Printf("Sending file... %d/%d bytes (%.2f%%)\n", bytesTotal, fileInfo.Size(), float64(bytesTotal)/float64(fileInfo.Size())*100)
	}

	fmt.Println("Finished sending file")
	return nil
}

// sendError sends an error to the client
func (c Client) sendError(err error) {
	_, _ = fmt.Fprintf(c.connection, "Error: %v", err)
}

// isPathSafe checks if a path is inside the allowed path. This prevents the client to access files outside the intended
// path
func isPathSafe(path string) bool {
	fullPath := filepath.Join(fileDirectory, path)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return false
	}

	absBasePath, err := filepath.Abs(fileDirectory)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, absBasePath)
}
