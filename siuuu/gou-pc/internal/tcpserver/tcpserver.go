package tcpserver

import (
	"encoding/binary"
	"encoding/json"
	"gou-pc/internal/agent"
	"gou-pc/internal/config"
	"gou-pc/internal/crypto"
	"gou-pc/internal/logcollector"
	"gou-pc/internal/logutil"
	"io"
	"net"
	"os"
	"time"
)

// Định nghĩa struct cho log archive
type ArchiveLogEntry = logcollector.ArchiveLogEntry

func Start(cfg *config.ServerConfig) error {
	ln, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		logutil.Error("failed to listen: %v", err)
		return err
	}
	defer ln.Close()
	logutil.Info("TCP server (AES) listening on %s", cfg.ListenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			logutil.Error("accept error: %v", err)
			continue
		}
		go handleConn(conn, cfg)
	}
}

func handleConn(conn net.Conn, cfg *config.ServerConfig) {
	defer conn.Close()
	for {
		lenBuf := make([]byte, 4)
		_, err := io.ReadFull(conn, lenBuf)
		if err != nil {
			if err != io.EOF {
				logutil.Error("read length error: %v", err)
			}
			return
		}
		msgLen := binary.BigEndian.Uint32(lenBuf)
		if msgLen == 0 || msgLen > 65536 {
			logutil.Error("invalid message length")
			return
		}
		msgBuf := make([]byte, msgLen)
		_, err = io.ReadFull(conn, msgBuf)
		if err != nil {
			logutil.Error("read message error: %v", err)
			return
		}
		decrypted, err := crypto.Decrypt(string(msgBuf))
		if err != nil {
			logutil.Error("decrypt error: %v", err)
			return
		}

		var req agent.Message
		if err := json.Unmarshal([]byte(decrypted), &req); err != nil {
			logutil.Error("invalid message format: %v", err)
			return
		}
		logutil.Info("Received: {type:%s, agent_id:%s, data:%v}", req.Type, getAgentIDFromMsg(req), req.Data)

		var resp agent.Message
		switch req.Type {
		case agent.TypeRegister:
			// Xử lý đăng ký: kiểm tra manager_client.json
			clients, _ := agent.LoadClients(cfg.ClientDBFile) // truyền file path
			devInfo := agent.DeviceInfo{}
			b, _ := json.Marshal(req.Data)
			_ = json.Unmarshal(b, &devInfo)
			found := agent.FindClientByDevice(devInfo, clients)
			if found == nil {
				clientID := agent.GenClientID()
				agentID := agent.GenAgentID()
				newClient := agent.ManagedClient{ClientID: clientID, AgentID: agentID, DeviceInfo: devInfo}
				clients = append(clients, newClient)
				_ = agent.SaveClients(cfg.ClientDBFile, clients) // truyền file path
				resp = agent.Message{
					Type: agent.TypeRegister,
					Data: map[string]string{"client_id": clientID, "agent_id": agentID},
				}
			} else {
				resp = agent.Message{
					Type: agent.TypeRegister,
					Data: map[string]string{"client_id": found.ClientID, "agent_id": found.AgentID},
				}
			}
		case agent.TypeRequestOTP:
			// Chuẩn hoá: luôn lấy agent_id từ req.Data nếu có
			var agentID string
			if m, ok := req.Data.(map[string]interface{}); ok {
				if v, ok := m["agent_id"].(string); ok {
					agentID = v
				}
			}
			logutil.Info("[REQUEST OTP] from agent_id=%s, data=%v", agentID, req.Data)
			resp = agent.Message{
				Type: agent.TypeRequestOTP,
				Data: map[string]interface{}{"agent_id": agentID, "otp": "123456"},
			}
		case agent.TypeHello:
			var agentID string
			if m, ok := req.Data.(map[string]interface{}); ok {
				if v, ok := m["agent_id"].(string); ok {
					agentID = v
				}
			}
			logutil.Info("[HELLO] from agent_id=%s, data=%v", agentID, req.Data)
			resp = agent.Message{
				Type: agent.TypeHello,
				Data: map[string]interface{}{"agent_id": agentID, "payload": req.Data},
			}
		case agent.TypeLog:
			var agentID, message string
			if m, ok := req.Data.(map[string]interface{}); ok {
				if v, ok := m["agent_id"].(string); ok {
					agentID = v
				}
				if payload, ok := m["payload"].(map[string]interface{}); ok {
					if msg, ok := payload["message"].(string); ok {
						message = msg
					}
				}
			}
			logEntry := ArchiveLogEntry{
				Time:    time.Now().Format(time.RFC3339),
				AgentID: agentID,
				Message: message,
			}
			if b, err := json.Marshal(logEntry); err == nil {
				f, err := os.OpenFile(cfg.ArchiveFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				if err == nil {
					f.Write(b)
					f.Write([]byte("\n"))
					f.Close()
				}
			}
			logutil.Info("[CLIENT LOG] %v", logEntry)
			resp = agent.Message{
				Type: agent.TypeLog,
				Data: map[string]interface{}{"agent_id": agentID, "result": "log received"},
			}
		default:
			resp = agent.Message{
				Type: agent.TypeError,
				Data: "Unknown request type",
			}
		}

		respJson, _ := json.Marshal(resp)
		encrypted, err := crypto.Encrypt(string(respJson))
		if err != nil {
			logutil.Error("encrypt error: %v", err)
			return
		}
		respBytes := []byte(encrypted)
		respLen := uint32(len(respBytes))
		lenResp := make([]byte, 4)
		binary.BigEndian.PutUint32(lenResp, respLen)
		if _, err := conn.Write(lenResp); err != nil {
			logutil.Error("write response length error: %v", err)
			return
		}
		if _, err := conn.Write(respBytes); err != nil {
			logutil.Error("write response error: %v", err)
			return
		}
		logutil.Info("Sent: {type:%s, agent_id:%v, data:%v}", resp.Type, getAgentIDFromResp(resp), resp.Data)
	}
}

func getAgentIDFromResp(resp agent.Message) interface{} {
	if m, ok := resp.Data.(map[string]interface{}); ok {
		return m["agent_id"]
	}
	return nil
}

func getAgentIDFromMsg(msg agent.Message) interface{} {
	if m, ok := msg.Data.(map[string]interface{}); ok {
		return m["agent_id"]
	}
	return nil
}

// LoadArchiveLogs đọc toàn bộ log từ file archive.log
var LoadArchiveLogs = logcollector.LoadArchiveLogs
