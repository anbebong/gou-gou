package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary" // Thêm import
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/hpcloud/tail"
)

// Khóa mã hóa (phải giống với khóa của server)
var aesKey = []byte("1234567897654321") // 16 byte cho AES-128

// Cấu trúc cho thông tin phần cứng
type HardwareInfo struct {
	HostID    string `json:"hostID"`
	IPAddress string `json:"ipAddress"`
}

// Cấu trúc cho tin nhắn
type Message struct {
	Type         string        `json:"type"`
	ClientID     string        `json:"clientID,omitempty"`
	Data         string        `json:"data,omitempty"`
	HardwareInfo *HardwareInfo `json:"hardwareInfo,omitempty"`
}

// Cấu trúc cho phản hồi
type Response struct {
	Status   string `json:"status"`
	ClientID string `json:"clientID,omitempty"`
	AgentID  string `json:"agentID,omitempty"`
	Data     string `json:"data,omitempty"`
	Message  string `json:"message,omitempty"`
}

var (
	clientID       string // Lưu trữ ID của client (UUID)
	agentID        string // Lưu trữ AgentID ngắn
	configFle      = "client_config.json"
	logFileToWatch = "events.log"
)

// Cấu trúc để lưu cấu hình client
type ClientConfig struct {
	ClientID string `json:"clientID"`
	AgentID  string `json:"agentID"`
}

func main() {
	loadConfig()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Không thể kết nối đến server: %v", err)
	}
	defer conn.Close()

	log.Println("Đã kết nối đến server.")

	// Nếu chưa có ClientID, thực hiện đăng ký
	if clientID == "" {
		log.Println("Chưa có ClientID, đang thực hiện đăng ký...")
		if err := register(conn); err != nil {
			log.Fatalf("Đăng ký thất bại: %v", err)
		}
	} else {
		log.Printf("Sử dụng AgentID đã có: %s (ClientID: %s)", agentID, clientID)
	}

	// Chạy một goroutine để đọc các phản hồi từ server
	go readResponses(conn)

	// Bắt đầu giám sát file log
	watchLogFile(conn)
}

// watchLogFile theo dõi file log và gửi các dòng mới đến server
func watchLogFile(conn net.Conn) {
	// Đảm bảo file tồn tại
	if _, err := os.Stat(logFileToWatch); os.IsNotExist(err) {
		log.Printf("File log '%s' không tồn tại, đang tạo file...", logFileToWatch)
		file, err := os.Create(logFileToWatch)
		if err != nil {
			log.Fatalf("Không thể tạo file log: %v", err)
		}
		file.Close()
	}

	log.Printf("Đang giám sát file log: %s", logFileToWatch)
	t, err := tail.TailFile(logFileToWatch, tail.Config{
		ReOpen:    true,                                          // Mở lại file nếu nó bị xoay vòng hoặc xóa
		Follow:    true,                                          // Tiếp tục theo dõi file
		Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // Bắt đầu đọc từ cuối file
		MustExist: true,                                          // Báo lỗi nếu file không tồn tại
		Poll:      true,                                          // Sử dụng polling để phát hiện thay đổi, tốt hơn cho một số hệ thống file
	})
	if err != nil {
		log.Fatalf("Không thể bắt đầu giám sát file: %v", err)
	}

	for line := range t.Lines {
		if line.Text != "" {
			log.Printf("Phát hiện dòng mới: %s", line.Text)
			msg := Message{Type: "message", ClientID: clientID, Data: line.Text}
			if err := sendMessage(conn, msg); err != nil {
				log.Printf("Lỗi khi gửi tin nhắn: %v", err)
				// Cân nhắc việc thử lại hoặc xử lý lỗi ở đây
			}
		}
	}
}

// register gửi yêu cầu đăng ký đến server
func register(conn net.Conn) error {
	hwInfo, err := getHardwareInfo()
	if err != nil {
		return fmt.Errorf("không thể lấy thông tin phần cứng: %w", err)
	}

	msg := Message{Type: "register", HardwareInfo: hwInfo}
	jsonMsg, _ := json.Marshal(msg)

	encryptedMsg, err := encrypt(string(jsonMsg))
	if err != nil {
		return err
	}

	// Gửi độ dài của tin nhắn trước
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(encryptedMsg)))
	if _, err := conn.Write(lenBytes); err != nil {
		return fmt.Errorf("lỗi gửi độ dài tin nhắn: %w", err)
	}

	// Gửi tin nhắn
	_, err = conn.Write([]byte(encryptedMsg))
	if err != nil {
		return fmt.Errorf("lỗi gửi tin nhắn đăng ký: %w", err)
	}

	// Đợi phản hồi đăng ký
	// Đọc 4 byte đầu tiên để lấy độ dài của tin nhắn
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		if err == io.EOF {
			return fmt.Errorf("server đã đóng kết nối khi đang chờ phản hồi đăng ký")
		}
		return fmt.Errorf("lỗi đọc độ dài phản hồi: %w", err)
	}
	length := binary.BigEndian.Uint32(lenBuf)

	// Đọc toàn bộ tin nhắn dựa trên độ dài đã nhận
	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return fmt.Errorf("lỗi đọc nội dung phản hồi: %w", err)
	}

	decryptedResp, err := decrypt(string(buf))
	if err != nil {
		return fmt.Errorf("lỗi giải mã phản hồi đăng ký: %w", err)
	}

	var resp Response
	if err := json.Unmarshal([]byte(decryptedResp), &resp); err != nil {
		return fmt.Errorf("lỗi unmarshal JSON phản hồi đăng ký: %w", err)
	}

	if resp.Status == "success" && resp.ClientID != "" {
		clientID = resp.ClientID
		agentID = resp.AgentID // Lưu lại AgentID
		saveConfig()
		log.Printf("Đăng ký thành công với AgentID: %s (ClientID: %s)", agentID, clientID)
		return nil
	} else {
		return fmt.Errorf("đăng ký thất bại: %s", resp.Message)
	}
}

// sendMessage mã hóa và gửi một tin nhắn đến server
func sendMessage(conn net.Conn, msg Message) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	encryptedMsg, err := encrypt(string(jsonMsg))
	if err != nil {
		return err
	}

	// Gửi độ dài của tin nhắn trước
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(encryptedMsg)))
	if _, err := conn.Write(lenBytes); err != nil {
		return fmt.Errorf("lỗi gửi độ dài tin nhắn: %w", err)
	}

	// Gửi tin nhắn
	_, err = conn.Write([]byte(encryptedMsg))
	return err
}

// readResponses đọc và giải mã các phản hồi từ server
func readResponses(conn net.Conn) {
	for {
		// Đọc 4 byte đầu tiên để lấy độ dài của tin nhắn
		lenBuf := make([]byte, 4)
		_, err := io.ReadFull(conn, lenBuf)
		if err != nil {
			if err == io.EOF {
				log.Println("Server đã đóng kết nối.")
				os.Exit(0)
			}
			log.Printf("Lỗi khi đọc độ dài phản hồi từ server: %v", err)
			return
		}

		length := binary.BigEndian.Uint32(lenBuf)
		if length > 4096 { // Thêm một giới hạn để tránh cấp phát bộ nhớ quá lớn
			log.Printf("Lỗi: Gói tin nhận được quá lớn: %d bytes", length)
			return
		}

		// Đọc toàn bộ tin nhắn dựa trên độ dài đã nhận
		buf := make([]byte, length)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			log.Printf("Lỗi khi đọc nội dung phản hồi từ server: %v", err)
			return
		}

		decryptedResp, err := decrypt(string(buf))
		if err != nil {
			log.Printf("Lỗi giải mã phản hồi: %v", err)
			continue // Bỏ qua tin nhắn lỗi và tiếp tục
		}

		var resp Response
		if err := json.Unmarshal([]byte(decryptedResp), &resp); err != nil {
			log.Printf("Lỗi unmarshal JSON phản hồi: %v. Dữ liệu đã giải mã: %s", err, decryptedResp)
			continue
		}

		// Xử lý phản hồi
		if resp.Message != "" {
			log.Printf("Phản hồi từ server: %s", resp.Message)
		}
		if resp.Data != "" {
			log.Printf("Dữ liệu từ server: %s", resp.Data)
		}
	}
}

// getHardwareInfo lấy thông tin phần cứng (HostID)
func getHardwareInfo() (*HardwareInfo, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, fmt.Errorf("không thể lấy machine id: %w", err)
	}
	return &HardwareInfo{HostID: id}, nil
}

// loadConfig tải ClientID từ file cấu hình
func loadConfig() {
	data, err := os.ReadFile(configFle)
	if err != nil {
		if os.IsNotExist(err) {
			return // Không có file, không sao cả
		}
		log.Printf("Lỗi khi đọc file cấu hình: %v", err)
		return
	}

	var config ClientConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("Lỗi unmarshal JSON cấu hình: %v", err)
		return
	}
	clientID = config.ClientID
	agentID = config.AgentID
}

// saveConfig lưu ClientID và AgentID vào file cấu hình
func saveConfig() {
	config := ClientConfig{ClientID: clientID, AgentID: agentID}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Printf("Lỗi marshal JSON cấu hình: %v", err)
		return
	}

	if err := os.WriteFile(configFle, data, 0644); err != nil {
		log.Printf("Lỗi khi ghi file cấu hình: %v", err)
	}
}

// encrypt mã hóa văn bản bằng AES
func encrypt(text string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt giải mã văn bản bằng AES
func decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext quá ngắn")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
