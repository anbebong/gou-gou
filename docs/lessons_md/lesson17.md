# Channels trong Go

Channel là cách để các goroutines giao tiếp và đồng bộ hóa với nhau. Channel hoạt động như một đường ống, cho phép truyền dữ liệu giữa các goroutines.

## 1. Khái niệm cơ bản

### Tạo channel
```go
ch := make(chan int)    // Channel không có buffer
ch := make(chan int, 3) // Channel có buffer size = 3
```

### Đóng channel
```go
close(ch) // Đóng channel, không thể gửi thêm dữ liệu
```

### Gửi và nhận dữ liệu
```go
ch <- value   // Gửi dữ liệu vào channel
value := <-ch // Nhận dữ liệu từ channel
```

## 2. Buffered vs Unbuffered Channels

### Unbuffered Channel (không buffer)
- Chỉ có thể gửi khi có người nhận sẵn sàng
- Gửi và nhận xảy ra đồng thời
- Tự động đồng bộ hóa

```go
ch := make(chan int)
go func() {
    ch <- 1 // Block cho đến khi có người nhận
}()
fmt.Println(<-ch) // Nhận dữ liệu
```

### Buffered Channel (có buffer)
- Có thể gửi khi buffer chưa đầy
- Block khi buffer đầy
- Buffer size > 0

```go
ch := make(chan int, 2)
ch <- 1  // OK, buffer còn trống
ch <- 2  // OK, buffer còn trống
// ch <- 3  // Block, buffer đã đầy
```

## 3. Cách sử dụng phổ biến

### Pipeline Pattern

```go
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}

func processor(in <-chan int, out chan<- int) {
    for num := range in {
        out <- num * 2
    }
    close(out)
}
```

### Fan-out Pattern

```go
// Phân chia công việc cho nhiều worker
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2
    }
}
```

## 4. Channel Direction

### Send-only channel
```go
func send(ch chan<- int) {
    ch <- 42
}
```

### Receive-only channel
```go
func receive(ch <-chan int) {
    val := <-ch
}
```

## 5. Xử lý timeout

```go
select {
case data := <-ch:
    // Xử lý dữ liệu
case <-time.After(2 * time.Second):
    // Timeout sau 2 giây
}
```

## 6. Best Practices

1. **Đóng channel đúng cách**
- Chỉ sender được đóng channel
- Không đóng channel nhiều lần
- Không gửi vào channel đã đóng

2. **Xử lý channel đã đóng**
```go
val, ok := <-ch
if !ok {
    // Channel đã đóng
}
```

3. **Range over channel**
```go
for value := range ch {
    // Tự động dừng khi channel đóng
}
```

## 7. Các lỗi thường gặp

1. Deadlock
```go
ch := make(chan int)
ch <- 1  // Deadlock! Không có ai nhận
```

2. Panic khi đóng channel đã đóng
```go
close(ch)
close(ch) // Panic!
```

3. Gửi vào channel đã đóng
```go
close(ch)
ch <- 1 // Panic!
```
