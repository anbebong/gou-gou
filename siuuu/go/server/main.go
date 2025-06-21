package main

import (
	"log"
	"net"
)

func main() {
	setupLogging()
	loadRegisteredClients()

	// Chạy API server trong một goroutine
	go startAPIServer()

	// Chạy CLI trong một goroutine
	go startCLI()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Không thể bắt đầu server TCP: %v", err)
	}
	defer listener.Close()
	InfoLogger.Println("Server TCP đang lắng nghe trên cổng 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			ErrorLogger.Printf("Không thể chấp nhận kết nối TCP: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
