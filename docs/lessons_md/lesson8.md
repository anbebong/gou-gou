# Con trỏ (Pointer) trong Go

Trong Go, con trỏ là một biến đặc biệt dùng để lưu trữ địa chỉ bộ nhớ của một biến khác. Khi bạn có con trỏ trỏ đến một biến, bạn có thể thay đổi giá trị của biến đó thông qua con trỏ.

## Cách sử dụng con trỏ

1. Để khai báo một con trỏ, ta sử dụng dấu * trước kiểu dữ liệu
2. Để lấy địa chỉ của một biến, ta sử dụng toán tử &
3. Để lấy giá trị mà con trỏ đang trỏ tới (dereference), ta sử dụng dấu *

### Ví dụ đơn giản

```go
x := 10          // Biến thường
p := &x          // p là con trỏ trỏ đến x
fmt.Println(*p)  // In ra giá trị mà p trỏ tới (10)
*p = 20          // Thay đổi giá trị của x thông qua p
fmt.Println(x)   // In ra 20
```

### Con trỏ với hàm

Khi truyền tham số vào hàm dưới dạng con trỏ, hàm có thể thay đổi giá trị của biến gốc:

```go
func changeValue(p *int) {
    *p = 100
}

x := 1
changeValue(&x)  // Truyền địa chỉ của x vào hàm
fmt.Println(x)   // In ra 100
```

## Lưu ý quan trọng

1. Giá trị mặc định của con trỏ là nil
2. Không thể thực hiện phép toán số học với con trỏ
3. Cẩn thận khi sử dụng con trỏ để tránh "nil pointer dereference"

### Sử dụng con trỏ với struct

Con trỏ thường được sử dụng với struct trong Go:

```go
type Person struct {
    Name string
    Age  int
}

p := &Person{Name: "Alice", Age: 25}
fmt.Println(p.Name)  // Có thể truy cập trực tiếp, không cần (*p).Name
```

## Khi nào nên sử dụng con trỏ?

1. Khi muốn một hàm có thể thay đổi giá trị của tham số truyền vào
2. Khi làm việc với struct lớn để tránh copy dữ liệu
3. Khi implement interface
4. Khi muốn thể hiện rằng một giá trị có thể là nil
