# Hệ thống Client-Server Giám sát Log bằng Go

## 1. Tổng quan

Dự án này là một hệ thống client-server được xây dựng bằng ngôn ngữ Go, sử dụng giao thức TCP cho việc giao tiếp. Chức năng chính của hệ thống là cho phép các "agent" (client) giám sát sự thay đổi trong file log và tự động gửi các dòng log mới đến server một cách an toàn.

- **Server**: Đóng vai trò trung tâm, quản lý việc đăng ký của các agent, nhận và lưu trữ log, đồng thời cung cấp các giao diện quản trị qua API và dòng lệnh (CLI).
- **Client (Agent)**: Là một ứng dụng nhẹ chạy trên máy cần giám sát. Nó tự động đăng ký với server, theo dõi một file log cụ thể (`events.log`), và gửi bất kỳ dòng log mới nào được thêm vào cho server.
- **Log Generator**: Một công cụ phụ trợ để tạo ra các dòng log ngẫu nhiên, phục vụ cho việc kiểm thử.

Toàn bộ dữ liệu truyền giữa client và server đều được mã hóa bằng AES để đảm bảo an toàn.

## 2. Tính năng chính

- **Đăng ký Agent tự động**: Agent tự động đăng ký với server dựa trên thông tin phần cứng (HostID) duy nhất và được cấp một `AgentID` ngắn gọn để dễ dàng quản lý.
- **Giao tiếp mã hóa**: Mọi tin nhắn giữa client và server đều được mã hóa (AES) và đóng gói (message framing) để đảm bảo an toàn và toàn vẹn dữ liệu.
- **Giám sát Log thời gian thực**: Agent sử dụng thư viện `tail` để theo dõi file `events.log` và gửi đi các dòng mới ngay lập tức.
- **Lưu trữ Log tập trung**: Server nhận log từ tất cả các agent và lưu vào file `archiver.log`, kèm theo thông tin `AgentID` và thời gian.
- **Hệ thống Log phân cấp**: Server ghi log hoạt động của chính nó ra file `service.log` và có thể cấu hình mức độ log (INFO, WARNING, ERROR) khi khởi động.
- **Quản trị đa kênh**: Có thể quản lý server (xem danh sách agent, gửi tin nhắn) thông qua:
  - **API HTTP** (Cổng `8081`)
  - **Giao diện dòng lệnh (CLI)** trực tiếp trên terminal của server.

## 3. Cấu trúc thư mục

```
/
├── client/
│   ├── client.go         # Mã nguồn của agent giám sát log
│   ├── client_config.json  # File lưu AgentID và ClientID sau khi đăng ký
│   └── events.log        # File log mà agent sẽ theo dõi
│
├── server/
│   ├── main.go             # Điểm khởi đầu của server
│   ├── tcp_handler.go      # Xử lý kết nối và giao tiếp TCP
│   ├── api_handler.go      # Xử lý các request API HTTP
│   ├── cli_handler.go      # Xử lý các lệnh từ CLI
│   ├── client_manager.go   # Quản lý danh sách client
│   ├── logging.go          # Cấu hình hệ thống log
│   ├── crypto.go           # Hàm mã hóa và giải mã
│   ├── models.go           # Các cấu trúc dữ liệu (structs)
│   ├── registered_clients.json # CSDL lưu thông tin các agent đã đăng ký
│   ├── service.log         # Log hoạt động của server
│   └── archiver.log        # Log do các agent gửi về
│
└── log_generator/
    └── main.go             # Công cụ tạo log ngẫu nhiên để test
```

## 4. Hướng dẫn sử dụng

**Yêu cầu:** Cài đặt Go (phiên bản 1.18 trở lên).

### Bước 1: Chạy Server

Mở một terminal, di chuyển vào thư mục `server` và chạy lệnh:

```bash
# Di chuyển vào thư mục server
cd server

# Chạy server với log level mặc định (INFO)
go run .
```

Server sẽ bắt đầu lắng nghe kết nối TCP ở cổng `8080` và API ở cổng `8081`.

**Tùy chọn cấu hình Log Level:**

Bạn có thể dùng cờ `-loglevel` để thay đổi mức độ chi tiết của log hệ thống:

```bash
# Chỉ ghi log WARNING và ERROR
go run . -loglevel=WARNING

# Chỉ ghi log ERROR
go run . -loglevel=ERROR
```

### Bước 2: Chạy Client (Agent)

Mở một terminal **khác**, di chuyển vào thư mục `client` và chạy lệnh:

```bash
# Di chuyển vào thư mục client
cd client

# Chạy agent
go run .
```

Lần đầu tiên chạy, agent sẽ đăng ký với server và tạo file `client_config.json` để lưu lại định danh của nó. Từ các lần sau, nó sẽ dùng lại định danh này.

### Bước 3: Kiểm thử - Tạo log tự động

Để kiểm tra xem agent có hoạt động hay không, hãy tạo ra một vài dòng log.

Mở một terminal **thứ ba**, di chuyển vào thư mục `log_generator` và chạy lệnh:

```bash
# Di chuyển vào thư mục log_generator
cd log_generator

# Chạy công cụ tạo log
go run .
```

Công cụ này sẽ bắt đầu ghi các dòng log ngẫu nhiên vào file `client/events.log`. Bạn sẽ thấy terminal của **Client (Agent)** thông báo "Phát hiện dòng mới..." và gửi đi, đồng thời terminal của **Server** cũng sẽ báo "Đã nhận tin nhắn từ Agent...".

### Bước 4: Kiểm thử - Quản trị Server

Bạn có thể tương tác với server qua CLI hoặc API.

**Sử dụng CLI:**

Tại terminal đang chạy **Server**, nhập các lệnh sau và nhấn Enter:

- `list`: Xem danh sách tất cả các agent đã đăng ký và trạng thái kết nối.
- `send <AgentID> <Nội dung tin nhắn>`: Gửi một tin nhắn đến một agent cụ thể (ví dụ: `send A001 Hello agent`).
- `help`: Xem các lệnh có sẵn.
- `exit`: Dừng server.

**Sử dụng API (với cURL):**

Mở một terminal **thứ tư** và thực hiện các request:

- **Lấy danh sách agent:**
  ```bash
  curl http://localhost:8081/clients
  ```

- **Gửi tin nhắn cho agent có ID là A001:**
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"message": "Hello from API"}' http://localhost:8081/send/A001
  ```
