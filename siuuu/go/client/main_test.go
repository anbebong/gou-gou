package main

import (
	"bytes"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Microsoft/go-winio"
)

func TestEncryptDecrypt(t *testing.T) {
	originalData := []byte("this is some secret data")

	encryptedData, err := encryptData(originalData)
	if err != nil {
		t.Fatalf("encryptData failed: %v", err)
	}

	decryptedData, err := decryptData(encryptedData)
	if err != nil {
		t.Fatalf("decryptData failed: %v", err)
	}

	if !bytes.Equal(originalData, decryptedData) {
		t.Errorf("decrypted data does not match original data. got=%s, want=%s", decryptedData, originalData)
	}
}

func TestHandleConnection(t *testing.T) {
	// Create an in-memory pipe to simulate a network connection
	clientConn, serverConn := net.Pipe()

	// Run the connection handler in a separate goroutine
	go func() {
		handleConnection(serverConn)
	}()

	// --- Client Side ---

	// Write a dummy request to the connection to trigger the handler
	_, err := clientConn.Write([]byte("GET"))
	if err != nil {
		t.Fatalf("Failed to write to pipe: %v", err)
	}

	// Read the response from the handler
	buf := make([]byte, 1024)
	n, err := clientConn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	// The handler should fetch data, encrypt it, then decrypt it, and send it back.
	// So the final result should be the original data from the "server".
	expectedResponse, _ := fetchDataFromServer()
	actualResponse := buf[:n]

	if !bytes.Equal(actualResponse, expectedResponse) {
		t.Errorf("handler returned unexpected body: got %q want %q", actualResponse, expectedResponse)
	}
}

func TestNamedPipeConnection(t *testing.T) {
	// Run the listener in a separate goroutine
	go func() {
		// In a real app, you'd manage the lifecycle of this listener.
		// We expect this to block, the test will exit before that's a problem.
		if err := startIPCListener(); err != nil {
			// Since this is in a goroutine, we can't fail the test directly.
			// We can log the error. If the client fails to connect,
			// the test will fail anyway.
			log.Printf("startIPCListener failed: %v", err)
		}
	}()

	// Give the listener a moment to start up.
	time.Sleep(200 * time.Millisecond)

	// --- Client Side ---
	pipePath := "\\\\.\\pipe\\my-credential-pipe"
	conn, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		t.Fatalf("Client failed to connect to named pipe: %v", err)
	}
	defer conn.Close()

	// Write a dummy request to the connection
	_, err = conn.Write([]byte("GET"))
	if err != nil {
		t.Fatalf("Client failed to write to pipe: %v", err)
	}

	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Client failed to read from pipe: %v", err)
	}

	// Verify the response
	expectedResponse, _ := fetchDataFromServer()
	actualResponse := buf[:n]

	if !bytes.Equal(actualResponse, expectedResponse) {
		t.Errorf("handler returned unexpected body: got %q want %q", actualResponse, expectedResponse)
	}
}
