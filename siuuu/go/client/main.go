package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
	"github.com/denisbrodbeck/machineid"
	"github.com/hpcloud/tail"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	serviceName        = "AnPTClientSvc"
	serviceDisplayName = "AnPT Client Service"
	serviceDescription = "Handles OTP requests and provides them via IPC."
)

// Cấu trúc cho service
type myService struct{}

// Cấu trúc cho thông tin phần cứng
type HardwareInfo struct {
	HostID    string `json:"hostID"`
	HostName  string `json:"hostName"`
	IPAddress string `json:"ipAddress"`
}

// Cấu trúc cho tin nhắn
type Message struct {
	Type         string        `json:"type"` // "register", "message", "request_otp"
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
	aesKey          = []byte("1234567897654321") // 16 byte cho AES-128
	clientID        string                       // Lưu trữ ID của client (UUID)
	agentID         string                       // Lưu trữ AgentID ngắn
	configFile      string                       // Đường dẫn tuyệt đối đến file cấu hình
	logFileToWatch  string                       // Đường dẫn đến file log (tải từ config)
	clientLogFile   string                       // Đường dẫn file log của chính client này
	serverAddress   string                       // Địa chỉ server (tải từ config)
	serverConn      net.Conn                     // Kết nối TCP đến server
	otpResponseChan = make(chan string, 1)       // Channel để nhận OTP cho IPC
	logTail         *tail.Tail                   // Đối tượng tail để có thể dừng nó
	ipcListener     net.Listener                 // Đối tượng listener để có thể đóng nó
	logFileHandle   *os.File                     // Giữ file handle của file log để tránh bị đóng sớm
	exeDir          string                       // Thư mục chứa file thực thi
)

// Cấu trúc để lưu cấu hình client
type ClientConfig struct {
	ClientID      string `json:"clientID"`
	AgentID       string `json:"agentID"`
	ServerAddress string `json:"serverAddress"`
	LogFile       string `json:"logFile"`
	ClientLogFile string `json:"clientLogFile"`
}

// initLogger khởi tạo logger để ghi ra file.
// Nếu isDebug là true, nó sẽ ghi cả ra console.
func initLogger(isDebug bool) {
	// Tải cấu hình chỉ để lấy đường dẫn file log
	loadConfigForLog()

	var err error
	logFileHandle, err = os.OpenFile(clientLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Lỗi khi mở file log %s: %v", clientLogFile, err)
	}

	var writer io.Writer
	if isDebug {
		// Khi debug, ghi ra cả console và file
		writer = io.MultiWriter(os.Stdout, logFileHandle)
	} else {
		// Khi chạy như service, chỉ ghi ra file
		writer = logFileHandle
	}

	log.SetOutput(writer)
	log.Println("--------------------")
	// log.Println("Logger đã được khởi tạo.")
}

// loadConfigForLog là phiên bản rút gọn của loadConfig chỉ để lấy thông tin file log khi khởi tạo.
func loadConfigForLog() {
	// Đường dẫn mặc định giờ đây là tuyệt đối
	clientLogFile = filepath.Join(exeDir, "client.log")

	data, err := os.ReadFile(configFile)
	if err != nil {
		return // Dùng giá trị mặc định
	}

	var config ClientConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return // Dùng giá trị mặc định
	}

	if config.ClientLogFile != "" {
		// Nếu đường dẫn trong config không phải là tuyệt đối, hãy coi nó tương đối so với thư mục exe
		if !filepath.IsAbs(config.ClientLogFile) {
			clientLogFile = filepath.Join(exeDir, config.ClientLogFile)
		} else {
			clientLogFile = config.ClientLogFile
		}
	}
}

func main() {
	var err error
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Không thể lấy đường dẫn file thực thi: %v", err)
	}
	exeDir = filepath.Dir(exePath)
	configFile = filepath.Join(exeDir, "client_config.json")

	isDebug := len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "debug"
	initLogger(isDebug)

	if logFileHandle != nil {
		defer logFileHandle.Close()
	}

	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("Không thể xác định môi trường: %v", err)
	}

	// Nếu không chạy trong môi trường tương tác (console), thì chạy như một service
	if !isInteractive {
		log.Printf("Bắt đầu service '%s'", serviceName)
		if err = svc.Run(serviceName, &myService{}); err != nil {
			log.Printf("Service '%s' thất bại: %v", serviceName, err)
		}
		log.Printf("Service '%s' đã dừng.", serviceName)
		return
	}

	// Xử lý các lệnh từ command line
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-client <command>")
		fmt.Println("Commands: install, remove, start, stop, debug")
		return
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		log.Println("Chạy service ở chế độ debug...")
		// svc.Run sẽ gọi Execute, nơi logic chính được bắt đầu.
		// Logger đã được thiết lập để ghi ra console và file.
		if err = svc.Run(serviceName, &myService{}); err != nil {
			log.Fatalf("Debug run failed: %v", err)
		}
	case "install":
		err = installService(serviceName, serviceDisplayName, serviceDescription)
		if err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		fmt.Printf("Service '%s' installed successfully.\n", serviceDisplayName)
	case "remove":
		err = removeService(serviceName)
		if err != nil {
			log.Fatalf("Failed to remove service: %v", err)
		}
		fmt.Printf("Service '%s' removed successfully.\n", serviceName)
	case "start":
		err = startService(serviceName)
		if err != nil {
			log.Fatalf("Failed to start service: %v", err)
		}
		fmt.Println("Service started.")
	case "stop":
		err = controlService(serviceName, svc.Stop, svc.Stopped)
		if err != nil {
			log.Fatalf("Failed to stop service: %v", err)
		}
		fmt.Println("Service stopped.")
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}

// Execute là hàm chính của service, chứa vòng lặp xử lý các yêu cầu điều khiển
func (s *myService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	go runClientLogic(make(chan struct{}))

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	log.Println("Service has started successfully.")

	<-r // Chờ tín hiệu dừng
	log.Println("Service is stopping.")
	changes <- svc.Status{State: svc.StopPending}
	// Logic dọn dẹp sẽ được thực hiện bởi các defer trong runClientLogic
	return
}

// connectAndRegister xử lý việc kết nối, đăng ký và tự động kết nối lại.
func connectAndRegister() {
	for {
		log.Printf("Đang cố gắng kết nối đến server tại %s...", serverAddress)
		conn, err := net.DialTimeout("tcp", serverAddress, 15*time.Second) // Thêm timeout cho kết nối
		if err != nil {
			log.Printf("Kết nối thất bại: %v. Thử lại sau 1 giây.", err)
			time.Sleep(1 * time.Second)
			continue // Thử lại
		}

		log.Println("Kết nối thành công. Gán kết nối cho biến toàn cục.")
		serverConn = conn // Gán kết nối thành công cho biến toàn cục

		// Nếu chưa có ClientID, thực hiện đăng ký
		if clientID == "" {
			log.Println("Chưa có ClientID, đang thực hiện đăng ký...")
			if err := register(); err != nil {
				log.Printf("Đăng ký thất bại: %v. Đóng kết nối và thử lại.", err)
				serverConn.Close() // Đóng kết nối hỏng
				serverConn = nil
				time.Sleep(10 * time.Second)
				continue // Thử lại
			}
			log.Println("Đăng ký thành công.")
		} else {
			log.Printf("Đã có AgentID: %s. Bỏ qua đăng ký.", agentID)
		}

		// Nếu đến được đây, kết nối và đăng ký (nếu cần) đã thành công.
		// Bắt đầu đọc phản hồi. Hàm này sẽ block cho đến khi kết nối bị mất.
		readResponses()

		// Nếu readResponses() kết thúc, nghĩa là kết nối đã mất.
		log.Println("Mất kết nối với server. Đang chuẩn bị kết nối lại...")
		serverConn.Close() // Đóng kết nối cũ
		serverConn = nil   // Đặt lại biến toàn cục
		// Vòng lặp for sẽ tự động lặp lại và thử kết nối lại.
	}
}

// runClientLogic chứa toàn bộ logic hoạt động của client
func runClientLogic(done chan struct{}) {
	loadConfig()

	// Chạy logic kết nối trong một goroutine riêng để không block các tác vụ khác
	go connectAndRegister()

	// Các tác vụ khác có thể chạy song song
	go watchLogFile()
	log.Println("Chuẩn bị khởi tạo IPC listener...")
	go startIPCListener()

	log.Println("Client logic is running.")
	// Chờ tín hiệu dừng từ service
	<-done
	log.Println("Client logic is shutting down.")

	// Dọn dẹp khi dừng
	if serverConn != nil {
		serverConn.Close()
	}
	if ipcListener != nil {
		ipcListener.Close()
	}
	if logTail != nil {
		logTail.Stop()
	}
}

func requestOTP() error {
	if serverConn == nil {
		return fmt.Errorf("không thể yêu cầu OTP: chưa kết nối đến server")
	}
	msg := Message{Type: "request_otp", ClientID: clientID}
	return sendMessage(msg)
}

func register() error {
	if serverConn == nil {
		return fmt.Errorf("không thể đăng ký: chưa kết nối đến server")
	}
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
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(encryptedMsg)))
	if _, err := serverConn.Write(lenBytes); err != nil {
		return fmt.Errorf("lỗi gửi độ dài tin nhắn: %w", err)
	}
	if _, err = serverConn.Write([]byte(encryptedMsg)); err != nil {
		return fmt.Errorf("lỗi gửi tin nhắn đăng ký: %w", err)
	}
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(serverConn, lenBuf); err != nil {
		return fmt.Errorf("lỗi đọc độ dài phản hồi: %w", err)
	}
	length := binary.BigEndian.Uint32(lenBuf)
	buf := make([]byte, length)
	if _, err := io.ReadFull(serverConn, buf); err != nil {
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
		log.Printf("Đăng ký thành công với AgentID: %s", agentID)
		return nil
	}
	return fmt.Errorf("đăng ký thất bại: %s", resp.Message)
}

func sendMessage(msg Message) error {
	if serverConn == nil {
		return fmt.Errorf("không thể gửi tin nhắn: chưa kết nối đến server")
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	encryptedMsg, err := encrypt(string(jsonMsg))
	if err != nil {
		return err
	}
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(encryptedMsg)))
	if _, err := serverConn.Write(lenBytes); err != nil {
		return fmt.Errorf("lỗi gửi độ dài tin nhắn: %w", err)
	}
	_, err = serverConn.Write([]byte(encryptedMsg))
	return err
}

func readResponses() {
	for {
		lenBuf := make([]byte, 4)
		_, err := io.ReadFull(serverConn, lenBuf)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("Server đã đóng kết nối.")
			} else {
				log.Printf("Lỗi khi đọc độ dài phản hồi: %v", err)
			}
			// Quan trọng: return ở đây để thoát khỏi hàm và kích hoạt logic kết nối lại
			return
		}
		length := binary.BigEndian.Uint32(lenBuf)
		buf := make([]byte, length)
		_, err = io.ReadFull(serverConn, buf)
		if err != nil {
			log.Printf("Lỗi khi đọc nội dung phản hồi: %v", err)
			return
		}
		decryptedResp, err := decrypt(string(buf))
		if err != nil {
			log.Printf("Lỗi giải mã phản hồi: %v", err)
			continue
		}
		var resp Response
		if err := json.Unmarshal([]byte(decryptedResp), &resp); err != nil {
			log.Printf("Lỗi unmarshal JSON phản hồi: %v", err)
			continue
		}
		// Xử lý các loại phản hồi khác nhau
		switch resp.Status {
		case "otp_generated":
			otp := resp.Data
			log.Println("Đã nhận OTP từ server.")
			// Gửi OTP đến channel với timeout để tránh block vô hạn.
			select {
			case otpResponseChan <- otp:
				log.Println("Đã gửi OTP trực tiếp cho trình xử lý IPC.")
			case <-time.After(5 * time.Second):
				log.Println("Hết thời gian chờ gửi OTP đến IPC handler.")
			}
		case "message":
			log.Printf("Tin nhắn từ server: %s", resp.Data)
		}
	}
}

// getHardwareInfo lấy thông tin phần cứng của máy
func getHardwareInfo() (*HardwareInfo, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, fmt.Errorf("không thể lấy machine id: %w", err)
	}
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	return &HardwareInfo{HostID: id, HostName: hostName}, nil
}

// loadConfig tải cấu hình từ file, sử dụng giá trị mặc định nếu cần.
func loadConfig() {
	// Giá trị mặc định
	serverAddress = "localhost:8080"
	logFileToWatch = filepath.Join(exeDir, "events.log")
	clientLogFile = filepath.Join(exeDir, "client.log") // Mặc định cho file log của client

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Không tìm thấy file cấu hình '%s'. Đang tạo file với giá trị mặc định.", configFile)
		saveConfig() // Lưu clientID, agentID (rỗng) và các giá trị mặc định
		return
	}

	var config ClientConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("Lỗi đọc file cấu hình: %v. Sử dụng giá trị mặc định.", err)
	}

	// Ghi đè giá trị mặc định bằng giá trị từ file nếu chúng không rỗng
	if config.ServerAddress != "" {
		serverAddress = config.ServerAddress
	}
	if config.LogFile != "" {
		if !filepath.IsAbs(config.LogFile) {
			logFileToWatch = filepath.Join(exeDir, config.LogFile)
		} else {
			logFileToWatch = config.LogFile
		}
	}
	if config.ClientLogFile != "" {
		if !filepath.IsAbs(config.ClientLogFile) {
			clientLogFile = filepath.Join(exeDir, config.ClientLogFile)
		} else {
			clientLogFile = config.ClientLogFile
		}
	}
	clientID = config.ClientID
	agentID = config.AgentID

	// Ghi lại file để đảm bảo nó chứa tất cả các trường mới nhất
	saveConfig()
}

// saveConfig lưu cấu hình hiện tại vào file
func saveConfig() {
	config := ClientConfig{
		ClientID:      clientID,
		AgentID:       agentID,
		ServerAddress: serverAddress,
		LogFile:       logFileToWatch,
		ClientLogFile: clientLogFile,
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Printf("Lỗi khi tạo JSON để lưu cấu hình: %v", err)
		return
	}
	if err := os.WriteFile(configFile, data, 0644); err != nil {
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

// watchLogFile theo dõi file log và gửi các dòng mới đến server
func watchLogFile() {
	// Đảm bảo file tồn tại
	if _, err := os.Stat(logFileToWatch); os.IsNotExist(err) {
		log.Printf("File log '%s' không tồn tại, đang tạo file...", logFileToWatch)
		file, err := os.Create(logFileToWatch)
		if err != nil {
			log.Printf("Không thể tạo file log: %v", err)
			return
		}
		file.Close()
	}
	var err error
	logTail, err = tail.TailFile(logFileToWatch, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
		MustExist: true,
		Poll:      true,
	})
	if err != nil {
		log.Printf("Không thể bắt đầu giám sát file: %v", err)
		return
	}
	for line := range logTail.Lines {
		if line.Text != "" {
			msg := Message{Type: "message", ClientID: clientID, Data: line.Text}
			if err := sendMessage(msg); err != nil {
				log.Printf("Lỗi khi gửi tin nhắn: %v", err)
			}
		}
	}
}

// startIPCListener bắt đầu lắng nghe các kết nối từ Credential Provider.
func startIPCListener() {
	log.Println("Bắt đầu thực thi hàm startIPCListener.")
	pipePath := `\\.\pipe\MySecretServicePipe`
	// Xóa pipe cũ nếu tồn tại để tránh lỗi "address already in use"
	os.Remove(pipePath)

	// SDDL (Security Descriptor Definition Language) cho phép mọi người dùng (Everyone)
	// có toàn quyền truy cập (Generic All) vào pipe. Điều này là cần thiết để
	// ứng dụng của người dùng (chạy ở Session 1) có thể kết nối với service (chạy ở Session 0).
	config := &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;WD)",
	}

	var err error
	ipcListener, err = winio.ListenPipe(pipePath, config)
	if err != nil {
		log.Printf("Không thể lắng nghe trên named pipe: %v", err)
		return
	}
	defer ipcListener.Close()

	log.Printf("Trình nghe IPC đã bắt đầu thành công tại %s", pipePath)

	for {
		conn, err := ipcListener.Accept()
		if err != nil {
			// Nếu listener đã bị đóng, thoát khỏi vòng lặp
			if strings.Contains(err.Error(), "use of closed network connection") || err == io.ErrClosedPipe {
				log.Println("Trình nghe IPC đã dừng.")
				return
			}
			log.Printf("Không thể chấp nhận kết nối IPC: %v", err)
			continue
		}
		go handleIPCConnection(conn)
	}
}

// handleIPCConnection xử lý một kết nối IPC đến.
func handleIPCConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Printf("Không thể đọc từ kết nối IPC: %v", err)
		}
		return
	}
	processedRequest := strings.TrimSpace(strings.ReplaceAll(string(buf[:n]), "\x00", ""))
	if processedRequest == "GET_SECRET" {
		log.Println("Yêu cầu 'GET_SECRET' hợp lệ. Đang yêu cầu OTP mới từ server...")
		if serverConn == nil {
			log.Println("Không thể yêu cầu OTP: Không có kết nối đến server.")
			conn.Write([]byte("ERROR: Not connected to server"))
			return
		}
		if err := requestOTP(); err != nil {
			log.Printf("Lỗi khi gửi yêu cầu OTP đến server: %v", err)
			conn.Write([]byte("ERROR: Failed to request OTP from server"))
			return
		}
		select {
		case otp := <-otpResponseChan:
			log.Println("Đã nhận được OTP, đang gửi cho client IPC.")
			conn.Write([]byte(otp))
		case <-time.After(10 * time.Second):
			log.Println("Lỗi: Hết thời gian chờ phản hồi OTP từ server.")
			conn.Write([]byte("ERROR: Timeout waiting for OTP from server"))
		}
	} else {
		log.Printf("Yêu cầu không xác định: '%s'", processedRequest)
		conn.Write([]byte("ERROR: Unknown request"))
	}
}

// --- Các hàm tiện ích để quản lý service ---

func installService(name, displayName, desc string) error {
	exepath, err := os.Executable()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}
	s, err = m.CreateService(name, exepath, mgr.Config{
		DisplayName: displayName,
		Description: desc,
		StartType:   mgr.StartAutomatic,
	})
	if err != nil {
		return err
	}
	defer s.Close()
	return nil
}

func removeService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()
	return s.Delete()
}

// startService correctly starts the service
func startService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start()
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if time.Now().After(timeout) {
			return fmt.Errorf("timeout waiting for service to reach state %d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}
