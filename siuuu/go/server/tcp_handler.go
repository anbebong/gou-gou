package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"
)

// var otps = make(map[string]*OTPInfo) // Bản đồ lưu trữ OTP cho từng client

// sendMessageToClient tìm một client đang hoạt động và gửi tin nhắn.
func sendMessageToClient(clientID, messageData string) error {
	activeClientsMutex.Lock()
	client, ok := activeClients[clientID]
	activeClientsMutex.Unlock()

	if !ok {
		// Tìm agentID để hiển thị lỗi cho thân thiện
		_, clientInfo, found := findClientByAnyID(clientID)
		if found {
			return fmt.Errorf("agent %s không được kết nối", clientInfo.AgentID)
		}
		return fmt.Errorf("client %s không được kết nối", clientID)
	}

	response := Response{
		Status:  "message", // Sử dụng một trạng thái rõ ràng cho tin nhắn từ server
		Data:    messageData,
		Message: "Tin nhắn từ server",
	}
	return sendResponse(client.Conn, response)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	InfoLogger.Println("Một client mới đã kết nối, đang chờ đăng ký...")
	var clientID string

	// Vòng lặp để đọc nhiều tin nhắn từ cùng một client
	for {
		// 1. Đọc 4 byte đầu tiên để lấy độ dài của tin nhắn
		lenBuf := make([]byte, 4)
		_, err := io.ReadFull(conn, lenBuf)
		if err != nil {
			if err == io.EOF {
				// Dọn dẹp khi client ngắt kết nối
				if clientID != "" {
					// Tìm agentID để ghi log cho đẹp
					_, clientInfo, found := findClientByAnyID(clientID)
					activeClientsMutex.Lock()
					delete(activeClients, clientID)
					activeClientsMutex.Unlock()
					if found {
						InfoLogger.Printf("Agent %s (ClientID: %s) đã ngắt kết nối.", clientInfo.AgentID, clientID)
					} else {
						WarningLogger.Printf("Client %s đã ngắt kết nối nhưng không tìm thấy thông tin đăng ký.", clientID)
					}
				} else {
					InfoLogger.Printf("Client tại %s đã ngắt kết nối trước khi đăng ký.", conn.RemoteAddr())
				}
				return
			}
			ErrorLogger.Printf("Lỗi khi đọc độ dài tin nhắn từ client %s: %v", conn.RemoteAddr(), err)
			return
		}

		// 2. Giải mã độ dài
		length := binary.BigEndian.Uint32(lenBuf)
		if length > 8192 { // Đặt giới hạn kích thước hợp lý để tránh tấn công DoS
			ErrorLogger.Printf("Lỗi: Gói tin từ client %s quá lớn: %d bytes", conn.RemoteAddr(), length)
			return // Đóng kết nối nếu gói tin quá lớn
		}

		// 3. Đọc toàn bộ nội dung tin nhắn dựa trên độ dài
		buf := make([]byte, length)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			ErrorLogger.Printf("Lỗi khi đọc nội dung tin nhắn từ client %s: %v", conn.RemoteAddr(), err)
			return
		}

		// Giải mã gói tin
		decryptedData, err := decrypt(string(buf))
		if err != nil {
			ErrorLogger.Printf("Lỗi giải mã từ %s: %v", conn.RemoteAddr(), err)
			continue
		}

		var msg Message
		if err := json.Unmarshal([]byte(decryptedData), &msg); err != nil {
			ErrorLogger.Printf("Lỗi unmarshal JSON từ %s: %v", conn.RemoteAddr(), err)
			continue
		}

		// Xử lý tin nhắn dựa trên loại
		switch msg.Type {
		case "request_otp":
			if msg.ClientID == "" {
				WarningLogger.Println("Yêu cầu OTP không có ClientID")
				continue
			}

			// Gọi hàm đã sửa để tạo và lưu OTP
			otp, err := generateAndSaveOTP(msg.ClientID)
			if err != nil {
				ErrorLogger.Printf("Lỗi khi tạo và lưu OTP cho %s: %v", msg.ClientID, err)
				response := Response{Status: "error", Message: "Lỗi hệ thống khi tạo OTP."}
				sendResponse(conn, response)
				continue
			}

			// Gửi OTP về cho client
			InfoLogger.Printf("Đã tạo và gửi OTP %s cho ClientID: %s", otp, msg.ClientID)
			response := Response{Status: "otp_generated", Data: otp, Message: "Mã OTP của bạn có hiệu lực trong 5 phút."}
			sendResponse(conn, response)

		case "register":
			if msg.HardwareInfo == nil || msg.HardwareInfo.HostID == "" {
				WarningLogger.Println("Yêu cầu đăng ký không có thông tin HostID.")
				response := Response{Status: "error", Message: "Thiếu thông tin HostID."}
				sendResponse(conn, response)
				continue
			}

			// Kiểm tra xem phần cứng đã tồn tại chưa
			registeredClientsMutex.Lock()
			for _, existingInfo := range registeredClients {
				// Sửa lỗi: Truy cập HostID thông qua struct lồng nhau ClientForJSON
				if existingInfo.ClientForJSON.HardwareInfo.HostID == msg.HardwareInfo.HostID {
					WarningLogger.Printf("Phần cứng đã được đăng ký với AgentID: %s (HostID: %s)", existingInfo.AgentID, msg.HardwareInfo.HostID)
					response := Response{Status: "error", Message: "Phần cứng này đã được đăng ký."}
					sendResponse(conn, response)
					registeredClientsMutex.Unlock()
					return // Đóng kết nối
				}
			}

			// Đăng ký client mới
			newID := uuid.New().String()
			clientIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
			agentID := generateNextAgentIDLocked()

			// Sửa lỗi: Khởi tạo struct RegisteredClientInfo với cấu trúc ClientForJSON lồng nhau
			newClientInfo := &RegisteredClientInfo{
				ClientForJSON: ClientForJSON{
					ClientID: newID,
					AgentID:  agentID,
					Username: "", // Mặc định username là rỗng khi đăng ký mới
					HardwareInfo: HardwareInfo{
						HostID:    msg.HardwareInfo.HostID,
						HostName:  msg.HardwareInfo.HostName,
						IPAddress: clientIP,
					},
				},
			}
			registeredClients[newID] = newClientInfo
			saveRegisteredClientsLocked() // Sửa lỗi: Gọi hàm không khóa để tránh deadlock
			registeredClientsMutex.Unlock()

			// Thêm vào danh sách active
			activeClientsMutex.Lock()
			activeClients[newID] = &ActiveClient{ID: newID, Conn: conn}
			activeClientsMutex.Unlock()

			clientID = newID // Lưu clientID cho kết nối này

			InfoLogger.Printf("Client từ IP %s đã đăng ký với AgentID %s (ClientID: %s, HostID: %s)", clientIP, agentID, newID, msg.HardwareInfo.HostID)
			// Gửi lại ID cho client
			response := Response{Status: "success", ClientID: newID, AgentID: agentID, Message: "Đăng ký thành công"}
			sendResponse(conn, response)

		case "message":
			if msg.ClientID == "" {
				WarningLogger.Println("Tin nhắn không có ClientID")
				continue
			}

			// Tìm thông tin client để lấy AgentID cho log
			_, clientInfo, found := findClientByAnyID(msg.ClientID)
			if !found {
				WarningLogger.Printf("Tin nhắn từ ClientID không xác định: %s", msg.ClientID)
				continue
			}

			InfoLogger.Printf("Đã nhận tin nhắn từ Agent %s", clientInfo.AgentID)
			archiveLog.Printf("%s: %s", clientInfo.AgentID, msg.Data)

			// Gửi phản hồi
			responseText := fmt.Sprintf("Server đã nhận được tin nhắn của bạn, Agent %s", clientInfo.AgentID)
			response := Response{Status: "success", Data: responseText}
			sendResponse(conn, response)

		default:
			WarningLogger.Printf("Loại tin nhắn không xác định: %s", msg.Type)
		}
	}
}

// sendResponse mã hóa và gửi phản hồi đến client
func sendResponse(conn net.Conn, resp Response) error {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ErrorLogger.Printf("Lỗi marshal JSON phản hồi: %v", err)
		return err
	}

	encryptedResp, err := encrypt(string(jsonResp))
	if err != nil {
		ErrorLogger.Printf("Lỗi mã hóa phản hồi: %v", err)
		return err
	}

	// Gửi độ dài của phản hồi trước
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(encryptedResp)))
	if _, err := conn.Write(lenBytes); err != nil {
		ErrorLogger.Printf("Lỗi gửi độ dài phản hồi: %v", err)
		return err
	}

	// Gửi nội dung phản hồi đã mã hóa
	_, err = conn.Write([]byte(encryptedResp))
	if err != nil {
		ErrorLogger.Printf("Lỗi gửi nội dung phản hồi: %v", err)
	}
	return err
}
