package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// startAPIServer khởi động một server HTTP cho việc quản trị.
func startAPIServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/clients", listClientsHandler)
	mux.HandleFunc("/clients/delete", deleteClientHandler)
	mux.HandleFunc("/message/send", sendMessageToClientHandler)

	InfoLogger.Println("API server đang khởi động trên cổng 8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("API server thất bại: %v", err)
	}
}

// listClientsHandler xử lý việc liệt kê tất cả các client đã đăng ký.
func listClientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Phương thức không được phép", http.StatusMethodNotAllowed)
		return
	}
	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()

	// Tạo một map mới để hiển thị, chỉ bao gồm AgentID và HardwareInfo
	displayClients := make(map[string]interface{})
	for _, client := range registeredClients {
		displayClients[client.AgentID] = client.HardwareInfo
	}

	data, err := json.MarshalIndent(displayClients, "", "  ")
	if err != nil {
		http.Error(w, "Không thể marshal clients", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
