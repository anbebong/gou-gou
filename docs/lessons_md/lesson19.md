# Tạo module đơn giản trong Go

### Các bước thực hiện:

1. Tạo và cài đặt thư mục module
```bash
mkdir mylog
cd mylog
go mod init mylog
```

2. Tạo file chứa code của module:
```go
package mylog

import "fmt"

func Output(a ...interface{}) {
    fmt.Println(a...)
}
```

3. Sử dụng module trong chương trình chính:
```go
package main

import "mylog"

func main() {
    mylog.Output("Hello from custom module!")
}
```

Giờ đây thay vì sử dụng `fmt.Println()`, bạn có thể dùng `mylog.Output()` từ module của riêng mình.
