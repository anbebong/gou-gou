# Vòng lặp (for, break, continue)

Go chỉ hỗ trợ một loại vòng lặp duy nhất là **for**, nhưng với nhiều hình thức linh hoạt.

### 1. For tiêu chuẩn (với 3 thành phần)
- **Cú pháp:**
  ```go
  for khởi_tạo; điều_kiện; hậu_xử_lý {
      khối_lệnh
  }
  ```
- Các thành phần đều có thể bỏ trống

**Ví dụ:**
```go
for i := 0; i < 10; i++ {
    // Mã thực thi lặp lại
}
```

### 2. For như while
- Chỉ giữ lại thành phần điều kiện
- **Cú pháp:**
  ```go
  for điều_kiện {
      khối_lệnh
  }
  ```

**Ví dụ:**
```go
i := 0
for i < 10 {
    // Mã thực thi
    i++
}
```

### 3. For vô hạn
- **Cú pháp:**
  ```go
  for {
      khối_lệnh
  }
  ```
- Vòng lặp sẽ chạy mãi mãi trừ khi gặp:
  - `break`: thoát khỏi vòng lặp
  - `return`: thoát khỏi hàm
  - `panic`: gây ra lỗi

**Ví dụ:**
```go
for {
    // Lặp vô hạn cho đến khi break
    if điều_kiện {
        break
    }
}
```

### 4. For với range
- Dùng để duyệt qua các phần tử trong collection
- **Cú pháp:**
  ```go
  for key, value := range collection {
      khối_lệnh
  }
  ```
- Hoạt động với: array, slice, string, map, channel

**Ví dụ:**
- Cho mảng/slice:
  ```go
  for i, value := range arr {
      // i là chỉ số, value là giá trị
  }
  ```
- Cho map:
  ```go
  for key, value := range myMap {
      // key là khóa, value là giá trị
  }
  ```
- Cho string:
  ```go
  for i, char := range str {
      // i là vị trí, char là ký tự (kiểu rune)
  }
  ```

### 5. break và continue
- `break`: thoát khỏi vòng lặp hiện tại
- `continue`: bỏ qua phần còn lại của lần lặp và tiếp tục lần lặp kế tiếp

**Ví dụ:**
```go
for i := 0; i < 10; i++ {
    if i == 5 {
        continue  // Bỏ qua khi i = 5
    }
    if i == 8 {
        break     // Thoát vòng lặp khi i = 8
    }
}
```

### 6. Label với break và continue
- Định nghĩa label: `TênLabel:`
- Dùng với break/continue để chỉ định vòng lặp nào cần ảnh hưởng
- Hữu ích với các vòng lặp lồng nhau

**Ví dụ:**
```go
OuterLoop:
for i := 0; i < 5; i++ {
    for j := 0; j < 5; j++ {
        if điều_kiện {
            break OuterLoop  // Thoát cả 2 vòng lặp
        }
    }
}
```
