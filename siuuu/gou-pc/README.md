# Gou-PC Service

## Mô tả
Service Golang gồm:
- TCP server sử dụng mã hoá TLS (port 9000)
- RESTful API với Gin (port 8080)

## Chạy thử
1. Cài Go >= 1.21
2. Cài dependency: `go mod tidy`
3. Tạo file `server.crt` và `server.key` (self-signed hoặc CA)
4. Chạy: `go run ./cmd/main.go`

## Thư mục
- `cmd/`: entrypoint
- `internal/rest/`: REST API
- `internal/tcpserver/`: TCP server (TLS)

## API mẫu
- `GET /ping` trả về `{ "message": "pong" }`

## TCP mẫu
- Kết nối TLS tới port 9000, gửi chuỗi bất kỳ, nhận về "Hello from server (TLS)!"
