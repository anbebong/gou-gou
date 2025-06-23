package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Microsoft/go-winio"
)

func main() {
	pipePath := "\\\\.\\pipe\\MySecretServicePipe"

	conn, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		log.Fatalf("Failed to connect to named pipe: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to the service.")

	// Send a request to the service
	request := []byte("GET_SECRET")
	_, err = conn.Write(request)
	if err != nil {
		log.Fatalf("Failed to write to pipe: %v", err)
	}

	fmt.Printf("Sent request: %s\n", request)

	// Read the response from the service
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read from pipe: %v", err)
	}

	response := buf[:n]
	fmt.Printf("Received response: %s\n", response)

	// You can add more logic here to handle different responses
	if string(response) == "ERROR: could not get data" {
		fmt.Println("The service reported an error.")
		os.Exit(1)
	}
}
