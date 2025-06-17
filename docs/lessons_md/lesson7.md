# Hàm (Function) trong Go

Hàm là khối code thực hiện một nhiệm vụ cụ thể, có thể nhận tham số đầu vào và trả về kết quả.

### 1. Khai báo hàm cơ bản
- **Cú pháp:**
  ```go
  func tên_hàm(tham_số kiểu_dữ_liệu) kiểu_trả_về {
      // Thân hàm
      return giá_trị
  }
  ```

**Ví dụ:**
```go
func sayHello() {
    fmt.Println("Hello, Go!")
}

func add(a, b int) int {
    return a + b
}
```

### 2. Tham số hàm
- Có thể có nhiều tham số cùng kiểu
- Tham số được truyền theo giá trị (pass by value)

**Ví dụ:**
```go
// Nhiều tham số cùng kiểu
func multiply(a, b int) int {
    return a * b
}

// Nhiều tham số khác kiểu
func greet(name string, age int) {
    fmt.Printf("Hello, %s! You are %d years old.\n", name, age)
}
```

### 3. Nhiều giá trị trả về
- Go cho phép hàm trả về nhiều giá trị
- Thường dùng để trả về kết quả và lỗi

**Ví dụ:**
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("không thể chia cho 0")
    }
    return a / b, nil
}

// Sử dụng
result, err := divide(10, 2)
if err != nil {
    fmt.Println("Lỗi:", err)
} else {
    fmt.Println("Kết quả:", result)
}
```

### 4. Named Return Values
- Có thể đặt tên cho giá trị trả về
- Tự động return các giá trị đã đặt tên

**Ví dụ:**
```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // Tự động trả về x và y
}
```

### 5. Variadic Functions
- Hàm nhận số lượng tham số không xác định
- Sử dụng dấu ... trước kiểu dữ liệu

**Ví dụ:**
```go
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}

// Sử dụng
fmt.Println(sum(1, 2))          // 3
fmt.Println(sum(1, 2, 3, 4, 5)) // 15
```

### 6. Hàm là kiểu dữ liệu bậc nhất
- Có thể gán hàm cho biến
- Truyền hàm như tham số
- Trả về hàm từ hàm khác

**Ví dụ:**
```go
// Hàm như tham số
func compute(fn func(float64, float64) float64) float64 {
    return fn(3, 4)
}

// Sử dụng
hypot := func(x, y float64) float64 {
    return math.Sqrt(x*x + y*y)
}
fmt.Println(compute(hypot))    // 5
```

### 7. Closure
- Hàm có thể truy cập biến bên ngoài phạm vi của nó

**Ví dụ:**
```go
func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

// Sử dụng
pos := adder()
fmt.Println(pos(1)) // 1
fmt.Println(pos(2)) // 3
fmt.Println(pos(3)) // 6
```

### 8. Defer
- Hoãn thực thi một hàm đến khi hàm chứa nó return
- Thường dùng để cleanup resources (đóng file, kết nối,...)
- Thực thi theo thứ tự LIFO (Last In First Out)

**Ví dụ:**
```go
func main() {
    defer fmt.Println("world")
    fmt.Println("hello")
}
// Output:
// hello
// world
```
