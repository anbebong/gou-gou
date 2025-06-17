# Goroutines trong Go

Goroutine là một đơn vị thực thi nhẹ được quản lý bởi Go runtime. Goroutines cho phép chúng ta thực hiện concurrent programming (lập trình đồng thời) một cách dễ dàng và hiệu quả.

## Goroutine là gì?

- Goroutine là một hàm có thể chạy đồng thời với các hàm khác
- Nhẹ hơn thread của hệ điều hành (chỉ vài KB bộ nhớ)
- Được quản lý bởi Go runtime, không phải hệ điều hành
- Có thể chạy hàng nghìn goroutines cùng lúc

## Cách tạo Goroutine

Để tạo một goroutine, chỉ cần thêm từ khóa `go` trước lời gọi hàm:

```go
// Hàm thông thường
sayHello()

// Chạy như goroutine
go sayHello()
```

## Ví dụ cơ bản

```go
func main() {
    // Chạy hàm như goroutine
    go printNumbers()
    go printLetters()
    
    // Đợi để thấy kết quả
    time.Sleep(time.Second)
}

func printNumbers() {
    for i := 1; i <= 5; i++ {
        fmt.Printf("%d ", i)
        time.Sleep(100 * time.Millisecond)
    }
}

func printLetters() {
    for i := 'a'; i <= 'e'; i++ {
        fmt.Printf("%c ", i)
        time.Sleep(100 * time.Millisecond)
    }
}
```

## Synchronization (Đồng bộ hóa)

### 1. WaitGroup
WaitGroup dùng để đợi một nhóm goroutines hoàn thành:

```go
var wg sync.WaitGroup

func main() {
    wg.Add(2)  // Đợi 2 goroutines
    
    go func() {
        defer wg.Done()
        // Do something
    }()
    
    go func() {
        defer wg.Done()
        // Do something else
    }()
    
    wg.Wait()  // Đợi tất cả hoàn thành
}
```

### 2. Mutex
Mutex dùng để bảo vệ dữ liệu khỏi data races:

```go
var (
    counter int
    mutex   sync.Mutex
)

func increment() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}
```

## Ứng dụng thực tế

### 1. Xử lý nhiều requests cùng lúc
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    go processRequest(r)  // Xử lý trong background
    w.Write([]byte("Request received"))
}

func processRequest(r *http.Request) {
    // Xử lý request mất nhiều thời gian
    time.Sleep(2 * time.Second)
    log.Printf("Processed request: %s", r.URL.Path)
}
```

### 2. Concurrent Data Processing
```go
func processData(data []int) {
    numCPU := runtime.NumCPU()
    chunk := len(data) / numCPU
    var wg sync.WaitGroup
    
    for i := 0; i < numCPU; i++ {
        wg.Add(1)
        go func(start int) {
            defer wg.Done()
            // Xử lý một phần của data
            for j := start; j < start+chunk; j++ {
                process(data[j])
            }
        }(i * chunk)
    }
    
    wg.Wait()
}
```

## Best Practices

1. Xử lý lỗi trong goroutines:
```go
func runTask() error {
    errChan := make(chan error, 1)
    
    go func() {
        if err := doSomething(); err != nil {
            errChan <- err
            return
        }
        errChan <- nil
    }()
    
    return <-errChan
}
```

2. Giới hạn số lượng goroutines:
```go
func processItems(items []Item) {
    semaphore := make(chan struct{}, 5) // Tối đa 5 goroutines
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        semaphore <- struct{}{} // Acquire
        
        go func(item Item) {
            defer func() {
                <-semaphore // Release
                wg.Done()
            }()
            process(item)
        }(item)
    }
    
    wg.Wait()
}
```

3. Graceful Shutdown:
```go
func main() {
    done := make(chan bool)
    go worker(done)
    
    // Đợi signal để shutdown
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    <-signals
    
    // Thông báo worker dừng
    close(done)
    // Đợi worker cleanup và thoát
    time.Sleep(time.Second)
}

func worker(done chan bool) {
    for {
        select {
        case <-done:
            fmt.Println("Cleaning up...")
            return
        default:
            // Do work
        }
    }
}
```

## Lưu ý quan trọng

1. Không nên dùng goroutine cho những tác vụ nhỏ
2. Cẩn thận với việc chia sẻ dữ liệu giữa các goroutines
3. Luôn đảm bảo goroutines được kết thúc đúng cách
4. Sử dụng context để quản lý lifecycle của goroutines
5. Xử lý panic trong goroutines để tránh crash chương trình

## Debugging Goroutines

1. Race Detector:
```bash
go run -race main.go
```

2. Trace goroutines:
```go
import "runtime/trace"

func main() {
    trace.Start(os.Stderr)
    defer trace.Stop()
    // Your code here
}
```

3. Print số lượng goroutines:
```go
fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
```
