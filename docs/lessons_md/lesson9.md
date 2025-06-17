# Mảng (Array) và Slice trong Go

## Mảng (Array)

Mảng trong Go là một tập hợp các phần tử cùng kiểu dữ liệu với kích thước cố định.

### Khai báo mảng

```go
var numbers [5]int              // Mảng 5 phần tử kiểu int, mặc định các giá trị là 0
colors := [3]string{"red", "green", "blue"}  // Khai báo và khởi tạo giá trị
matrix := [2][3]int{{1, 2, 3}, {4, 5, 6}}   // Mảng đa chiều
```

### Đặc điểm của mảng
- Kích thước cố định khi khai báo
- Các phần tử có cùng kiểu dữ liệu
- Khi truyền mảng vào hàm, Go sẽ copy toàn bộ mảng (tốn bộ nhớ)

## Slice

Slice là một cấu trúc dữ liệu linh hoạt hơn mảng, cho phép thay đổi kích thước động.

### Khai báo slice

```go
var s []int                     // Slice rỗng, giá trị là nil
numbers := []int{1, 2, 3, 4, 5} // Khai báo và khởi tạo trực tiếp
s := make([]int, 5)             // Tạo slice với len=5, cap=5
s := make([]int, 3, 5)          // Tạo slice với len=3, cap=5
```

### Cắt slice (slicing)

```go
array := [5]int{1, 2, 3, 4, 5}
slice := array[1:4]     // Lấy phần tử từ index 1 đến 3 (4-1)
slice := array[:3]      // Từ đầu đến index 2
slice := array[2:]      // Từ index 2 đến cuối
slice := array[:]       // Toàn bộ mảng
```

### Các thao tác với slice

1. Thêm phần tử:
```go
slice = append(slice, 6)        // Thêm một phần tử
slice = append(slice, 7, 8, 9)  // Thêm nhiều phần tử
```

2. Copy slice:
```go
slice2 := make([]int, len(slice))
copy(slice2, slice)
```

### Capacity và Length

- `len(slice)`: Số phần tử hiện có trong slice
- `cap(slice)`: Dung lượng tối đa của slice trước khi cần mở rộng bộ nhớ

```go
s := make([]int, 3, 5)
fmt.Println(len(s))  // 3
fmt.Println(cap(s))  // 5
```

### Lưu ý quan trọng

1. Slice là tham chiếu đến một mảng
2. Khi bạn thay đổi phần tử trong slice, mảng gốc cũng bị thay đổi
3. Khi append vượt quá capacity, Go sẽ tạo một mảng mới với capacity lớn hơn
4. Slice có thể là nil, khác với mảng rỗng

### Khi nào dùng Array vs Slice?

- Dùng Array khi:
  - Cần kích thước cố định
  - Muốn đảm bảo không có thay đổi về số lượng phần tử
  - Làm việc với dữ liệu nhỏ và cố định

- Dùng Slice khi:
  - Cần một tập hợp động có thể thay đổi kích thước
  - Truyền tham chiếu qua các hàm (tối ưu bộ nhớ)
  - Làm việc với dữ liệu có kích thước không xác định trước

## Best Practices

1. Ưu tiên sử dụng slice thay vì array trong hầu hết các trường hợp
2. Dùng `make()` để khởi tạo slice với capacity phù hợp nếu biết trước kích thước
3. Sử dụng `append()` để thêm phần tử thay vì truy cập trực tiếp
4. Cẩn thận với slice of slice để tránh giữ tham chiếu không cần thiết
