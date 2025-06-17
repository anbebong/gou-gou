# Map trong Go

Map là một cấu trúc dữ liệu lưu trữ các cặp key-value, trong đó mỗi key là duy nhất. Map trong Go tương tự như dictionary trong Python hoặc object trong JavaScript.

## Khai báo Map

Có nhiều cách để khai báo một map:

```go
// Cách 1: Khai báo map rỗng
var scores map[string]int

// Cách 2: Sử dụng make
ages := make(map[string]int)

// Cách 3: Khai báo và khởi tạo cùng lúc
colors := map[string]string{
    "red":   "#ff0000",
    "green": "#00ff00",
    "blue":  "#0000ff",
}
```

## Thao tác với Map

### 1. Thêm hoặc cập nhật phần tử

```go
ages := make(map[string]int)

// Thêm phần tử mới
ages["Alice"] = 25
ages["Bob"] = 30

// Cập nhật giá trị
ages["Alice"] = 26
```

### 2. Truy cập phần tử

```go
age := ages["Alice"]    // Lấy giá trị

// Kiểm tra key có tồn tại không
age, exists := ages["Charlie"]
if exists {
    fmt.Println("Tuổi của Charlie:", age)
} else {
    fmt.Println("Không tìm thấy Charlie")
}
```

### 3. Xóa phần tử

```go
delete(ages, "Bob")    // Xóa phần tử với key "Bob"
```

### 4. Duyệt qua map

```go
for key, value := range colors {
    fmt.Printf("Key: %s, Value: %s\n", key, value)
}

// Chỉ duyệt qua keys
for key := range colors {
    fmt.Println("Key:", key)
}
```

## Lưu ý quan trọng

1. Map là kiểu tham chiếu
   - Khi truyền map vào hàm, các thay đổi sẽ ảnh hưởng đến map gốc
   - Giá trị zero của map là nil

2. Key của map phải là kiểu dữ liệu có thể so sánh được
   - Số, string, bool, interface, pointer, channel, array, struct (nếu các trường có thể so sánh)
   - KHÔNG thể dùng slice, map hoặc function làm key

3. Value của map có thể là bất kỳ kiểu dữ liệu nào
   - Bao gồm cả map khác hoặc struct

## Các trường hợp sử dụng phổ biến

1. Lưu trữ dữ liệu dạng từ điển
```go
dictionary := map[string]string{
    "hello": "xin chào",
    "world": "thế giới",
}
```

2. Đếm tần suất
```go
wordCount := make(map[string]int)
for _, word := range words {
    wordCount[word]++
}
```

3. Cache dữ liệu
```go
cache := make(map[string]interface{})
cache["user"] = getUser()  // Lưu kết quả để dùng lại
```

## Best Practices

1. Khởi tạo map với kích thước phù hợp
```go
// Nếu biết trước số lượng phần tử
users := make(map[string]User, 100)
```

2. Kiểm tra sự tồn tại của key trước khi sử dụng
```go
if value, ok := map["key"]; ok {
    // Sử dụng value
}
```

3. Xử lý map nil
```go
var m map[string]int    // map nil
if m == nil {
    m = make(map[string]int)
}
```

4. Clear map
```go
// Cách 1: Tạo map mới
m = make(map[string]int)

// Cách 2: Xóa từng phần tử
for k := range m {
    delete(m, k)
}
```

5. Map là không an toàn cho concurrent access
   - Sử dụng sync.Map cho trường hợp concurrent
   - Hoặc dùng mutex để bảo vệ truy cập
