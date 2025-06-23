package main

import (
	"net"
	"time"
)

// Cấu trúc cho thông tin phần cứng
type HardwareInfo struct {
	HostID    string `json:"hostID"`
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

// Cấu trúc để lưu thông tin client đang hoạt động
type ActiveClient struct {
	ID   string
	Conn net.Conn
}

// Cấu trúc để lưu thông tin client đã đăng ký vào file
type RegisteredClientInfo struct {
	ClientForJSON
	Conn       net.Conn  `json:"-"` // Connection object, ignored in JSON
	LastSeen   time.Time `json:"lastSeen"`
	IsOnline   bool      `json:"isOnline"`
	RemoteAddr string    `json:"remoteAddr"`
}

// ClientForJSON là cấu trúc chuyên dụng để lưu và tải thông tin client từ file JSON.
type ClientForJSON struct {
	AgentID      string       `json:"agentID"`
	ClientID     string       `json:"clientID"`
	HardwareInfo HardwareInfo `json:"hardwareInfo"`
	Username     string       `json:"username"` // Thêm trường Username
}

// OTPInfo lưu trữ mã OTP và thời gian hết hạn
type OTPInfo struct {
	Code      string
	ExpiresAt time.Time
}

// GetOTPResponse định nghĩa cấu trúc cho phản hồi API OTP
type GetOTPResponse struct {
	ClientID         string `json:"clientID"`
	AgentID          string `json:"agentID"`
	OTP              string `json:"otp"`
	ExpiresInSeconds int64  `json:"expiresTime"` // Số giây còn lại cho OTP
}
