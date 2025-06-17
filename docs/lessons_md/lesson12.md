# Method trong Go

Method trong Go là một hàm đặc biệt gắn với một kiểu dữ liệu cụ thể (receiver type). Method cho phép chúng ta thêm hành vi cho các kiểu dữ liệu, tương tự như phương thức trong lập trình hướng đối tượng.

Methods giúp code trở nên rõ ràng và có tổ chức hơn bằng cách:
- Nhóm các hàm liên quan với kiểu dữ liệu
- Cải thiện khả năng đọc và bảo trì code
- Hỗ trợ tính đóng gói và trừu tượng hóa

## Định nghĩa Method

Cấu trúc cơ bản của một method:
```go
func (receiver receiverType) MethodName(parameters) returnType {
    // Method body
}
```

Ví dụ với một struct đơn giản:

```go
type Rectangle struct {
    width  float64
    height float64
}

// Method Area có receiver type là Rectangle
func (r Rectangle) Area() float64 {
    return r.width * r.height
}
```

## Value Receivers vs Pointer Receivers

### Value Receivers
Value receiver tạo một bản sao của receiver khi method được gọi:
```go
// Method không thay đổi receiver
func (r Rectangle) Area() float64 {
    return r.width * r.height
}

// Sử dụng
rect := Rectangle{width: 10, height: 5}
area := rect.Area()  // rect không bị thay đổi
```

### Pointer Receivers
Pointer receiver cho phép method thay đổi giá trị của receiver:
```go
// Method có thể thay đổi receiver
func (r *Rectangle) Scale(factor float64) {
    r.width *= factor
    r.height *= factor
}

// Sử dụng
rect := Rectangle{width: 10, height: 5}
rect.Scale(2)  // rect.width và rect.height bị thay đổi
```

## Khi nào dùng Value/Pointer Receiver?

Sự khác biệt chính giữa value và pointer receiver:

### Dùng Value Receiver khi:
- Method không cần thay đổi receiver
```go
func (s String) Length() int {
    return len(string(s))
}
```
- Receiver là kiểu dữ liệu cơ bản hoặc nhỏ
```go
type MyInt int
func (m MyInt) IsPositive() bool {
    return m > 0
}
```
- Receiver là map, func, chan (vì bản thân chúng đã là reference type)
```go
type Cache map[string]string
func (c Cache) Has(key string) bool {
    _, exists := c[key]
    return exists
}
```

### Dùng Pointer Receiver khi:
- Method cần thay đổi receiver
```go
func (u *User) SetPassword(password string) {
    u.password = hash(password)
}
```
- Receiver là struct lớn (để tránh copy)
```go
type LargeStruct struct {
    Data [1024]int
}
func (l *LargeStruct) Process() {
    // Xử lý trên con trỏ tránh copy dữ liệu lớn
}
```
- Muốn tất cả method của type đều nhất quán
```go
type Counter struct {
    count int
}
func (c *Counter) Increment() { c.count++ }
func (c *Counter) Decrement() { c.count-- }
func (c *Counter) Value() int { return c.count }
```

## Method Set

Method set xác định những method nào có thể được gọi trên một type:

Method set là tập hợp các method có thể được gọi trên một type:

```go
type Person struct {
    Name string
    Age  int
}

// Method set của Person
func (p Person) GetInfo() string {
    return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

// Method set của *Person
func (p *Person) SetAge(age int) {
    p.Age = age
}
```

## Method với Non-struct Types

Go cho phép định nghĩa method cho bất kỳ type nào được định nghĩa trong cùng package (trừ type từ package khác). Điều này rất hữu ích khi bạn muốn mở rộng chức năng cho các kiểu dữ liệu đơn giản:

```go
// Định nghĩa type mới dựa trên type có sẵn
type Celcius float64

// Thêm methods cho type
func (c Celcius) ToFahrenheit() float64 {
    return float64(c)*9/5 + 32
}

func (c Celcius) String() string {
    return fmt.Sprintf("%.2f°C", c)
}

// Sử dụng
temp := Celcius(25.5)
fmt.Println(temp)                // "25.50°C"
fmt.Println(temp.ToFahrenheit()) // 77.9
```

## Method Chaining

Method chaining là một kỹ thuật cho phép gọi nhiều method liên tiếp, làm cho code trở nên dễ đọc và súc tích hơn:

Method chaining cho phép gọi nhiều method liên tiếp:

```go
type Builder struct {
    str string
}

func (b *Builder) Append(s string) *Builder {
    b.str += s
    return b
}

// Sử dụng: 
// builder.Append("Hello").Append(" ").Append("World")
```

## Functional Options Pattern

Functional Options là một pattern phổ biến trong Go để cấu hình đối tượng một cách linh hoạt:

```go
type Server struct {
    host     string
    port     int
    timeout  time.Duration
    maxConns int
}

// Option là một function type để cấu hình Server
type Option func(*Server)

// Các function tạo Option
func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func WithMaxConns(n int) Option {
    return func(s *Server) {
        s.maxConns = n
    }
}

// Constructor sử dụng variadic options
func NewServer(opts ...Option) *Server {
    // Giá trị mặc định
    s := &Server{
        host:     "localhost",
        port:     8080,
        timeout:  30 * time.Second,
        maxConns: 100,
    }
    
    // Áp dụng các option
    for _, opt := range opts {
        opt(s)
    }
    
    return s
}

// Sử dụng:
server := NewServer(
    WithHost("example.com"),
    WithPort(443),
    WithTimeout(1 * time.Minute),
    WithMaxConns(1000),
)
```

Pattern này có nhiều ưu điểm:
- Cho phép cấu hình linh hoạt với nhiều tùy chọn
- Dễ dàng thêm tùy chọn mới mà không phải thay đổi code hiện có
- Giá trị mặc định hợp lý
- API rõ ràng và dễ sử dụng

## Method với Interface

Method là cách để type implement interface:

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle implements Shape
type Rectangle struct {
    width, height float64
}

func (r Rectangle) Area() float64 {
    return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.width + r.height)
}
```

## Best Practices

1. Đặt tên method rõ ràng và nhất quán:
```go
// Tốt ✅
type User struct {
    name  string
    admin bool
}

func (u User) IsAdmin() bool            // Rõ ràng về chức năng
func (u User) HasPermission(p Permission) bool  // Nhất quán với quy ước đặt tên
func (u *User) SetName(name string)     // Prefix Set cho method thay đổi giá trị
func (u User) String() string           // Implement interface chuẩn

// Không tốt ❌
func (u User) Admin() bool             // Không rõ là kiểm tra hay đặt giá trị
func (u User) CheckPerm(p Permission)  // Không nhất quán với tên khác
func (u *User) UpdateName(name string) // Nên dùng SetName để nhất quán
```

2. Group các method liên quan và comment rõ ràng:
```go
// User struct và các method của nó
type User struct {
    email    string
    password string
}

// Validation methods
func (u User) ValidateEmail() error {
    if !strings.Contains(u.email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}

func (u User) ValidatePassword() error {
    if len(u.password) < 8 {
        return errors.New("password too short")
    }
    return nil
}

// Validate thực hiện tất cả các validation
func (u User) Validate() error {
    if err := u.ValidateEmail(); err != nil {
        return err
    }
    return u.ValidatePassword()
}

// Database operations
func (u *User) Save() error {
    if err := u.Validate(); err != nil {
        return err
    }
    // Save to database
    return nil
}

func (u *User) Delete() error {
    // Delete from database
    return nil
}
```

3. Tránh receiver name dài và sử dụng tên có ý nghĩa:
```go
// Tốt ✅
func (u User) Save() error       // Ngắn gọn
func (db Database) Query() error // Mô tả rõ vai trò
func (s *Server) Start() error   // Tên có ý nghĩa

// Không tốt ❌
func (thisUser User) Save() error      // Quá dài
func (x Database) Query() error        // Không mô tả
func (serverInstance *Server) Start()  // Thừa thãi
```

4. Nhất quán trong việc sử dụng receiver type:
```go
// Tốt ✅: Nhất quán sử dụng pointer receiver cho tất cả method
type Counter struct {
    value int
}

func (c *Counter) Increment() { c.value++ }
func (c *Counter) Decrement() { c.value-- }
func (c *Counter) Value() int { return c.value }

// Không tốt ❌: Trộn lẫn value và pointer receiver
type Counter struct {
    value int
}

func (c Counter) Value() int { return c.value }
func (c *Counter) Increment() { c.value++ }
func (c Counter) IsZero() bool { return c.value == 0 }
```

5. Xử lý lỗi trong method:
```go
type Account struct {
    balance float64
}

func (a *Account) Withdraw(amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    if amount > a.balance {
        return fmt.Errorf("insufficient funds: have %.2f, need %.2f", 
            a.balance, amount)
    }
    a.balance -= amount
    return nil
}
```

## Ví dụ thực tế

1. Logger với Method Chaining:
```go
type Logger struct {
    prefix    string
    level     string
    output    io.Writer
    timestamp bool
}

func (l *Logger) WithPrefix(prefix string) *Logger {
    l.prefix = prefix
    return l
}

func (l *Logger) WithLevel(level string) *Logger {
    l.level = level
    return l
}

func (l *Logger) WithTimestamp() *Logger {
    l.timestamp = true
    return l
}

func (l *Logger) Log(message string) {
    timestamp := ""
    if l.timestamp {
        timestamp = time.Now().Format("2006-01-02 15:04:05 ")
    }
    fmt.Fprintf(l.output, "%s[%s] %s: %s\n",
        timestamp, l.level, l.prefix, message)
}

// Sử dụng:
logger := &Logger{output: os.Stdout}
logger.WithPrefix("APP").
       WithLevel("ERROR").
       WithTimestamp().
       Log("Something went wrong!")
```

2. Cache với TTL (Time To Live):
```go
type CacheItem struct {
    Value      interface{}
    Expiration time.Time
}

type Cache struct {
    items map[string]CacheItem
    mu    sync.RWMutex
}

func NewCache() *Cache {
    c := &Cache{
        items: make(map[string]CacheItem),
    }
    go c.startCleanup()
    return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.items[key] = CacheItem{
        Value:      value,
        Expiration: time.Now().Add(ttl),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.items[key]
    if !exists {
        return nil, false
    }
    
    if time.Now().After(item.Expiration) {
        return nil, false
    }
    
    return item.Value, true
}

func (c *Cache) startCleanup() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        c.cleanup()
    }
}

func (c *Cache) cleanup() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    now := time.Now()
    for key, item := range c.items {
        if now.After(item.Expiration) {
            delete(c.items, key)
        }
    }
}
```

3. Config Manager với Validation:
```go
type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
}

type ServerConfig struct {
    Host    string
    Port    int
    Timeout time.Duration
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
}

func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return errors.New("database host is required")
    }
    if c.Database.Port <= 0 {
        return errors.New("invalid database port")
    }
    if c.Server.Port <= 0 {
        return errors.New("invalid server port")
    }
    if c.Server.Timeout <= 0 {
        return errors.New("invalid server timeout")
    }
    return nil
}

func (c *Config) WithDefaultTimeout() *Config {
    if c.Server.Timeout == 0 {
        c.Server.Timeout = 30 * time.Second
    }
    return c
}

// Sử dụng:
config := &Config{
    Database: DatabaseConfig{
        Host: "localhost",
        Port: 5432,
    },
    Server: ServerConfig{
        Host: "0.0.0.0",
        Port: 8080,
    },
}

err := config.WithDefaultTimeout().Validate()
if err != nil {
    log.Fatal(err)
}
```
