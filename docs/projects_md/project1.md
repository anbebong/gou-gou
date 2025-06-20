## I. Goroutines
Goroutines là một trong những tính năng đặc biệt nhất của Golang để lập trình Concurrency cực kỳ đơn giản. Goroutines bản chất là các hàm (function) hay method được thực thi một các độc lập và đồng thời nhưng vẫn có thể kết nối với nhau. Một cách ngắn gọn, những thực thi đồng thời được gọi là Goroutines trong Go (Golang).

### Bản chất của Goroutines

Trong bất kỳ một chương trình Golang đều tồn tại ít nhất một Goroutine, gọi là main Goroutine. Nếu main goroutines này kết thúc, toàn bộ các goroutines khác trong chương trình cũng đều bị dừng và kết thúc ngay.

Goroutine bản chất là một lightweight execution thread (luồng thực thi gọn nhẹ). Vì thế việc sử dụng các Goroutines trong Golang có chi phí cực kì thấp so với cách sử dụng các Thread truyền thống (OS Thread).

#### 1.Cách khai báo một Goroutine
Bất kì một hàm nào trong Golang cũng đều có thể chạy đồng thời hay Goroutines với việc thêm vào từ khoá go.

Từ khoá go cũng chính là tên của ngôn ngữ Go, là first-class keyword. Nghĩa là trong Golang bạn không cần phải import bất kì một package nào để sử dụng nó.

```go
func name(){
// statements
}

// using go keyword as the 
// prefix of your function call
go name()
```
Go runtime sẽ:

- Tạo một stack nhỏ (~2KB)
- Đưa function vào scheduler để lên lịch thực thi
- Thực thi function trong một goroutine mới
> Goroutine chỉ kết thúc khi function của nó kết thúc, Hoặc khi main goroutine kết thúc, Không có cách để "kill" một goroutine từ bên ngoài

```go
package main

import (
	"fmt"
	"time"
)

func sayHello(name string) {
	for i := 0; i <= 5; i++ {
		fmt.Printf("%d : Hello %s\n", time.Now().Second(), name)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go sayHello("Viet") // 1s in 1 lần
	sayHello("Nam") // 2s sẽ in ra 1 lần
	time.Sleep(2 * time.Second)
}
```

## II. Channel 

Channel là kênh giao tiếp trung gian giữa các Goroutines có thể gởi và nhận dữ liệu cho nhau một cách an toàn thông qua cơ chế lock-free. Mặc định, Channel là kênh giao tiếp 2 chiều. Nghĩa là Channel có thể dùng cho cả gởi và nhận dữ liệu.

#### Cách khai báo Channel

Để sử dụng Channel, chúng ta dùng keyword chan được hỗ trợ mặc định từ trong chính ngôn ngữ Golang (thường được hiểu là first-class). Chúng ta không cần phải import thêm bất kì một package nào để sử dụng chan.

```go
var channelName chan Type
``` 
hoặc
```go
channelName := make(chan Type)
``` 

#### Gửi nhận dữ liệu qua Channel

- Để gởi và nhận dữ liệu qua Channel, chúng ta sẽ dùng toán tử <-. Toán tử này hoạt động như một chỉ hướng dữ liệu sẽ đi từ đâu đến đâu. Chỉ hướng này giúp ta xác định được dữ liệu đang được gởi đi hay nhận về.

```go
channelName <- value
```
- Nhận dữ liệu từ Channel

```go
myVar := <- channelName
```
- Khi gửi nhận thì bị block
- Ví dụ 
```go
package main

import "fmt"

func main() {
	myChan := make(chan int)

	go func() {
        time.Sleep(time.Second * 2) //
		myChan <- 1 // channel đang nhận 1
	}()

	fmt.Println(<-myChan) // Đợi 2s rồi mới lấy giá trị 1 ra

}
``` 
### Mutex

