package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time" // Thêm time để sử dụng trong ClientDisplayInfo

	"github.com/golang-jwt/jwt/v5"
)

var usersFile = filepath.Join(getExecDir(), "users.json")

// getExecDir trả về thư mục chứa file thực thi
func getExecDir() string {
	exePath, err := os.Executable()
	if err != nil {
		return "." // fallback
	}
	return filepath.Dir(exePath)
}

// Danh sách user mẫu
var users = []User{
	{Username: "admin", Password: "adminpass", Role: "admin"},
	{Username: "user1", Password: "user1pass", Role: "user"},
}

// Hàm xác thực user
func authenticateAPIUser(username, password string) *User {
	for _, u := range users {
		if u.Username == username && u.Password == password {
			return &u
		}
	}
	return nil
}

var jwtSecret = []byte("your_secret_key") // Đổi secret khi deploy

func generateJWT(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func jwtAuthMiddleware(next http.HandlerFunc, requireAdmin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if requireAdmin && claims["role"] != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// Truyền claims vào context để handler lấy ra được
		ctx := context.WithValue(r.Context(), "claims", claims)
		next(w, r.WithContext(ctx))
	}
}

// Endpoint đăng nhập trả về JWT
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	user := authenticateAPIUser(req.Username, req.Password)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token, err := generateJWT(user.Username, user.Role)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   token,
		"user": map[string]interface{}{
			"username": user.Username,
			"role":     user.Role,
			"email":    user.Email,
			"name":     user.Name,
		},
	})
}

// Middleware CORS cho tất cả các route
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingResponseWriter để lấy status code trả về
// (nên đặt ở đầu file hoặc gần các middleware)
type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Middleware log mọi request API vào InfoLogger
func apiLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: 200}
		start := time.Now()
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		InfoLogger.Printf("[API] %s %s - %d (%s)", r.Method, r.URL.Path, lrw.StatusCode, duration)
	})
}

// API lấy log archive dạng JSON, có thể lọc theo trường (level, agentID, ...)
func getArchiveLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Chỉ hỗ trợ GET", http.StatusMethodNotAllowed)
		return
	}
	logPath := filepath.Join(getExecDir(), "archiver.log")
	InfoLogger.Printf("[API] Đọc log archive từ %s", logPath)
	file, err := os.Open(logPath)
	if err != nil {
		ErrorLogger.Printf("[API] Không thể đọc file log: %v", err)
		http.Error(w, "Không thể đọc file log", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var logs []map[string]string
	scanner := bufio.NewScanner(file)
	levelFilter := r.URL.Query().Get("level")
	agentFilter := r.URL.Query().Get("agent")
	for scanner.Scan() {
		line := scanner.Text()
		// Ví dụ: 2025/06/24 10:41:05 001: 2025-06-24 10:41:05 [WARNING] - Message
		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 4 {
			continue
		}
		ts := parts[0] + " " + parts[1]
		agent := strings.TrimSuffix(parts[2], ":")
		rest := parts[3]
		// Tìm level và message
		level := ""
		msg := ""
		levelStart := strings.Index(rest, "[")
		levelEnd := strings.Index(rest, "]")
		if levelStart != -1 && levelEnd > levelStart {
			level = rest[levelStart+1 : levelEnd]
			msg = strings.TrimSpace(rest[levelEnd+1:])
		} else {
			msg = rest
		}
		logEntry := map[string]string{
			"timestamp": ts,
			"agentID":   agent,
			"level":     level,
			"message":   msg,
		}
		if (levelFilter == "" || strings.EqualFold(level, levelFilter)) &&
			(agentFilter == "" || agent == agentFilter) {
			logs = append(logs, logEntry)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// API lấy log của thiết bị đã đăng nhập
func getMyDeviceLogHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	registeredClientsMutex.Lock()
	defer registeredClientsMutex.Unlock()
	var agentID string
	for _, client := range registeredClients {
		if client.Username == username {
			agentID = client.AgentID
			break
		}
	}
	if agentID == "" {
		http.Error(w, "User chưa được gán thiết bị", http.StatusForbidden)
		return
	}

	logPath := filepath.Join(getExecDir(), "archiver.log")
	file, err := os.Open(logPath)
	if err != nil {
		http.Error(w, "Không thể đọc file log", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var logs []map[string]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 4 {
			continue
		}
		ts := parts[0] + " " + parts[1]
		agent := strings.TrimSuffix(parts[2], ":")
		if agent != agentID {
			continue
		}
		rest := parts[3]
		level := ""
		msg := ""
		levelStart := strings.Index(rest, "[")
		levelEnd := strings.Index(rest, "]")
		if levelStart != -1 && levelEnd > levelStart {
			level = rest[levelStart+1 : levelEnd]
			msg = strings.TrimSpace(rest[levelEnd+1:])
		} else {
			msg = rest
		}
		logEntry := map[string]string{
			"timestamp": ts,
			"agentID":   agent,
			"level":     level,
			"message":   msg,
		}
		logs = append(logs, logEntry)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// API lấy danh sách user
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Chỉ hỗ trợ GET", http.StatusMethodNotAllowed)
		return
	}
	if err := loadUsers(); err != nil {
		ErrorLogger.Printf("[API] Không thể load users: %v", err)
		http.Error(w, "Không thể load users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// startAPIServer khởi động một server HTTP cho việc quản trị.
func startAPIServer(apiPort string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", loginHandler)
	mux.HandleFunc("/api/clients", jwtAuthMiddleware(listClientsHandler, false))
	mux.HandleFunc("/api/clients/delete", jwtAuthMiddleware(deleteClientHandler, true))
	mux.HandleFunc("/message/send", jwtAuthMiddleware(sendMessageToClientHandler, true))
	mux.HandleFunc("/api/otp", jwtAuthMiddleware(handleGetOTP, false))
	mux.HandleFunc("/api/clients/assign-user", jwtAuthMiddleware(handleAssignUser, true))
	mux.HandleFunc("/api/users/create", jwtAuthMiddleware(createUserHandler, true))
	mux.HandleFunc("/api/users/change-password", jwtAuthMiddleware(changePasswordHandler, false))
	mux.HandleFunc("/api/users/update", jwtAuthMiddleware(updateUserHandler, true))
	mux.HandleFunc("/api/logs/archive", jwtAuthMiddleware(getArchiveLogHandler, true))
	mux.HandleFunc("/api/users", jwtAuthMiddleware(listUsersHandler, true))
	mux.HandleFunc("/api/users/update-info", jwtAuthMiddleware(updateUserInfoHandler, false))
	mux.HandleFunc("/api/logs/my-device", jwtAuthMiddleware(getMyDeviceLogHandler, false))

	InfoLogger.Printf("API server đang khởi động trên cổng %s...", apiPort)
	// Bọc toàn bộ mux bằng logging + CORS middleware
	if err := http.ListenAndServe(fmt.Sprintf(":%s", apiPort), apiLoggingMiddleware(corsMiddleware(mux))); err != nil {
		log.Fatalf("API server thất bại: %v", err)
	}
}

// Load users từ file JSON
func loadUsers() error {
	f, err := os.Open(usersFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&users)
}

// Save users ra file JSON
func saveUsers() error {
	f, err := os.Create(usersFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(users)
}

// Hàm khởi tạo, gọi loadUsers khi start API server
func init() {
	if err := loadUsers(); err != nil {
		log.Printf("Không thể load users: %v", err)
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

// Thêm API tạo user mới
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	for _, u := range users {
		if u.Username == req.Username {
			http.Error(w, "User đã tồn tại", http.StatusBadRequest)
			return
		}
	}
	users = append(users, req)
	saveUsers()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

// Đổi mật khẩu: user chỉ đổi được mật khẩu của chính mình, admin đổi được của bất kỳ ai
func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username    string `json:"username"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Lấy user từ JWT
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	jwtUsername := claims["username"].(string)
	jwtRole := claims["role"].(string)

	// Luôn load lại users từ file để đảm bảo dữ liệu mới nhất
	if err := loadUsers(); err != nil {
		http.Error(w, "Không thể load users", http.StatusInternalServerError)
		return
	}

	for i, u := range users {
		if u.Username == req.Username {
			if jwtRole == "admin" || jwtUsername == req.Username {
				if jwtRole == "admin" || u.Password == req.OldPassword {
					users[i].Password = req.NewPassword
					saveUsers()
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"status": "changed"})
					return
				} else {
					http.Error(w, "Sai mật khẩu cũ", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Không đủ quyền", http.StatusForbidden)
				return
			}
		}
	}
	http.Error(w, "Không tìm thấy user", http.StatusNotFound)
}

// Sửa thông tin user (chỉ đổi role, chỉ admin được phép)
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Lấy user từ JWT
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	jwtRole := claims["role"].(string)
	if jwtRole != "admin" {
		http.Error(w, "Chỉ admin được phép sửa quyền user", http.StatusForbidden)
		return
	}
	// Luôn load lại users từ file để đảm bảo dữ liệu mới nhất
	if err := loadUsers(); err != nil {
		http.Error(w, "Không thể load users", http.StatusInternalServerError)
		return
	}
	for i, u := range users {
		if u.Username == req.Username {
			users[i].Role = req.Role
			saveUsers()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
			return
		}
	}
	http.Error(w, "Không tìm thấy user", http.StatusNotFound)
}

// Sửa thông tin user (user tự đổi tên/email, admin đổi được cho bất kỳ ai, KHÔNG cho đổi username)
func updateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Lấy user từ JWT
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	jwtUsername := claims["username"].(string)
	jwtRole := claims["role"].(string)

	if err := loadUsers(); err != nil {
		http.Error(w, "Không thể load users", http.StatusInternalServerError)
		return
	}

	for i, u := range users {
		if u.Username == req.Username {
			if jwtRole == "admin" || jwtUsername == req.Username {
				// KHÔNG cho phép đổi username, chỉ đổi name và email
				users[i].Name = req.Name
				users[i].Email = req.Email
				saveUsers()
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
				return
			} else {
				http.Error(w, "Không đủ quyền", http.StatusForbidden)
				return
			}
		}
	}
	http.Error(w, "Không tìm thấy user", http.StatusNotFound)
}
