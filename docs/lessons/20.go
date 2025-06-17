package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Bài tập 1: Đọc file log và tìm ERROR
func findErrors(filename string) {
	// Tạo file log mẫu
	logContent := `2023-06-17 10:00:00 INFO: Server started
2023-06-17 10:01:15 ERROR: Connection failed
2023-06-17 10:01:30 INFO: Retry connection
2023-06-17 10:01:45 ERROR: Database error
2023-06-17 10:02:00 INFO: Connection restored`

	os.WriteFile(filename, []byte(logContent), 0644)

	// Đọc và lọc các dòng ERROR
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Lỗi mở file:", err)
		return
	}
	defer file.Close()

	fmt.Println("Các dòng log có ERROR:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ERROR") {
			fmt.Println(line)
		}
	}
}

// Bài tập 2: Config struct và các hàm thao tác
type Config struct {
	Host string
	Port int
	User string
}

func SaveConfig(cfg Config, filename string) error {
	content := fmt.Sprintf("host=%s\nport=%d\nuser=%s", cfg.Host, cfg.Port, cfg.User)
	return os.WriteFile(filename, []byte(content), 0644)
}

func LoadConfig(filename string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	// Parse config đơn giản
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "host":
			cfg.Host = parts[1]
		case "port":
			fmt.Sscanf(parts[1], "%d", &cfg.Port)
		case "user":
			cfg.User = parts[1]
		}
	}
	return cfg, nil
}

func main() {
	// Demo bài tập 1: Tìm ERROR trong log
	fmt.Println("=== Bài tập 1: Tìm ERROR trong log ===")
	findErrors("log.txt")

	// Demo bài tập 2: Lưu và đọc config
	fmt.Println("\n=== Bài tập 2: Config file ===")
	cfg := Config{
		Host: "localhost",
		Port: 8080,
		User: "admin",
	}

	SaveConfig(cfg, "config.txt")

	loadedCfg, _ := LoadConfig("config.txt")
	fmt.Printf("Loaded config: %+v\n", loadedCfg)
}
