package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Microsoft/go-winio"
)

var (
	// Mutex để đảm bảo truy cập an toàn vào dữ liệu bí mật từ nhiều goroutine
	secretMutex sync.RWMutex
	// Biến toàn cục để lưu trữ dữ liệu đã được mã hóa
	storedEncryptedSecret []byte
)

func main() {
	fmt.Println("Bắt đầu dịch vụ credential...")

	// Chạy một goroutine để định kỳ lấy và cập nhật bí mật
	go updateSecretPeriodically()

	// Bắt đầu lắng nghe các kết nối IPC
	// Đây là một trình giữ chỗ cho cơ chế IPC thực tế.
	// Đối với Windows, named pipe là một lựa chọn phổ biến cho việc này.
	if err := startIPCListener(); err != nil {
		log.Fatalf("Không thể bắt đầu trình nghe IPC: %v", err)
	}
}

// startIPCListener bắt đầu lắng nghe các kết nối từ Credential Provider.
// Đây là một ví dụ đơn giản hóa.
// Đối với một Credential Provider thực tế của Windows, bạn có thể sẽ sử dụng Named Pipes.
func startIPCListener() error {
	pipePath := "\\\\.\\pipe\\MySecretServicePipe"
	listener, err := winio.ListenPipe(pipePath, nil)
	if err != nil {
		return fmt.Errorf("không thể lắng nghe trên named pipe: %w", err)
	}
	defer listener.Close()

	log.Printf("Trình nghe IPC đã bắt đầu tại %s", pipePath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Không thể chấp nhận kết nối: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

// handleConnection xử lý một kết nối IPC đến.
func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Đã chấp nhận kết nối.")

	// 1. Đọc yêu cầu từ Credential Provider
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Không thể đọc từ kết nối: %v", err)
		return
	}

	// Credential Providers thường gửi chuỗi UTF-16 (wide strings).
	// Khi đọc dưới dạng byte thô, chúng có thể chứa các byte null xen kẽ.
	// Chúng ta sẽ chuẩn hóa chuỗi bằng cách xóa các byte null và khoảng trắng.
	rawRequest := string(buf[:n])
	processedRequest := strings.ReplaceAll(rawRequest, "\x00", "")   // Xóa byte null
	processedRequest = strings.ReplaceAll(processedRequest, " ", "") // Xóa khoảng trắng
	processedRequest = strings.TrimSpace(processedRequest)

	log.Printf("Đã nhận được yêu cầu thô: '%s'", rawRequest)
	log.Printf("Yêu cầu đã xử lý: '%s'", processedRequest)

	// 2. Kiểm tra xem yêu cầu có phải là "GET_SECRET" không
	if processedRequest == "GET_SECRET" {
		log.Println("Yêu cầu 'GET_SECRET' hợp lệ. Đang xử lý...")

		secretMutex.RLock() // Khóa để đọc an toàn
		if storedEncryptedSecret == nil {
			secretMutex.RUnlock()
			log.Println("Lỗi: Bí mật chưa được khởi tạo.")
			conn.Write([]byte("LỖI: Bí mật chưa sẵn sàng"))
			return
		}

		// Giải mã dữ liệu được lưu trữ
		decryptedData, err := decryptData(storedEncryptedSecret)
		secretMutex.RUnlock()

		if err != nil {
			log.Printf("Không thể giải mã dữ liệu: %v", err)
			conn.Write([]byte("LỖI: giải mã thất bại"))
			return
		}

		// 3. Phản hồi cho Credential Provider với dữ liệu đã giải mã
		_, err = conn.Write(decryptedData)
		if err != nil {
			log.Printf("Không thể ghi vào kết nối: %v", err)
		}
		log.Println("Đã gửi dữ liệu bí mật thành công.")

	} else {
		log.Printf("Yêu cầu không hợp lệ nhận được: '%s'", processedRequest)
	}
}

// updateSecretPeriodically là một vòng lặp chạy nền để lấy dữ liệu từ máy chủ,
// mã hóa và lưu trữ nó.
func updateSecretPeriodically() {
	// Chạy ngay một lần khi khởi động
	log.Println("Cập nhật bí mật lần đầu...")
	updateSecret()

	// Sau đó chạy định kỳ mỗi 5 phút
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		log.Println("Đang cập nhật bí mật định kỳ...")
		updateSecret()
	}
}

// updateSecret thực hiện logic lấy, mã hóa và lưu trữ bí mật.
func updateSecret() {
	serverData, err := fetchDataFromServer()
	if err != nil {
		log.Printf("Lỗi khi lấy dữ liệu từ máy chủ: %v", err)
		return
	}

	encryptedData, err := encryptData(serverData)
	if err != nil {
		log.Printf("Lỗi khi mã hóa dữ liệu: %v", err)
		return
	}

	secretMutex.Lock() // Khóa để ghi an toàn
	storedEncryptedSecret = encryptedData
	secretMutex.Unlock()

	log.Println("Đã lấy và mã hóa bí mật mới thành công.")
}

// fetchDataFromServer mô phỏng việc lấy dữ liệu từ một máy chủ bên ngoài.
func fetchDataFromServer() ([]byte, error) {
	log.Println("Đang lấy dữ liệu từ máy chủ...")
	// Trong một triển khai thực tế, điều này sẽ tạo một yêu cầu HTTP hoặc tương tự.
	return []byte("888888"), nil
}

// encryptData là một trình giữ chỗ cho logic mã hóa của bạn.
func encryptData(data []byte) ([]byte, error) {
	log.Println("Đang mã hóa dữ liệu...")
	// QUAN TRỌNG: Sử dụng một thư viện mã hóa mạnh, sẵn sàng cho sản xuất như
	// crypto/aes, crypto/cipher, và golang.org/x/crypto/nacl/secretbox.
	// Đây chỉ là một trình giữ chỗ và KHÔNG an toàn.
	encrypted := make([]byte, len(data))
	for i, b := range data {
		encrypted[i] = b + 1 // "Mã hóa" đơn giản không an toàn
	}
	return encrypted, nil
}

// decryptData là một trình giữ chỗ cho logic giải mã của bạn.
func decryptData(data []byte) ([]byte, error) {
	log.Println("Đang giải mã dữ liệu...")
	// Điều này sẽ đảo ngược logic trong encryptData.
	decrypted := make([]byte, len(data))
	for i, b := range data {
		decrypted[i] = b - 1
	}
	return decrypted, nil
}
