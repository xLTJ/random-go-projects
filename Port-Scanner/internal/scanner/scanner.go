package scanner

import (
	"fmt"
	"net"
	"time"
)

type ScanResult struct {
	Port   int
	IsOpen bool
}

func RunScan(ports []int, targetHost string, resultChan chan ScanResult) {
	portChan := make(chan int, 100)

	// start workers
	for i := 0; i < cap(portChan); i++ {
		go worker(portChan, resultChan, targetHost)
	}

	// send ports to chan
	go func() {
		for _, port := range ports {
			portChan <- port
		}
		close(portChan)
	}()

	return
}

func worker(portChan chan int, resultChan chan ScanResult, targetHost string) {
	timeOut := time.Second

	for port := range portChan {
		address := fmt.Sprintf("%s:%d", targetHost, port)
		conn, err := net.DialTimeout("tcp", address, timeOut)
		if err != nil {
			resultChan <- ScanResult{Port: port, IsOpen: false}
			continue
		}
		_ = conn.Close()
		resultChan <- ScanResult{Port: port, IsOpen: true}
	}
}
