# Package và Import trong Go

Package là cách Go tổ chức và tái sử dụng code. Mỗi chương trình Go được tạo thành từ các package, và Go cung cấp một hệ thống quản lý package mạnh mẽ.

## Ví dụ Cơ Bản về Package

### 1. Main Package
```go
// main.go
package main

import (
    "fmt"
    "myapp/calculator"  // Import package local
)

func main() {
    result := calculator.Add(5, 3)
    fmt.Printf("5 + 3 = %d\n", result)
    
    fmt.Printf("5 - 3 = %d\n", calculator.Subtract(5, 3))
}
```

### 2. Calculator Package
```go
// calculator/calculator.go
package calculator

// Add trả về tổng của hai số
func Add(a, b int) int {
    return a + b
}

// Subtract trả về hiệu của hai số
func Subtract(a, b int) int {
    return a - b
}

// Hàm này không thể được gọi từ package khác
func multiply(a, b int) int {
    return a * b
}
```

### 3. Multiple Files trong Cùng Package
```go
// calculator/advanced.go
package calculator

// Power tính lũy thừa
func Power(base, exp int) int {
    result := 1
    for i := 0; i < exp; i++ {
        result *= base
    }
    return result
}
```

## Cấu trúc Package

### Package Declaration
```go
// Mỗi file .go phải bắt đầu với package declaration
package main  // Chương trình thực thi
// hoặc
package utils // Package thư viện
```

### Quy tắc đặt tên Package
- Tên package nên ngắn gọn, súc tích
- Sử dụng chữ thường
- Không có dấu gạch ngang hoặc dấu gạch dưới
- Thường là danh từ đơn số

```go
// Tốt ✅
package database
package models
package utils

// Không tốt ❌
package userDatabase
package my_package
package DBUtils
```

## Import Package

### Các Cách Import và Ví Dụ Sử Dụng

```go
// calculator_app.go
package main

// 1. Import đơn
import "fmt"

// 2. Import nhiều package
import (
    "fmt"
    "strings"
    "time"
)

// 3. Import với alias
import (
    "fmt"
    str "strings"        // Sử dụng alias 'str' cho package strings
    . "math"            // Import trực tiếp (không khuyến khích)
    _ "image/png"       // Import cho side-effects
)

func main() {
    // Sử dụng fmt package
    fmt.Println("Hello, World!")
    
    // Sử dụng strings package với alias
    message := "hello, world"
    fmt.Println(str.ToUpper(message))  // HELLO, WORLD
    
    // Sử dụng time package
    now := time.Now()
    fmt.Printf("Hiện tại là: %v\n", now)
    
    // Sử dụng math package (imported với .)
    fmt.Println(Pi)  // Có thể dùng Pi thay vì math.Pi
}
```

### Ví dụ về Import Local Package
```go
// main.go
package main

import (
    "fmt"
    "myapp/models"      // Import package local
    "myapp/utils"
    db "myapp/database" // Sử dụng alias cho package
)

func main() {
    // Sử dụng models package
    user := models.NewUser("john", "john@example.com")
    
    // Sử dụng utils package
    if utils.ValidateEmail(user.Email) {
        // Sử dụng database package với alias
        db.SaveUser(user)
    }
}
```

### Import Path
- Package trong thư viện chuẩn: `"fmt"`, `"os"`, `"net/http"`
- Package bên ngoài: `"github.com/user/repo"`
- Package local: `"myapp/utils"`

## Tổ chức Code

### Cấu trúc thư mục điển hình
```
myproject/
├── main.go
├── go.mod
├── go.sum
├── pkg/
│   ├── models/
│   │   └── user.go
│   ├── database/
│   │   └── db.go
│   └── utils/
│       └── helpers.go
└── internal/
    └── service/
        └── handler.go
```

### File go.mod
```go
module github.com/username/project

go 1.21

require (
    github.com/lib/pq v1.10.9
    github.com/gorilla/mux v1.8.0
)
```

## Exported Names và Private Functions

Trong Go, chỉ những identifier (tên biến, hàm, type) bắt đầu bằng chữ in hoa mới có thể được truy cập từ package khác.

### Ví dụ về Exported vs Unexported Names
```go
// models/user.go
package models

import "time"

// User là exported type, có thể được sử dụng từ package khác
type User struct {
    Name      string    // Exported field
    Email     string    // Exported field
    createdAt time.Time // Unexported field (private)
}

// NewUser là exported function
func NewUser(name, email string) *User {
    return &User{
        Name:      name,
        Email:     email,
        createdAt: time.Now(),
    }
}

// GetCreatedAt là exported method để truy cập unexported field
func (u *User) GetCreatedAt() time.Time {
    return u.createdAt
}

// validateEmail là unexported function (private)
func validateEmail(email string) bool {
    // validation logic
    return true
}

// ValidateUser là exported function sử dụng private function
func ValidateUser(u *User) bool {
    return validateEmail(u.Email)
}
```

### Sử dụng trong Main Package
```go
// main.go
package main

import (
    "fmt"
    "myapp/models"
)

func main() {
    user := models.NewUser("John", "john@example.com")
    
    // Có thể truy cập các exported fields
    fmt.Println(user.Name)   // OK
    fmt.Println(user.Email)  // OK
    
    // KHÔNG thể truy cập unexported fields
    // fmt.Println(user.createdAt)  // Error!
    
    // Nhưng có thể truy cập thông qua exported method
    fmt.Println(user.GetCreatedAt())  // OK
    
    // Có thể gọi exported function
    if models.ValidateUser(user) {
        fmt.Println("User is valid")
    }
    
    // KHÔNG thể gọi unexported function
    // models.validateEmail(user.Email)  // Error!
}
```

## Package Documentation

Go sử dụng comments để tạo documentation. Dưới đây là ví dụ về cách viết documentation đúng chuẩn:

### Ví dụ về Package Documentation
```go
// Package database cung cấp các hàm và interface để tương tác với cơ sở dữ liệu.
// Package này hỗ trợ nhiều loại database khác nhau và cung cấp một interface
// thống nhất cho các thao tác CRUD.
package database

import (
    "context"
    "errors"
)

// ErrNotFound được trả về khi không tìm thấy record trong database
var ErrNotFound = errors.New("record not found")

// Client đại diện cho một kết nối database.
// Nó cung cấp các phương thức để thực hiện các thao tác database cơ bản.
type Client struct {
    connStr string
    timeout int
}

// NewClient tạo một client database mới với connection string được chỉ định.
// Timeout là thời gian tối đa (tính bằng giây) cho mỗi operation.
// Trả về error nếu không thể kết nối đến database.
func NewClient(connStr string, timeout int) (*Client, error) {
    // implementation
}

// Query thực hiện một SQL query và trả về kết quả.
// Các tham số:
//   - ctx: context để control timeout và cancellation
//   - query: câu SQL query
//   - args: các tham số cho prepared statement
//
// Trả về:
//   - Kết quả của query
//   - Error nếu có lỗi xảy ra hoặc context bị cancel
func (c *Client) Query(ctx context.Context, query string, args ...interface{}) (Result, error) {
    // implementation
}
```

### Sử dụng godoc
Sau khi viết documentation, bạn có thể xem nó bằng lệnh `go doc`:

```bash
# Xem documentation của package
go doc database

# Xem chi tiết một type
go doc database.Client

# Xem documentation của một method
go doc database.Client.Query
```

## Ví dụ Thực tế

### 1. Package Models
```go
// pkg/models/user.go
package models

type User struct {
    ID       int
    Username string
    Email    string
}

func NewUser(username, email string) *User {
    return &User{
        Username: username,
        Email:    email,
    }
}
```

### 2. Package Database
```go
// pkg/database/db.go
package database

import "myapp/pkg/models"

type UserRepository interface {
    GetUser(id int) (*models.User, error)
    SaveUser(user *models.User) error
}

type DBHandler struct {
    // configuration
}

func (db *DBHandler) GetUser(id int) (*models.User, error) {
    // implementation
}
```

### 3. Package Utils
```go
// pkg/utils/validator.go
package utils

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func ValidateEmail(email string) bool {
    return emailRegex.MatchString(email)
}
```

## Best Practices

1. Package Organization:
```go
// Tổ chức theo chức năng
package models    // Định nghĩa cấu trúc dữ liệu
package handlers // Xử lý HTTP requests
package database // Tương tác với database
package utils    // Các hàm tiện ích
```

2. Circular Dependencies:
- Tránh import cycle giữa các package
- Sử dụng interface để giảm sự phụ thuộc

3. Package Names:
```go
// Tốt
package models
package utils
package handlers

// Tránh
package model  // dùng số nhiều
package utilityFunctions  // quá dài
```

4. Internal Packages:
```
myapp/
├── internal/  // Chỉ có thể import bởi myapp
│   └── auth/
├── pkg/      // Có thể import bởi các app khác
│   └── utils/
```

5. Tài liệu hóa:
```go
// Package config cung cấp các hàm để đọc
// và xử lý cấu hình ứng dụng từ các nguồn khác nhau.
package config

// LoadConfig đọc file cấu hình và trả về Config struct.
// Nếu file không tồn tại, trả về error.
func LoadConfig(filename string) (*Config, error) {
    // implementation
}
```

## Tools hữu ích

1. `go doc`: Xem documentation
```bash
go doc fmt.Println
go doc -all fmt
```

2. `go mod`: Quản lý dependencies
```bash
go mod init myproject
go mod tidy
go mod vendor
```

3. `go list`: Liệt kê packages
```bash
go list ./...
go list -m all
```
