package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func runGoCodeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Nhận code từ client
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	code := string(msg)

	// Tạo file tạm
	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Không tạo được file tạm"))
		return
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(code)
	tmpFile.Close()

	// Chạy go run bằng exec.Command
	cmd := exec.Command("go", "run", tmpFile.Name())
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Giới hạn thời gian chạy (3s)
	if err := cmd.Start(); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Không thể thực thi mã Go"))
		return
	}
	done := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		}
		done <- struct{}{}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		// Một trong hai stream đã xong
	case <-time.After(3 * time.Second):
		cmd.Process.Kill()
		conn.WriteMessage(websocket.TextMessage, []byte("[TIMEOUT: quá 3 giây]"))
	}
	cmd.Wait()
	conn.WriteMessage(websocket.TextMessage, []byte("__DONE__"))
}

func formatGoCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	code, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Lỗi đọc dữ liệu"))
		return
	}
	cmd := exec.Command("gofmt")
	cmd.Stdin = bytes.NewReader(code)
	out, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Lỗi format: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(out)
}

func main() {
	// Serve static index.html
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", runGoCodeWS)
	http.HandleFunc("/format", formatGoCode)
	log.Println("Server đang chạy ở http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
