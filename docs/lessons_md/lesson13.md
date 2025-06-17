# Interface trong Go

Interface là một tập hợp các method signatures mà một type cần implement. Interface trong Go rất linh hoạt và là nền tảng cho tính trừu tượng và đa hình.

## Khái niệm cơ bản

### Định nghĩa Interface
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Interface có thể kết hợp từ nhiều interface khác
type ReadWriter interface {
    Reader
    Writer
}
```

### Cách Implement Interface
Trong Go, interface được implement một cách ngầm định. Type không cần khai báo rõ ràng việc implement interface:

```go
type File struct {
    // ...
}

// File tự động implement Reader interface
func (f *File) Read(p []byte) (n int, err error) {
    // Implementation
    return len(p), nil
}
```

## Interface thông dụng

### 1. Stringer Interface
```go
type Stringer interface {
    String() string
}

type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}
```

### 2. Error Interface
```go
type error interface {
    Error() string
}

type ValidationError struct {
    Field string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

## Empty Interface

Empty interface (interface{}) không có method nào và được implement bởi tất cả các type:

```go
// Có thể nhận bất kỳ giá trị nào
func PrintAnything(v interface{}) {
    fmt.Printf("Type: %T, Value: %v\n", v, v)
}

// Map với value là empty interface
data := map[string]interface{}{
    "name": "John",
    "age":  30,
    "city": "New York",
}
```

## Type Assertion

Type assertion cho phép kiểm tra và truy cập type cụ thể của một interface:

```go
var i interface{} = "hello"

// Kiểm tra type an toàn
if str, ok := i.(string); ok {
    fmt.Println(str)
} else {
    fmt.Println("i is not a string")
}

// Switch type
switch v := i.(type) {
case string:
    fmt.Printf("String: %s\n", v)
case int:
    fmt.Printf("Integer: %d\n", v)
default:
    fmt.Printf("Unknown type\n")
}
```

## Interface nhỏ là interface tốt

Go khuyến khích tạo các interface nhỏ, tập trung vào một chức năng cụ thể:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Kết hợp khi cần
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

## Ví dụ thực tế về Interface

### 1. Xử lý nhiều loại Storage
```go
type Storage interface {
    Save(data []byte) error
    Load() ([]byte, error)
    Delete() error
}

// Implement cho File Storage
type FileStorage struct {
    path string
}

func (fs FileStorage) Save(data []byte) error {
    return ioutil.WriteFile(fs.path, data, 0644)
}

func (fs FileStorage) Load() ([]byte, error) {
    return ioutil.ReadFile(fs.path)
}

func (fs FileStorage) Delete() error {
    return os.Remove(fs.path)
}

// Implement cho Memory Storage
type MemoryStorage struct {
    data []byte
}

func (ms *MemoryStorage) Save(data []byte) error {
    ms.data = make([]byte, len(data))
    copy(ms.data, data)
    return nil
}

func (ms *MemoryStorage) Load() ([]byte, error) {
    return ms.data, nil
}

func (ms *MemoryStorage) Delete() error {
    ms.data = nil
    return nil
}
```

### 2. Plugin System
```go
type Plugin interface {
    Name() string
    Initialize() error
    Execute() error
    Cleanup() error
}

// Implement cho các plugin cụ thể
type LoggerPlugin struct {
    // ...
}

func (p LoggerPlugin) Name() string { return "Logger" }
func (p LoggerPlugin) Initialize() error { /* ... */ return nil }
func (p LoggerPlugin) Execute() error { /* ... */ return nil }
func (p LoggerPlugin) Cleanup() error { /* ... */ return nil }

// Plugin Manager
type PluginManager struct {
    plugins []Plugin
}

func (pm *PluginManager) Register(p Plugin) {
    pm.plugins = append(pm.plugins, p)
}

func (pm *PluginManager) ExecuteAll() error {
    for _, p := range pm.plugins {
        if err := p.Execute(); err != nil {
            return fmt.Errorf("plugin %s failed: %v", p.Name(), err)
        }
    }
    return nil
}
```

## Best Practices

1. Giữ interface nhỏ và tập trung:
```go
// Tốt
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Không tốt - quá nhiều trách nhiệm
type DataHandler interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
    Close() error
    Flush() error
    // ...nhiều method khác
}
```

2. Đặt interface ở nơi sử dụng:
```go
// Package client định nghĩa interface cho nhu cầu của nó
package client

type DataSource interface {
    Fetch() ([]byte, error)
}

// Hàm chỉ cần method Fetch
func ProcessData(ds DataSource) error {
    data, err := ds.Fetch()
    if err != nil {
        return err
    }
    // Process data
    return nil
}
```

3. Sử dụng interface cho dependency injection:
```go
type Logger interface {
    Log(message string) error
}

type Service struct {
    logger Logger
}

func NewService(logger Logger) *Service {
    return &Service{logger: logger}
}
```

4. Xử lý error với interface:
```go
type CustomError interface {
    error
    Code() int
}

type APIError struct {
    code    int
    message string
}

func (e APIError) Error() string { return e.message }
func (e APIError) Code() int    { return e.code }

func handleError(err error) {
    if apiErr, ok := err.(CustomError); ok {
        fmt.Printf("API Error %d: %s\n", apiErr.Code(), apiErr.Error())
    } else {
        fmt.Printf("General Error: %s\n", err.Error())
    }
}
```
