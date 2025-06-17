# Biến, hằng số và kiểu dữ liệu cơ bản

Go là **ngôn ngữ kiểu tĩnh** (statically typed), yêu cầu khai báo kiểu dữ liệu cho biến.

### 1. Khai báo biến
Go có nhiều cách để khai báo biến:

#### Khai báo đầy đủ
```go
var tên_biến kiểu_dữ_liệu = giá_trị
var message string = "Hello Go"
```

#### Khai báo ngắn gọn (tự suy luận kiểu)
```go
tên_biến := giá_trị
message := "Hello Go" // Tự động là string
```

#### Khai báo nhiều biến
```go
var a, b int = 1, 2
x, y := 10, "hello"
```

### 2. Các kiểu dữ liệu cơ bản
| Kiểu dữ liệu | Mô tả | Ví dụ |
|--------------|-------|-------|
| bool | Giá trị logic | true/false |
| string | Chuỗi ký tự | "Hello" |
| int, int8, int16, int32, int64 | Số nguyên có dấu | -10, 255 |
| uint, uint8, uint16, uint32, uint64 | Số nguyên không dấu | 0, 255 |
| float32, float64 | Số thực | 3.14159 |
| complex64, complex128 | Số phức | 1+2i |
| byte | Alias cho uint8 | 97 ('a') |
| rune | Alias cho int32 | 'あ' |

### 3. Hằng số
Hằng số được khai báo với từ khóa `const` và giá trị không thể thay đổi trong quá trình chạy.
```go
const Pi = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)
```

### 4. Zero Values
Khi khai báo biến mà không gán giá trị, Go tự động gán **zero value** theo kiểu dữ liệu:
| Kiểu dữ liệu | Zero Value |
|--------------|------------|
| Số | 0 |
| bool | false |
| string | "" (chuỗi rỗng) |
| pointer | nil |

> **Lưu ý**: Khác với các ngôn ngữ khác, Go không cho phép biến chưa được sử dụng trong code.

### Ví dụ tổng hợp
```go
package main
import "fmt"
func main() {
    // Khai báo và sử dụng biến
    var age int = 25
    name := "Gopher"
    // Zero values
    var score int    // = 0
    var isActive bool // = false
    // Hằng số
    const MaxUsers = 1000
    fmt.Println(name, "is", age, "years old")
}
```
