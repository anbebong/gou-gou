# Đọc và Ghi File trong Go

Go cung cấp nhiều cách để làm việc với file thông qua package os và bufio. Trong bài này, chúng ta sẽ học các thao tác cơ bản nhất với file.

## 1. Đọc File

### Đọc toàn bộ file một lần
```go
data, err := os.ReadFile("test.txt")
if err != nil {
    fmt.Println("Lỗi đọc file:", err)
    return
}
fmt.Println(string(data))
```
- Hàm os.ReadFile() đọc toàn bộ nội dung file vào bộ nhớ
- Trả về dữ liệu dạng []byte và error (nếu có)
- Phù hợp với file nhỏ

### Đọc file theo từng dòng
```go
file, err := os.Open("test.txt")
if err != nil {
    fmt.Println("Lỗi mở file:", err)
    return
}
defer file.Close()  // Đảm bảo file được đóng

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}
```
- Dùng bufio.Scanner để đọc từng dòng
- defer file.Close() đảm bảo file luôn được đóng
- Phù hợp với file lớn hoặc cần xử lý theo dòng

## 2. Ghi File

### Ghi file một lần
```go
content := []byte("Xin chào từ Go!")
err := os.WriteFile("test.txt", content, 0644)
if err != nil {
    fmt.Println("Lỗi ghi file:", err)
    return
}
```
- os.WriteFile() ghi toàn bộ dữ liệu vào file
- 0644 là quyền truy cập file (read/write cho owner, read cho others)
- File sẽ được tạo mới nếu chưa tồn tại

### Ghi file theo từng phần
```go
file, err := os.Create("test.txt")
if err != nil {
    fmt.Println("Lỗi tạo file:", err)
    return
}
defer file.Close()

writer := bufio.NewWriter(file)
writer.WriteString("Dòng 1\n")
writer.WriteString("Dòng 2\n")
writer.Flush()  // Đẩy dữ liệu từ buffer xuống file
```
- Dùng bufio.Writer để ghi từng phần
- Flush() đảm bảo dữ liệu được ghi xuống file
- Phù hợp khi cần ghi nhiều lần

## 3. Xử lý Lỗi

Luôn kiểm tra lỗi khi làm việc với file:
- Lỗi mở file (file không tồn tại, không có quyền)
- Lỗi đọc/ghi (ổ đĩa đầy, file bị khóa)
- Đảm bảo đóng file sau khi dùng xong
