package agent

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gou-pc/internal/crypto"
	"io"
	"net"
	"os"
	"time"

	"gou-pc/internal/logutil"

	"github.com/denisbrodbeck/machineid"
)

const (
	TypeRegister   = "register"
	TypeRequestOTP = "request_otp"
	TypeHello      = "hello"
	TypeLog        = "log"
	TypeError      = "error"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type Agent struct {
	Conn net.Conn
}

type DeviceInfo struct {
	HostName   string `json:"hostName"`
	IPAddress  string `json:"ipAddress"`
	MacAddress string `json:"macAddress"`
	HardwareID string `json:"hardwareID"`
}

// Chuẩn hoá struct cho mọi message trao đổi (ngoại trừ đăng ký): luôn có AgentID
// Dùng cho log, hello, request_otp, ...
type AgentMessageData struct {
	AgentID string      `json:"agent_id"`
	Payload interface{} `json:"payload,omitempty"`
}

// LogData dùng cho bản tin log
// (có thể dùng AgentMessageData.Payload = LogData)
type LogData struct {
	Message string `json:"message"`
}

func (a *Agent) Connect(addr string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	a.Conn = conn
	return nil
}

func (a *Agent) Send(msg Message) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		logutil.Error("Send: JSON marshal failed: %v", err)
		return err
	}
	logutil.Info("Send: {type:%s, agent_id:%s, payload:%v}", msg.Type, getAgentIDFromMsg(msg), getPayloadFromMsg(msg))
	encryptedMsg, err := crypto.Encrypt(string(jsonMsg))
	if err != nil {
		logutil.Error("Send: Encryption failed: %v", err)
		return err
	}
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(encryptedMsg)))
	if _, err := a.Conn.Write(lenBytes); err != nil {
		logutil.Error("Send: Write length failed: %v", err)
		return err
	}
	_, err = a.Conn.Write([]byte(encryptedMsg))
	if err != nil {
		logutil.Error("Send: Write encrypted message failed: %v", err)
	}
	return err
}

func (a *Agent) Receive() (Message, error) {
	var msg Message
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(a.Conn, lenBuf); err != nil {
		return msg, err
	}
	length := binary.BigEndian.Uint32(lenBuf)
	buf := make([]byte, length)
	if _, err := io.ReadFull(a.Conn, buf); err != nil {
		return msg, err
	}
	decryptedResp, err := crypto.Decrypt(string(buf))
	if err != nil {
		return msg, err
	}
	if err := json.Unmarshal([]byte(decryptedResp), &msg); err != nil {
		return msg, err
	}
	logutil.Info("Received: {type:%s, agent_id:%s, data:%v}", msg.Type, getAgentIDFromMsg(msg), msg.Data)
	return msg, nil
}

func (a *Agent) Close() error {
	if a.Conn != nil {
		return a.Conn.Close()
	}
	return nil
}

func GetDeviceInfo() (*DeviceInfo, error) {
	host, _ := os.Hostname()
	ip := getLocalIP()
	mac := getMacAddress()
	hwid, _ := machineid.ID()
	return &DeviceInfo{
		HostName:   host,
		IPAddress:  ip,
		MacAddress: mac,
		HardwareID: hwid,
	}, nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

func getMacAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String()
		}
	}
	return ""
}

func RegisterAgent(a *Agent, configPath string) (clientID, agentID string, err error) {
	dev, _ := GetDeviceInfo()
	logutil.Info("RegisterAgent: Registering device info: %+v", dev)
	msg := Message{Type: TypeRegister, Data: dev}
	if err := a.Send(msg); err != nil {
		logutil.Error("RegisterAgent: Send register message failed: %v", err)
		return "", "", err
	}
	resp, err := a.Receive()
	if err != nil {
		logutil.Error("RegisterAgent: Receive response failed: %v", err)
		return "", "", err
	}
	logutil.Info("RegisterAgent: Received response type=%s data=%v", resp.Type, resp.Data)
	if resp.Type == TypeRegister {
		var regInfo struct {
			ClientID string `json:"client_id"`
			AgentID  string `json:"agent_id"`
		}
		b, _ := json.Marshal(resp.Data)
		_ = json.Unmarshal(b, &regInfo)
		logutil.Info("RegisterAgent: Registration success client_id=%s agent_id=%s", regInfo.ClientID, regInfo.AgentID)
		_ = os.WriteFile(configPath, []byte(fmt.Sprintf(`{"client_id":"%s","agent_id":"%s"}`, regInfo.ClientID, regInfo.AgentID)), 0644)
		return regInfo.ClientID, regInfo.AgentID, nil
	}
	logutil.Error("RegisterAgent: Registration failed, response: %v", resp.Data)
	return "", "", fmt.Errorf("đăng ký thất bại: %v", resp.Data)
}

// WatchLogAndSend theo dõi file log, gửi dòng mới cho server
func (a *Agent) WatchLogAndSend(logPath string, interval time.Duration, agentID string) {
	// Lưu offset vào cùng thư mục với logPath, tên file: <logPath>.offset
	offsetPath := logPath + ".offset"
	if !isAbsPath(offsetPath) {
		cwd, _ := os.Getwd()
		offsetPath = cwd + string(os.PathSeparator) + offsetPath
	}
	var lastSize int64 = 0
	// Đọc offset từ file nếu có
	if b, err := os.ReadFile(offsetPath); err == nil {
		fmt.Sscanf(string(b), "%d", &lastSize)
	}
	for {
		file, err := os.Open(logPath)
		if err != nil {
			logutil.Error("WatchLogAndSend: open log file error: %v", err)
			time.Sleep(interval)
			continue
		}
		stat, err := file.Stat()
		if err != nil {
			file.Close()
			logutil.Error("WatchLogAndSend: stat log file error: %v", err)
			time.Sleep(interval)
			continue
		}
		if stat.Size() < lastSize {
			// Nếu file bị truncate, chỉ gửi phần mới
			lastSize = 0
		}
		if stat.Size() > lastSize {
			file.Seek(lastSize, 0)
			buf := make([]byte, stat.Size()-lastSize)
			_, err := file.Read(buf)
			if err == nil {
				lines := splitLines(string(buf))
				for _, line := range lines {
					if line == "" {
						continue
					}
					logMsg := LogData{Message: line}
					msgData := AgentMessageData{AgentID: agentID, Payload: logMsg}
					msg := Message{Type: TypeLog, Data: msgData}
					err := a.Send(msg)
					if err != nil {
						logutil.Error("WatchLogAndSend: send log line error: %v", err)
					}
				}
			}
			lastSize = stat.Size()
			// Lưu lại offset ra file
			os.WriteFile(offsetPath, []byte(fmt.Sprintf("%d", lastSize)), 0644)
		}
		file.Close()
		time.Sleep(interval)
	}
}

func isAbsPath(path string) bool {
	return len(path) > 1 && (path[1] == ':' || path[0] == '/' || path[0] == '\\')
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// getAgentIDFromMsg trích agent_id từ msg.Data nếu có
func getAgentIDFromMsg(msg Message) string {
	if data, ok := msg.Data.(AgentMessageData); ok {
		return data.AgentID
	}
	if m, ok := msg.Data.(map[string]interface{}); ok {
		if id, ok := m["agent_id"].(string); ok {
			return id
		}
	}
	return ""
}

// getPayloadFromMsg trích payload từ msg.Data nếu có
func getPayloadFromMsg(msg Message) interface{} {
	if data, ok := msg.Data.(AgentMessageData); ok {
		return data.Payload
	}
	if m, ok := msg.Data.(map[string]interface{}); ok {
		if p, ok := m["payload"]; ok {
			return p
		}
	}
	return nil
}
