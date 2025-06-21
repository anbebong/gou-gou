package main

import "net"

// Cấu trúc cho thông tin phần cứng
type HardwareInfo struct {
	HostID    string `json:"hostID"`
	IPAddress string `json:"ipAddress"`
}

// Cấu trúc cho tin nhắn
type Message struct {
	Type         string        `json:"type"` // "register", "message"
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
	AgentID      string       `json:"agentID"`
	HardwareInfo HardwareInfo `json:"hardwareInfo"`
}
