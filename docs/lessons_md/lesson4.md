# Câu lệnh điều kiện (if, switch)

Go cung cấp câu lệnh điều kiện if và switch để kiểm soát luồng thực thi của chương trình.

### 1. Câu lệnh if-else

- **Cấu trúc cơ bản:**
  ```go
  if điều_kiện {
      khối_lệnh
  }
  ```
- Có thể kết hợp với *else* hoặc *else if*
- Không cần dấu ngoặc đơn `()` quanh điều kiện, nhưng phải có ngoặc nhọn `{}`

**Ví dụ:**
```go
if x > 10 {
    // Xử lý khi x > 10
} else if x > 5 {
    // Xử lý khi 5 < x ≤ 10
} else {
    // Xử lý khi x ≤ 5
}
```

### 2. If với khai báo biến
- Go cho phép khai báo biến ngay trong câu lệnh if
- Biến này chỉ tồn tại trong phạm vi của khối if-else

**Ví dụ:**
```go
if y := 2*x; y > 10 {
    // y tồn tại ở đây
} else {
    // và ở đây
}
// Không thể sử dụng y ở đây
```

### 3. Câu lệnh switch

- Dùng để so sánh một biến với nhiều giá trị khác nhau
- **Cú pháp:**
  ```go
  switch biểu_thức {
  case giá_trị_1:
      // Mã thực thi khi biểu_thức = giá_trị_1
  case giá_trị_2, giá_trị_3:
      // Mã thực thi khi biểu_thức = giá_trị_2 hoặc 3
  default:
      // Mã thực thi khi không có case nào khớp
  }
  ```

### 4. Switch không có biểu thức

- Go cho phép switch không cần biểu thức
- Tự động so sánh với "true"
- Tương đương với chuỗi if-else if

**Ví dụ:**
```go
switch {
case x > 10:
    // Thực thi khi x > 10
case x > 5:
    // Thực thi khi 5 < x ≤ 10
default:
    // Thực thi khi x ≤ 5
}
```

### 5. Fallthrough

- Thông thường, khi một case được khớp và thực thi, switch sẽ kết thúc
- Từ khóa "fallthrough" cho phép tiếp tục thực thi case tiếp theo

**Ví dụ:**
```go
switch {
case x < 10:
    // Mã thực thi khi x < 10
    fallthrough
case x < 20:
    // Mã này LUÔN thực thi nếu case trước đúng và có fallthrough
}
```

**Lưu ý:**  
Không như các ngôn ngữ khác, Go tự động "break" sau mỗi case, không cần thêm "break".
