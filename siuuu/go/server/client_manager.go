package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// "Cơ sở dữ liệu" trong bộ nhớ và file
var (
	activeClients      = make(map[string]*ActiveClient) // Client đang kết nối
	activeClientsMutex sync.Mutex

	registeredClients      = make(map[string]*RegisteredClientInfo) // Tất cả client đã từng đăng ký
	registeredClientsMutex sync.Mutex
	dbFile                 = "clients.json"
)

// loadRegisteredClients tải danh sách client từ file JSON
func loadRegisteredClients() {
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			InfoLogger.Println("Không tìm thấy file DB, sẽ tạo file mới.")
			return
		}
		ErrorLogger.Printf("Lỗi khi đọc file DB: %v", err)
		return
	}

	if err := json.Unmarshal(data, &registeredClients); err != nil {
		ErrorLogger.Printf("Lỗi unmarshal DB JSON: %v", err)
	}
	InfoLogger.Printf("Đã tải %d client đã đăng ký.", len(registeredClients))
}

// saveRegisteredClients lưu danh sách client vào file JSON
func saveRegisteredClients() {
	data, err := json.MarshalIndent(registeredClients, "", "  ")
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
