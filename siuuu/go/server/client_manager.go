package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// "Cơ sở dữ liệu" trong bộ nhớ và file
var (
	activeClients      = make(map[string]*ActiveClient) // Client đang kết nối
	activeClientsMutex sync.Mutex

	registeredClients      = make(map[string]*RegisteredClientInfo) // Tất cả client đã từng đăng ký
	registeredClientsMutex sync.Mutex
	dbFile                 = "clients.json"

	otps      = make(map[string]*OTPInfo) // Thêm map để lưu trữ OTP
	otpsMutex sync.Mutex                  // Mutex để bảo vệ map OTP
)

// loadRegisteredClients tải danh sách client từ file JSON
func loadRegisteredClients() {
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	// Khởi tạo map để đảm bảo nó không nil
	registeredClients = make(map[string]*RegisteredClientInfo)

	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			InfoLogger.Println("Không tìm thấy file DB, sẽ tạo file mới.")
			return
		}
		ErrorLogger.Printf("Lỗi khi đọc file DB: %v", err)
		return
	}

	if len(data) == 0 {
		return // File rỗng, không cần làm gì
	}

	var clientsFromJSON []ClientForJSON
	if err := json.Unmarshal(data, &clientsFromJSON); err != nil {
		ErrorLogger.Printf("Lỗi unmarshal DB JSON: %v", err)
		return
	}

	for _, clientJSON := range clientsFromJSON {
		// Tạo một bản sao để tránh các vấn đề về con trỏ biến lặp
		c := clientJSON
		registeredClients[c.ClientID] = &RegisteredClientInfo{
			ClientForJSON: c,     // Sử dụng trực tiếp struct đã unmarshal
			IsOnline:      false, // Đặt giá trị runtime mặc định
		}
	}

	InfoLogger.Printf("Đã tải %d client đã đăng ký.", len(registeredClients))
}

// saveRegisteredClients lưu danh sách client vào file JSON một cách an toàn (thread-safe)
func saveRegisteredClients() {
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()
	saveRegisteredClientsLocked()
}

// saveRegisteredClientsLocked chứa logic lưu file, yêu cầu caller phải giữ mutex
func saveRegisteredClientsLocked() {
	var clientsToSave []ClientForJSON
	for _, clientInfo := range registeredClients {
		// Struct ClientForJSON đầy đủ đã có trong clientInfo, bao gồm cả Username
		clientsToSave = append(clientsToSave, clientInfo.ClientForJSON)
	}

	data, err := json.MarshalIndent(clientsToSave, "", "  ")
	if err != nil {
		ErrorLogger.Printf("Lỗi marshal DB JSON: %v", err)
		return
	}

	if err := os.WriteFile(dbFile, data, 0644); err != nil {
		ErrorLogger.Printf("Lỗi khi ghi file DB: %v", err)
	}
}

// generateNextAgentIDLocked tạo ra một ID agent mới dựa trên ID lớn nhất hiện có.
// Caller phải giữ registeredClientsMutex.
func generateNextAgentIDLocked() string {
	maxID := 0
	for _, client := range registeredClients {
		if client.AgentID != "" {
			id, err := strconv.Atoi(client.AgentID)
			if err == nil && id > maxID {
				maxID = id
			}
		}
	}
	return fmt.Sprintf("%03d", maxID+1)
}

// assignUserToClient gán một username cho một client.
func assignUserToClient(clientID, username string) bool {
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	if clientInfo, ok := registeredClients[clientID]; ok {
		clientInfo.Username = username
		saveRegisteredClientsLocked() // Gọi phiên bản không khóa để tránh deadlock
		return true
	}
	return false
}

// findClientByAnyID tìm một client bằng ClientID (uuid) hoặc AgentID (số ngắn).
func findClientByAnyID(id string) (string, *RegisteredClientInfo, bool) {
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	// Đầu tiên, kiểm tra xem id có phải là ClientID đầy đủ không
	if clientInfo, ok := registeredClients[id]; ok {
		return id, clientInfo, true
	}

	// Nếu không, lặp để tìm bằng AgentID
	for clientID, clientInfo := range registeredClients {
		if clientInfo.AgentID == id {
			return clientID, clientInfo, true
		}
	}

	return "", nil, false
}

// Tạo struct OTPInfo để lưu trữ thông tin OTP

// generateAndSaveOTP tạo, lưu và trả về một mã OTP mới cho một client.
func generateAndSaveOTP(clientID string) (string, error) {
	// Tạo secret từ clientID theo thuật toán đã định
	otpSecret := generateOTPSecretFromClientID(clientID)

	// Tạo mã TOTP từ secret
	otp, err := generateTOTP(otpSecret)
	if err != nil {
		return "", fmt.Errorf("không thể tạo mã TOTP: %w", err)
	}

	// Tính toán thời gian hết hạn thực sự của TOTP (kết thúc của khung thời gian 30s hiện tại)
	now := time.Now()
	step := int64(30) // Khung thời gian 30 giây
	remainingSeconds := step - (now.Unix() % step)
	expiresAt := now.Add(time.Duration(remainingSeconds) * time.Second)

	otpsMutex.Lock()
	// Lưu OTP với thời gian hết hạn chính xác
	otps[clientID] = &OTPInfo{
		Code:      otp,
		ExpiresAt: expiresAt,
	}
	otpsMutex.Unlock()

	return otp, nil
}

// removeOTP xóa OTP khỏi map
func removeOTP(clientID string) {
	otpsMutex.Lock()
	delete(otps, clientID)
	otpsMutex.Unlock()
}

// getAllOTPs trả về thông tin OTP của tất cả các client hiện có.
func getAllOTPs() []GetOTPResponse {
	otpsMutex.Lock()
	defer otpsMutex.Unlock()
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	// Khởi tạo một slice rỗng, không phải nil, để đảm bảo JSON trả về là [] thay vì null
	result := make([]GetOTPResponse, 0)

	for clientID, otpInfo := range otps {
		if clientInfo, ok := registeredClients[clientID]; ok {
			remainingSeconds := int64(time.Until(otpInfo.ExpiresAt).Seconds())
			if remainingSeconds < 0 {
				remainingSeconds = 0
			}
			result = append(result, GetOTPResponse{
				ClientID:         clientID,
				AgentID:          clientInfo.AgentID,
				OTP:              otpInfo.Code,
				ExpiresInSeconds: remainingSeconds,
			})
		}
	}
	return result
}

// getOTPForClient trả về thông tin OTP cho một client cụ thể.
func getOTPForClient(clientID string) (*GetOTPResponse, bool) {
	otpsMutex.Lock()
	defer otpsMutex.Unlock()
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	otpInfo, ok := otps[clientID]
	if !ok {
		return nil, false
	}

	if clientInfo, ok := registeredClients[clientID]; ok {
		remainingSeconds := int64(time.Until(otpInfo.ExpiresAt).Seconds())
		if remainingSeconds < 0 {
			remainingSeconds = 0
		}
		return &GetOTPResponse{
			ClientID:         clientID,
			AgentID:          clientInfo.AgentID,
			OTP:              otpInfo.Code,
			ExpiresInSeconds: remainingSeconds,
		}, true
	}

	return nil, false // Không tìm thấy thông tin client đã đăng ký
}
