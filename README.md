# Project Index
This is primarily for me personally so i know where to look for stuff if i need to reference something

## General Concepts

### Network Interactions
*Projects involving network communication like TCP, UDP, sockets.*

#### TCP Server
*Listens for and handles incoming TCP connections.*

- [Tcp-reciever](./Tcp-reciever/) - (Depth: Basic)
- [TCP-chat-server](./TCP-chat-server/) - (Depth: Basic)
- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Basic)

#### UDP Server
*Listens for and handles incoming UDP packets.*

- [UDP-reciever](./UDP-reciever/) - (Depth: Basic)
- [UDP chat server](./UDP-chat-server/) - (Depth: Intermediate)

#### TCP Client
*Connects to and interacts with a TCP server.*

- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Basic)

#### UDP Client
*Sends and receives UDP packets to/from a server.*

- [UDP chat server](./UDP-chat-server/) - (Depth: Basic)

#### Custom Protocol
*Implementation of simple application-layer protocols over network sockets.*

- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Basic)

### Application Logic
*Projects implementing specific application-level features or state.*

#### Multi-User State Management
*Managing state for multiple concurrent users or clients.*

- [TCP-chat-server](./TCP-chat-server/) - (Depth: Intermediate)
- [UDP chat server](./UDP-chat-server/) - (Depth: Intermediate)

### File Operations
*Projects focused on manipulating files or the filesystem.*

#### File Transfer
*Sending or receiving files over a network connection.*

- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Intermediate)

### Security
*Projects demonstrating specific security techniques or countermeasures.*

#### Path Traversal Prevention
*Basic checks to prevent accessing files outside intended directories.*

- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Basic)

### Cryptography
*Projects implementing cryptographic algorithms or techniques.*

#### Classical Ciphers
*Implementation of historical ciphers like Caesar, Vigenere, etc.*

- [Caesar-Cipher-Tool](./Caesar-Cipher-Tool/) - (Depth: Basic)

### Command-Line Interface (CLI)
*Projects primarily focused on providing a command-line interface.*

#### Command/Flag Parsing
*Parsing command-line arguments, flags, and subcommands.*

- [Caesar-Cipher-Tool](./Caesar-Cipher-Tool/) - (Depth: Intermediate)

## Go-Specific Concepts

### Concurrency
*Projects utilizing Go's concurrency primitives.*

#### Goroutines for Client Handling
*Using a separate goroutine for each connected client (typically TCP).*

- [Tcp-reciever](./Tcp-reciever/) - (Depth: Basic)
- [TCP-chat-server](./TCP-chat-server/) - (Depth: Basic)
- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Basic)

#### Goroutines for Concurrent Tasks
*Using goroutines to perform multiple operations concurrently (e.g., I/O, processing).*

- [UDP chat server](./UDP-chat-server/) - (Depth: Basic)

#### Channel-Based Communication
*Using channels for synchronization or data passing between goroutines.*

- [TCP-chat-server](./TCP-chat-server/) - (Depth: Intermediate)
- [UDP chat server](./UDP-chat-server/) - (Depth: Intermediate)

### OS Interaction
*Projects interacting directly with the operating system.*

#### Signal Handling
*Handling OS signals, often for graceful shutdown.*

- [Tcp-reciever](./Tcp-reciever/) - (Depth: Basic)
- [UDP-reciever](./UDP-reciever/) - (Depth: Basic)
- [TCP-chat-server](./TCP-chat-server/) - (Depth: Basic)
- [UDP chat server](./UDP-chat-server/) - (Depth: Basic)
- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Basic)

### IO & Buffering
*Projects demonstrating specific Input/Output techniques.*

#### Bufio Usage
*Using buffered I/O (`bufio` package) for efficiency.*

- [File-Transfer-Server](./File-Transfer-Server/) - (Depth: Basic)

#### Stdin/Stdout Handling
*Reading from standard input and writing to standard output.*

- [Caesar-Cipher-Tool](./Caesar-Cipher-Tool/) - (Depth: Basic)

### External Libraries
*Projects demonstrating usage of specific third-party Go libraries.*

#### Cobra
*Using the Cobra library for building CLI applications.*

- [Caesar-Cipher-Tool](./Caesar-Cipher-Tool/) - (Depth: Intermediate)
