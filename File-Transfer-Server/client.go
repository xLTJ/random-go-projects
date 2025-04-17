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
	connection    net.Conn
	currentFolder string
}

func StartClient() {
	connection, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer func() { _ = connection.Close() }()

	fmt.Println("Connected to server")
	fmt.Printf("type %s for help\n", CmdHelp)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		serverConnection := ServerConnection{connection: connection, currentFolder: ""}
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
	case CmdHelp:
		showHelp()

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

	case CmdList:
		err := c.listFiles()
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

	_, err = fmt.Fprintf(c.connection, "%s %s%s %d\n", CmdUpload, c.currentFolder, filepath.Base(filename), fileInfo.Size())
	if err != nil {
		return err
	}

	buffer := make([]byte, 8192)
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
	_, err := fmt.Fprintf(c.connection, "%s %s%s\n", CmdDownload, c.currentFolder, filename)
	if err != nil {
		return err
	}

	buffer := make([]byte, 8192)
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

func (c ServerConnection) listFiles() error {
	_, err := fmt.Fprintf(c.connection, "%s %s\n", CmdList, c.currentFolder)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(c.connection)
	statusMessage, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	if strings.HasPrefix(statusMessage, "Error") {
		return fmt.Errorf(statusMessage)
	}

	fileCount, err := strconv.Atoi(strings.TrimSpace(statusMessage))
	if err != nil {
		return err
	}

	fmt.Println("Files on server:")
	fmt.Println("----------------")
	for i := 0; i < fileCount; i++ {
		fileInfo, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		parts := strings.Split(fileInfo, "|")
		if len(parts) != 3 {
			continue
		}

		fileName := parts[0]
		fileSize, _ := strconv.Atoi(parts[1])
		isDir, _ := strconv.ParseBool(parts[2])

		if isDir {
			fmt.Printf("[DIR] %s\n", fileName)
		} else {
			fmt.Printf("%s (%d bytes)\n", fileName, fileSize)
		}
	}

	return nil
}

func (c ServerConnection) changeCurrentFolder() {

}

func showHelp() {
	fmt.Printf("Upload file: %s file\n", CmdUpload)
	fmt.Printf("Download file: %s file\n", CmdDownload)
}
