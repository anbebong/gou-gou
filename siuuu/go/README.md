# Hệ thống Client-Server Giám sát Log bằng Go

## 1. Tổng quan

- **Client** là một Windows service, giao tiếp với server qua TCP (mã hóa AES), đồng thời cung cấp IPC (named pipe) để các ứng dụng hệ thống (như Credential Provider) lấy OTP hoặc thông tin bảo mật.
- **Server** quản lý agent, nhận log, cung cấp API HTTP (JWT) và CLI để quản trị, phân quyền động qua file user.
- **Log Generator** là công cụ tạo log test cho agent.

## 2. Cấu trúc thư mục
```
/
├── client/
│   ├── main.go                # Service client, giao tiếp TCP & IPC
│   ├── client_config.json     # Lưu ClientID, AgentID, serverAddress, logFile
│   ├── client.log             # Log hoạt động của client
│   ├── events.log             # File log agent theo dõi
│   ├── ipc_test_client.go     # Test giao tiếp IPC (named pipe)
│   └── ... (các file phụ trợ)
│
├── server/
│   ├── main.go                # Khởi động server, lắng nghe TCP & API
│   ├── tcp_handler.go         # Xử lý giao tiếp TCP với client
│   ├── api_handler.go         # Xử lý API HTTP (JWT, phân quyền)
│   ├── cli_handler.go         # CLI quản trị server
│   ├── client_manager.go      # Quản lý client đã đăng ký
│   ├── models.go              # Định nghĩa struct dữ liệu
│   ├── clients.json           # Lưu thông tin client đã đăng ký
│   ├── users.json             # Lưu thông tin user, phân quyền
│   ├── service.log            # Log hoạt động server
│   ├── archiver.log           # Log tập trung từ agent
│   └── ... (các file phụ trợ)
│
├── log_generator/
│   └── main.go                # Công cụ tạo log test cho agent
│
└── web/
    └── index.html             # Giao diện web (nếu có)
```

## 3. Hướng dẫn sử dụng

### Server
```bash
cd server
go run .
```
- Server lắng nghe TCP (mặc định 8080), API HTTP (8081).
- Có thể cấu hình cổng và loglevel qua tham số dòng lệnh.

### Client (Agent Service)
Client là một Windows service, giao tiếp với server qua TCP và với hệ thống qua IPC (named pipe).

- **Cài đặt service:**
  ```bash
  cd client
  go run . install
  ```
- **Khởi động service:**
  ```bash
  go run . start
  ```
- **Dừng service:**
  ```bash
  go run . stop
  ```
- **Gỡ cài đặt service:**
  ```bash
  go run . remove
  ```
- **Chạy debug (foreground, log ra console):**
  ```bash
  go run . debug
  ```

Sau khi cài đặt và khởi động, agent sẽ tự động đăng ký với server, tạo file `client_config.json` lưu định danh và cấu hình.

- **Giao tiếp IPC:**
  - Service lắng nghe named pipe `\\.\pipe\MySecretServicePipe`.
  - Ứng dụng khác (ví dụ Credential Provider) có thể gửi chuỗi `GET_SECRET` vào pipe này để nhận OTP mới từ server.
  - Có thể test IPC bằng file `ipc_test_client.go` trong thư mục client.

### Log Generator
```bash
cd log_generator
go run .
```
- Tự động ghi log vào `client/events.log`, agent sẽ gửi log mới về server.

## 4. Hệ thống user, phân quyền và xác thực JWT (API)
- User lưu trong `server/users.json` với trường: `username`, `password`, `role`.
- 2 quyền: `admin` (toàn quyền), `user` (chỉ xem và đổi mật khẩu của mình).
- Đăng nhập qua `/api/login` nhận JWT, gửi JWT qua header `Authorization: Bearer <token>` cho các API cần xác thực.
- Dữ liệu user luôn được load lại từ file trước mỗi thao tác quan trọng.

### Phân quyền:
- **Admin**: Toàn quyền (tạo user, đổi mật khẩu bất kỳ, đổi role, xóa client, gửi tin nhắn...)
- **User thường**: Chỉ xem thông tin, đổi mật khẩu của chính mình.

## 5. Tài liệu API và hướng dẫn test

### Đăng nhập lấy JWT
```
POST http://localhost:8081/api/login
Content-Type: application/json
{
  "username": "admin",
  "password": "adminpass"
}
```
**Response:**
```
{
  "token": "<JWT_TOKEN>"
}
```

### Lấy danh sách client (cần JWT)
```
curl -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8081/clients
```

### Đổi mật khẩu (user chỉ đổi được của mình, admin đổi được của bất kỳ ai)
```
POST http://localhost:8081/api/users/change-password
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
{
  "username": "user1",
  "oldPassword": "user1pass",
  "newPassword": "newpass123"
}
```

### Tạo user mới (chỉ admin)
```
POST http://localhost:8081/api/users/create
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
{
  "username": "newuser",
  "password": "newpass",
  "role": "user"
}
```

### Đổi role user (chỉ admin)
```
POST http://localhost:8081/api/users/update
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
{
  "username": "user1",
  "role": "admin"
}
```

### Gán user cho client (chỉ admin)
```
POST http://localhost:8081/api/clients/assign-user
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
{
  "id": "<clientID hoặc agentID>",
  "username": "user1"
}
```

### Xóa client (chỉ admin)
```
curl -X DELETE -H "Authorization: Bearer <JWT_TOKEN>" "http://localhost:8081/clients/delete?id=<agentID>"
```

### Gửi tin nhắn tới client (chỉ admin)
```
POST http://localhost:8081/message/send
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
{
  "targetID": "<agentID>",
  "data": "Nội dung tin nhắn"
}
```

### Lấy OTP cho client (cần JWT)
```
curl -H "Authorization: Bearer <JWT_TOKEN>" "http://localhost:8081/api/otp?id=<agentID>"
```

## 6. Hướng dẫn test IPC
- Có thể dùng file `ipc_test_client.go` để test gửi yêu cầu `GET_SECRET` tới pipe `\\.\pipe\MySecretServicePipe` và nhận về OTP từ server.
- Đảm bảo service client đang chạy trước khi test IPC.

## 7. Lưu ý
- Luôn gửi JWT token trong header cho các API cần xác thực.
- Nếu đổi mật khẩu hoặc role, cần đăng nhập lại để lấy token mới.
- Dữ liệu user được cập nhật động, không cần restart server.
- Client chỉ nên thao tác qua lệnh service, không chạy trực tiếp bằng `go run .` nếu không phải debug.

---
