# Xử lý lỗi (Error Handling) trong Go

Go xử lý lỗi một cách đơn giản và trực tiếp thông qua interface `error`, không dùng exceptions như các ngôn ngữ khác.

## Error Interface
```go
type error interface {
    Error() string
}
```

## Các cách xử lý lỗi cơ bản

### 1. Return error
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("không thể chia cho 0")
    }
    return a / b, nil
}

// Sử dụng
result, err := divide(10, 0)
if err != nil {
    // Xử lý lỗi
    return err
}
```

### 2. Custom Error
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

### 3. Multiple Errors
```go
type ValidationErrors []ValidationError

// Tập hợp nhiều lỗi vào một slice
func (v ValidationErrors) Error() string {
    if len(v) == 0 {
        return ""
    }
    
    msgs := make([]string, len(v))
    for i, err := range v {
        msgs[i] = err.Error()
    }
    return strings.Join(msgs, "; ")
}
```

## Các Pattern Xử lý Lỗi

### 1. Wrap Error
```go
// Thêm context cho error
if err := readFile(path); err != nil {
    return fmt.Errorf("đọc file %s: %w", path, err)
}
```

### 2. errors.Is và errors.As
```go
// Kiểm tra error cụ thể
if errors.Is(err, os.ErrNotExist) {
    // File không tồn tại
}

// Chuyển đổi error type
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Printf("Lỗi ở: %s\n", valErr.Field)
}
```

### 3. Panic và Recover
```go
func doSomething() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    if critical {
        panic("lỗi nghiêm trọng")
    }
    return nil
}
```

## Best Practices

### 1. Luôn xử lý error
```go
// Tốt ✅
if err != nil {
    return fmt.Errorf("validate: %w", err)
}

// Không tốt ❌
result, _ := someFunc() // Bỏ qua error
```

### 2. Dùng defer để cleanup 
```go
func processFile(path string) (err error) {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()
    return process(f)
}
```

## Ví dụ Thực Tế

### 1. Validation Chain
```go
v := &Validator{}
err := v.Check(user.Name != "", "name", "required").
       Check(user.Age >= 0, "age", "must be positive").
       Errors()
```

### 2. HTTP Error Handler 
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    err := processRequest(r)
    if err == nil {
        return
    }

    switch {
    case errors.Is(err, ErrInvalid):
        http.Error(w, "Invalid input", http.StatusBadRequest)
    case errors.Is(err, ErrNotFound):
        http.Error(w, "Not found", http.StatusNotFound)
    default:
        http.Error(w, "Server error", http.StatusInternalServerError)
    }
}
```
