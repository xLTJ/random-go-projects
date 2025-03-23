package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Message struct {
	sender  User
	content string
	time    time.Time
}

type User struct {
	addr     *net.UDPAddr
	username string
}

type Chatroom struct {
	users           map[string]User
	addOrModifyUser chan User
	removeUser      chan User
	messages        chan Message
}

func (c *Chatroom) startChatroom(listener *net.UDPConn) {
	for {
		select {
		case newUser := <-c.addOrModifyUser:
			user, exists := c.users[newUser.addr.String()]

			if exists {
				c.messages <- Message{content: user.username + " has changed their username to: " + newUser.username}
				user.username = newUser.username
			} else {
				c.messages <- Message{content: newUser.username + " has joined"}
				c.users[newUser.addr.String()] = newUser
			}

			fmt.Println("Added new user")

		case user := <-c.removeUser:
			delete(c.users, user.addr.String())
			fmt.Printf("Removed user: %s\n", user.username)

			c.messages <- Message{content: user.username + " has left"}

		case message := <-c.messages:
			fmt.Printf("A message was sent from %s: %s\n", message.sender.addr, message.content)
			formattedMessage := message.content + "\n"

			if message.sender.addr != nil {
				formattedMessage = message.sender.username + ": " + formattedMessage
			}

			for userAddr, user := range c.users {
				if userAddr == message.sender.addr.String() {
					continue
				}

				_, err := listener.WriteToUDP([]byte(formattedMessage), user.addr)
				if err != nil {
					log.Printf("Failed to send to %s: %v", userAddr, err)
				}
			}
		}
	}
}

func startServer() {
	addr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 1234,
	}

	listener, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}
	fmt.Printf("Chat server listening on %s\n", listener.LocalAddr().String())

	chatroom := Chatroom{
		users:           make(map[string]User),
		addOrModifyUser: make(chan User),
		removeUser:      make(chan User),
		messages:        make(chan Message, 10),
	}
	go chatroom.startChatroom(listener)
	go handleMessages(listener, chatroom)

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

func handleMessages(listener *net.UDPConn, chatroom Chatroom) {
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := listener.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		message := string(buffer[:n])
		switch {
		case strings.HasPrefix(message, "!SetUsername"):
			user := User{
				addr:     remoteAddr,
				username: strings.TrimSpace(strings.TrimPrefix(message, "!SetUsername")),
			}

			chatroom.addOrModifyUser <- user
			continue

		case strings.HasPrefix(message, "!Disconnect"):
			chatroom.removeUser <- chatroom.users[remoteAddr.String()]
		}

		chatroom.messages <- Message{
			sender:  chatroom.users[remoteAddr.String()],
			content: message,
			time:    time.Now(),
		}
	}
}
