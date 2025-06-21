package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// startCLI khởi động giao diện dòng lệnh của server.
func startCLI() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Server CLI đã khởi động. Gõ 'help' để xem các lệnh.")
	for {
		fmt.Print("-> ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			ErrorLogger.Printf("Lỗi khi đọc từ CLI: %v", err)
			continue
		}
		cmdString = strings.TrimSpace(cmdString)
		parts := strings.Split(cmdString, " ")
		command := parts[0]

		switch command {
		case "help":
			fmt.Println("Các lệnh có sẵn:")
			fmt.Println("  list              - Liệt kê tất cả các client đã đăng ký.")
			fmt.Println("  delete <agentID>  - Xóa một client đã đăng ký bằng AgentID.")
			fmt.Println("  send <agentID> <message> - Gửi tin nhắn đến một client đang kết nối.")
			fmt.Println("  exit              - Tắt server.")
		case "list":
			registeredClientsMutex.Lock()

			displayClients := make(map[string]interface{})
			for _, client := range registeredClients {
				displayClients[client.AgentID] = client.HardwareInfo
			}

			data, _ := json.MarshalIndent(displayClients, "", "  ")
			fmt.Println(string(data))
			registeredClientsMutex.Unlock()
		case "delete":
			if len(parts) < 2 {
				fmt.Println("Sử dụng: delete <agentID>")
				continue
			}

			idToDelete := parts[1]
			fullClientID, clientInfo, found := findClientByAnyID(idToDelete)

			if found {
				registeredClientsMutex.Lock()
				delete(registeredClients, fullClientID)
				saveRegisteredClients()
				registeredClientsMutex.Unlock()
				fmt.Printf("Đã xóa Agent %s (ClientID: %s)\n", clientInfo.AgentID, fullClientID)
			} else {
				fmt.Printf("Client với AgentID '%s' không tìm thấy\n", idToDelete)
			}
		case "send":
			if len(parts) < 3 {
				fmt.Println("Sử dụng: send <agentID> <message>")
				continue
			}
			targetID := parts[1]
			message := strings.Join(parts[2:], " ")

			fullClientID, clientInfo, found := findClientByAnyID(targetID)
			if !found {
				fmt.Printf("Không tìm thấy client với AgentID: %s\n", targetID)
				continue
			}

			if err := sendMessageToClient(fullClientID, message); err != nil {
				fmt.Printf("Lỗi khi gửi tin nhắn: %v\n", err)
			} else {
				fmt.Printf("Tin nhắn đã được gửi đến Agent %s\n", clientInfo.AgentID)
			}
		case "exit":
			InfoLogger.Println("Đang tắt server...")
			os.Exit(0)
		case "":
			// không làm gì
		default:
			fmt.Println("Lệnh không xác định. Gõ 'help'.")
		}
	}
}
