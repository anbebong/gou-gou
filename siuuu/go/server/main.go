package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	port     = flag.String("port", "8080", "Port to listen on")
	apiPort  = flag.String("apiport", "8081", "Port for the API server")
	logLevel = flag.String("loglevel", "info", "Log level (debug, info, warning, error)")
)

func main() {
	flag.Parse()

	setupLogging(*logLevel)
	loadRegisteredClients()

	// Chạy API server trong một goroutine
	go startAPIServer(*apiPort)

	// Chạy CLI trong một goroutine
	go startCLI()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("Không thể bắt đầu server TCP: %v", err)
	}
	defer listener.Close()
	InfoLogger.Println("Server TCP đang lắng nghe trên cổng", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			ErrorLogger.Printf("Không thể chấp nhận kết nối TCP: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
