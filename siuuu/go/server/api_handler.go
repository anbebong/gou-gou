package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time" // Thêm time để sử dụng trong ClientDisplayInfo
)

// startAPIServer khởi động một server HTTP cho việc quản trị.
func startAPIServer(apiPort string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/clients", listClientsHandler)
	mux.HandleFunc("/clients/delete", deleteClientHandler)
	mux.HandleFunc("/message/send", sendMessageToClientHandler)
	mux.HandleFunc("/api/otp", handleGetOTP)
	mux.HandleFunc("/api/clients/assign-user", handleAssignUser) // Giữ lại handler gán user

	InfoLogger.Printf("API server đang khởi động trên cổng %s...", apiPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", apiPort), mux); err != nil {
		log.Fatalf("API server thất bại: %v", err)
	}
}

// listClientsHandler xử lý việc liệt kê tất cả các client đã đăng ký với thông tin chi tiết.
// Phản hồi API sẽ KHÔNG chứa ClientID và remoteAddr riêng biệt.
func listClientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Phương thức không được phép", http.StatusMethodNotAllowed)
		return
	}

	// Định nghĩa một struct chỉ dành cho việc hiển thị API, không chứa ClientID
	type ClientDisplayInfo struct {
		AgentID      string       `json:"agentID"`
		Username     string       `json:"username"`
		IsOnline     bool         `json:"isOnline"`
		LastSeen     time.Time    `json:"lastSeen"`
		HardwareInfo HardwareInfo `json:"hardwareInfo"` // IP Address sẽ nằm trong đây
	}

	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	clients := make([]ClientDisplayInfo, 0, len(registeredClients))

	for _, client := range registeredClients {
		// Sao chép HardwareInfo để tránh thay đổi dữ liệu gốc
		hwInfo := client.HardwareInfo

		displayInfo := ClientDisplayInfo{
			AgentID:      client.AgentID,
			Username:     client.Username,
			LastSeen:     client.LastSeen,
			IsOnline:     false,  // Mặc định là offline
			HardwareInfo: hwInfo, // Gán bản sao
		}

		// Kiểm tra trạng thái online và cập nhật IP nếu client đang kết nối
		activeClientsMutex.Lock()
		if activeClient, ok := activeClients[client.ClientID]; ok {
			displayInfo.IsOnline = true
			// Cập nhật địa chỉ IP hiện tại của client nếu đang online
			if activeClient.Conn != nil {
				displayInfo.HardwareInfo.IPAddress = activeClient.Conn.RemoteAddr().String()
			}
		}
		activeClientsMutex.Unlock()

		clients = append(clients, displayInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		http.Error(w, "Không thể marshal clients", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// deleteClientHandler xử lý việc xóa một client.
func deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Phương thức không được phép", http.StatusMethodNotAllowed)
		return
	}
	idToDelete := r.URL.Query().Get("id")
	if idToDelete == "" {
		http.Error(w, "Thiếu tham số 'id' (AgentID)", http.StatusBadRequest)
		return
	}

	fullClientID, clientInfo, found := findClientByAnyID(idToDelete)
	if !found {
		http.Error(w, fmt.Sprintf("Agent với ID '%s' không tìm thấy", idToDelete), http.StatusNotFound)
		return
	}

	registeredClientsMutex.Lock()
	delete(registeredClients, fullClientID)
	saveRegisteredClients()
	registeredClientsMutex.Unlock()

	InfoLogger.Printf("[API] Đã xóa Agent %s (ClientID: %s)", clientInfo.AgentID, fullClientID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Agent %s đã được xóa thành công", clientInfo.AgentID)
}

// sendMessageToClientHandler xử lý việc gửi tin nhắn đến một client.
func sendMessageToClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Phương thức không được phép", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		TargetID string `json:"targetID"`
		Data     string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request body không hợp lệ", http.StatusBadRequest)
		return
	}

	fullClientID, clientInfo, found := findClientByAnyID(req.TargetID)
	if !found {
		http.Error(w, fmt.Sprintf("Client với ID '%s' không tìm thấy", req.TargetID), http.StatusNotFound)
		return
	}

	if err := sendMessageToClient(fullClientID, req.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	InfoLogger.Printf("[API] Đã gửi tin nhắn đến Agent %s", clientInfo.AgentID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Tin nhắn đã được gửi đến Agent %s", clientInfo.AgentID)
}

// handleGetOTP xử lý yêu cầu lấy mã OTP.
// GET /api/otp -> trả về OTP của tất cả client.
// GET /api/otp?id=<clientID_or_agentID> -> trả về OTP của một client cụ thể.
func handleGetOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Chỉ hỗ trợ phương thức GET"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	clientIDParam := r.URL.Query().Get("id")

	if clientIDParam == "" {
		// Lấy OTP của tất cả client
		// Lấy danh sách tất cả các client đã đăng ký
		registeredClientsMutex.Lock()
		clientIDs := make([]string, 0, len(registeredClients))
		for id := range registeredClients {
			clientIDs = append(clientIDs, id)
		}
		registeredClientsMutex.Unlock()

		// Tạo OTP mới cho từng client
		for _, id := range clientIDs {
			_, err := generateAndSaveOTP(id)
			if err != nil {
				// Ghi log lỗi nhưng vẫn tiếp tục cho các client khác
				ErrorLogger.Printf("API failed to generate OTP for %s in get-all: %v", id, err)
			}
		}

		// Lấy lại toàn bộ danh sách OTP vừa tạo
		AllOTPs := getAllOTPs()
		json.NewEncoder(w).Encode(AllOTPs)
	} else {
		// Tìm clientID đầy đủ từ param (có thể là agentID)
		clientID, _, ok := findClientByAnyID(clientIDParam)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{Status: "error", Message: "Không tìm thấy client"})
			return
		}

		// Luôn tạo một mã OTP mới theo yêu cầu
		_, err := generateAndSaveOTP(clientID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorLogger.Printf("API failed to generate OTP for %s: %v", clientID, err)
			json.NewEncoder(w).Encode(Response{Status: "error", Message: "Lỗi khi tạo mã OTP"})
			return
		}

		// Lấy thông tin của mã OTP vừa tạo để trả về
		otpInfo, ok := getOTPForClient(clientID)
		if !ok {
			// Trường hợp này gần như không thể xảy ra nếu bước trên thành công
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Status: "error", Message: "Không thể lấy thông tin OTP sau khi tạo"})
			return
		}
		json.NewEncoder(w).Encode(otpInfo)
	}
}

// handleAssignUser xử lý việc gán người dùng cho một client.
// POST /api/clients/assign-user
// Body: {"id": "<clientID_or_agentID>", "username": "<username>"}
type AssignUserRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func handleAssignUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Chỉ hỗ trợ phương thức POST"})
		return
	}

	var req AssignUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Request body không hợp lệ"})
		return
	}

	if req.ID == "" || req.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "id và username là bắt buộc"})
		return
	}

	// Tìm clientID đầy đủ từ id được cung cấp
	clientID, clientInfo, ok := findClientByAnyID(req.ID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Không tìm thấy client"})
		return
	}

	// Gán user cho client
	if assignUserToClient(clientID, req.Username) {
		InfoLogger.Printf("Đã gán username '%s' cho client %s (ID: %s)", req.Username, req.ID, clientID)
		json.NewEncoder(w).Encode(Response{Status: "success", Message: fmt.Sprintf("Đã gán username '%s' cho agent %s", req.Username, clientInfo.AgentID)})
	} else {
		// Trường hợp này gần như không thể xảy ra nếu findClientByAnyID thành công
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Lỗi khi gán người dùng"})
	}
}
