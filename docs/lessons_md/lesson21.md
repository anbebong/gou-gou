## 1. Ví dụ với net/http (chuẩn Go)

```go
package main

import (
    "fmt"
    "net/http"
)

// Hàm handler nhận request và ghi ra response
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Xin chào từ net/http!")
}

func main() {
    // Đăng ký handler cho đường dẫn "/hello"
    http.HandleFunc("/hello", helloHandler)
    fmt.Println("Server chạy tại http://localhost:8080/hello")
    // Khởi động HTTP server tại cổng 8080
    http.ListenAndServe(":8080", nil)
}
```

**Giải thích:**
- `http.HandleFunc("/hello", helloHandler)`: Khi có request đến `/hello`, hàm `helloHandler` sẽ được gọi.
- `helloHandler`: Nhận hai tham số `w` (ghi dữ liệu trả về cho client), `r` (chứa thông tin request).
- `http.ListenAndServe(":8080", nil)`: Khởi động server tại cổng 8080.

**Cách chạy:**
1. Lưu file, chạy: `go run main.go`
2. Truy cập [http://localhost:8080/hello](http://localhost:8080/hello) để xem kết quả.

---

## 2. Ví dụ với Gin framework

**Gin** là một framework mạnh mẽ và phổ biến để xây dựng ứng dụng web và RESTful API bằng ngôn ngữ Go (Golang). Gin được thiết kế để:

- Đơn giản hóa việc xây dựng web server/API so với net/http thuần.
- Giúp code ngắn gọn, dễ đọc, dễ mở rộng.
- Hỗ trợ nhiều tính năng hiện đại như: routing, middleware, binding dữ liệu, validation, trả về JSON/XML, nhóm route, logging, error handling...

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    // Tạo router với middleware logger và recovery mặc định
    r := gin.Default()
    // Đăng ký route GET /hello
    r.GET("/hello", func(c *gin.Context) {
        // Trả về chuỗi text cho client
        c.String(200, "Xin chào từ Gin!")
    })
    // Khởi động server tại cổng 8080
    r.Run(":8080")
}
```

**Giải thích:**
- `gin.Default()`: Tạo router có sẵn middleware log và tự động bắt lỗi (recovery).
- `r.GET("/hello", ...)`: Đăng ký route GET cho đường dẫn `/hello`.
- `c *gin.Context`: Đối tượng chứa thông tin request, response, params, v.v.
- `c.String(200, "...")`: Trả về chuỗi text với mã trạng thái 200.

**Cách chạy:**
1. Cài Gin: `go get github.com/gin-gonic/gin`
2. Lưu file, chạy: `go run main.go`
3. Truy cập [http://localhost:8080/hello](http://localhost:8080/hello) để xem kết quả.

---

## 3. Giải thích về `gin.Default()`

- `gin.Default()` là hàm khởi tạo Gin router với **middleware mặc định**:
    - **Logger**: Tự động log mọi request.
    - **Recovery**: Tự động bắt panic và trả về lỗi 500, tránh server bị crash.
- Bạn có thể thay thế bằng `gin.New()` nếu muốn tuỳ biến middleware.

**Ví dụ:**
```go
r := gin.New()
r.Use(gin.Logger())
r.Use(gin.Recovery())
```
Cách này tương đương với `gin.Default()` nhưng bạn tự cấu hình từng middleware.

---

## 4. So sánh nhanh

| net/http              | Gin                        |
|-----------------------|----------------------------|
| Có sẵn trong Go       | Cần cài thêm Gin           |
| Cú pháp thuần Go      | Cú pháp gọn, dễ mở rộng    |
| Tính năng cơ bản      | Hỗ trợ middleware, nhóm route, trả về JSON, ... |
| Ít tính năng nâng cao | Nhiều tiện ích hiện đại    |

---

## 5. Khi nào chọn net/http, khi nào chọn Gin?

- **net/http**: Phù hợp dự án nhỏ, học tập, hoặc khi bạn muốn kiểm soát mọi thứ chi tiết.
- **Gin**: Phù hợp dự án vừa và lớn, cần phát triển nhanh, muốn có sẵn nhiều tính năng tiện lợi như middleware, binding dữ liệu, trả về JSON, group routes...

---