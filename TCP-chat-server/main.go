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

type User struct {
	username         string
	connection       net.Conn
	incomingMessages chan string
}

type Chatroom struct {
	users      map[*User]bool
	addUser    chan *User
	removeUser chan *User
	messages   chan string
}

func (c *Chatroom) StartChatroom() {
	for {
		select {
		case user := <-c.addUser:
			c.users[user] = true
			fmt.Printf("Added new user\n")

		case user := <-c.removeUser:
			delete(c.users, user)
			close(user.incomingMessages)
			fmt.Printf("Removed user: %s\n", user.username)

			for otherUser := range c.users {
				otherUser.incomingMessages <- user.username + " has left"
			}

		case message := <-c.messages:
			fmt.Printf("A message was sent: %s\n", message)
			for user := range c.users {
				user.incomingMessages <- message
			}
		}
	}
}

func main() {
	// Setup listener to listen on random port
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		log.Fatalf("Failed to setup listener: %v\n", err)
	}
	fmt.Printf("Chat Server listening on: %s\n", listener.Addr().String())

	// Setup chatroom and run it in a goroutine
	chatroom := Chatroom{
		users:      make(map[*User]bool),
		addUser:    make(chan *User),
		removeUser: make(chan *User),
		messages:   make(chan string),
	}
	go chatroom.StartChatroom()

	go func() {
		for {
			// Wait for a user to connect
			connection, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v\n", err)
				continue
			}

			// Setup user and handle them in a goroutine so we can handle any other new users joining
			user := User{
				username:         "Anonymous",
				connection:       connection,
				incomingMessages: make(chan string),
			}
			go handleUser(&user, &chatroom)
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

func handleUser(user *User, chatroom *Chatroom) {
	chatroom.addUser <- user

	go sendMessages(user, chatroom)

	viewMessages(user)
}

func sendMessages(user *User, chatroom *Chatroom) {
	reader := bufio.NewReader(user.connection)

	// The first message a user sends is their username
	user.incomingMessages <- "Enter username: "
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading username: %v\n", err)
	} else {
		user.username = username[:len(username)-1]
	}

	chatroom.messages <- user.username + " has joined"

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			chatroom.removeUser <- user
			break
		}

		chatroom.messages <- user.username + ": " + message[:len(message)-1]
	}
}

func viewMessages(user *User) {
	for {
		message := <-user.incomingMessages
		_, err := user.connection.Write([]byte(message + "\n"))
		if err != nil {
			break
		}
	}
}
