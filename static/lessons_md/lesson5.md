# Arrays và Slices trong Go

Arrays (mảng) và Slices là hai cấu trúc dữ liệu quan trọng trong Go để lưu trữ và quản lý tập hợp các phần tử cùng kiểu.

## 1. Arrays (Mảng)

Arrays trong Go có kích thước cố định và được khai báo với độ dài xác định. Một khi đã được khai báo, kích thước của array không thể thay đổi.

Cú pháp khai báo array:
```go
var arr [5]int                     // Mảng 5 số nguyên, mặc định các phần tử = 0
arr2 := [3]string{"Go", "is", "fun"} // Khai báo và khởi tạo giá trị
arr3 := [...]int{1, 2, 3, 4}      // Compiler tự suy ra độ dài = 4
```

Một số đặc điểm của array:
- Độ dài là một phần của kiểu dữ liệu: `[5]int` và `[6]int` là hai kiểu khác nhau
- Arrays được truyền vào hàm theo giá trị (pass by value), không phải tham chiếu
- Có thể so sánh hai array cùng kiểu bằng toán tử `==`

## 2. Slices

Slices là cấu trúc dữ liệu linh hoạt hơn arrays, cho phép thay đổi kích thước động. Slice là một tham chiếu đến một phần của array.

Cú pháp khai báo slice:
```go
var s []int                     // Slice rỗng, nil
s1 := []int{1, 2, 3}           // Slice với giá trị khởi tạo
s2 := make([]int, 5)           // Slice độ dài 5, capacity 5
s3 := make([]int, 3, 5)        // Slice độ dài 3, capacity 5
```

Các thao tác với slice:

1. Cắt slice (slicing):
```go
arr := [5]int{1, 2, 3, 4, 5}
s1 := arr[1:4]    // [2 3 4]     - từ index 1 đến 3
s2 := arr[:3]     // [1 2 3]     - từ đầu đến index 2
s3 := arr[2:]     // [3 4 5]     - từ index 2 đến hết
```

2. Thêm phần tử với append:
```go
s := []int{1, 2, 3}
s = append(s, 4)        // [1 2 3 4]
s = append(s, 5, 6, 7)  // [1 2 3 4 5 6 7]
```

3. Nối hai slice:
```go
s1 := []int{1, 2}
s2 := []int{3, 4}
s1 = append(s1, s2...)  // [1 2 3 4]
```

## 3. Các khái niệm quan trọng về Slice

### Length và Capacity

- Length (len): số phần tử hiện có trong slice
- Capacity (cap): số phần tử tối đa slice có thể chứa trước khi cần mở rộng
```go
s := make([]int, 3, 5)
fmt.Println(len(s))  // 3
fmt.Println(cap(s))  // 5
```

### Zero Value và nil

- Zero value của slice là nil
- Một nil slice không có array nền tảng
- Một empty slice có length và capacity = 0 nhưng có array nền tảng
```go
var s []int         // nil slice
fmt.Println(s == nil)  // true

s = []int{}        // empty slice
fmt.Println(s == nil)  // false
```

## 4. Best Practices

1. **Sử dụng Slice thay vì Array**: Trong hầu hết các trường hợp, nên sử dụng slice vì tính linh hoạt của nó.

2. **Chỉ định capacity khi biết trước**: Khi biết được số lượng phần tử cần thêm vào, nên chỉ định capacity để tránh việc phải cấp phát lại bộ nhớ nhiều lần.
```go
s := make([]int, 0, 1000)  // slice rỗng với capacity 1000
```

3. **Cẩn thận với slice của slice**: Khi tạo slice từ một slice khác, chúng chia sẻ cùng array nền tảng. Thay đổi một slice có thể ảnh hưởng đến slice khác.

4. **Copy khi cần**: Sử dụng `copy()` để tạo một bản sao độc lập của slice:
```go
s1 := []int{1, 2, 3}
s2 := make([]int, len(s1))
copy(s2, s1)
```

## 5. Common Pitfalls

1. **Truy cập ngoài phạm vi**:
```go
s := []int{1, 2, 3}
fmt.Println(s[5])  // panic: runtime error: index out of range
```

2. **Quên kiểm tra nil**:
```go
var s []int
if s != nil {  // luôn kiểm tra nil trước khi thao tác
    // do something
}
```

3. **Giữ tham chiếu không cần thiết**:
```go
s := []int{1, 2, 3, 4, 5}
s = append(s[:2], s[3:]...)  // Xóa phần tử thứ 3
```

Arrays và Slices là nền tảng cho việc xử lý dữ liệu trong Go. Hiểu rõ về chúng sẽ giúp bạn viết code hiệu quả và tránh được các lỗi phổ biến.
