# Toán tử và biểu thức trong Go

Go cung cấp nhiều loại toán tử để thực hiện các phép tính và biểu thức logic.

### 1. Toán tử số học
| Toán tử | Mô tả | Ví dụ |
|---------|-------|-------|
| + | Cộng | 5 + 3 = 8 |
| - | Trừ | 5 - 3 = 2 |
| * | Nhân | 5 * 3 = 15 |
| / | Chia | 5 / 3 = 1 (với số nguyên) |
| % | Chia lấy dư | 5 % 3 = 2 |
| ++ | Tăng thêm 1 | a++ |
| -- | Giảm đi 1 | a-- |

> **Lưu ý**: Trong Go, ++ và -- là câu lệnh, không phải biểu thức. Nghĩa là bạn không thể viết b = a++ như trong C/Java.

### 2. Toán tử so sánh
| Toán tử | Mô tả | Ví dụ |
|------|-------|-------|
| == | Bằng | a == b |
| != | Khác | a != b |
| < | Nhỏ hơn | a < b |
| <= | Nhỏ hơn hoặc bằng | a <= b |
| > | Lớn hơn | a > b |
| >= | Lớn hơn hoặc bằng | a >= b |

Tất cả các toán tử so sánh đều trả về kiểu **bool** (true/false).

### 3. Toán tử logic
| Toán tử | Mô tả | Ví dụ |
|---------|-------|-------|
| && | VÀ logic | a && b (cả hai điều kiện phải đúng) |
| || | HOẶC logic | a || b (ít nhất một điều kiện đúng) |
| ! | PHỦĐỊNH | !a (đúng thành sai, sai thành đúng) |

Go sử dụng **đánh giá ngắn mạch** (short-circuit evaluation): Khi đánh giá a && b, nếu a đã sai, sẽ không cần đánh giá b.

### 4. Toán tử bit
| Toán tử | Mô tả | Ví dụ |
|---------|-------|-------|
| & | AND bit | 5 & 3 = 1 |
| | | OR bit | 5 | 3 = 7 |
| ^ | XOR bit | 5 ^ 3 = 6 |
| << | Dịch trái | 5 << 1 = 10 |
| >> | Dịch phải | 5 >> 1 = 2 |
| &^ | Bit Clear (AND NOT) | 5 &^ 3 = 4 |

Các phép toán bit rất hữu ích khi làm việc với cờ, mặt nạ bit và tối ưu hóa.

### 5. Toán tử gán
| Toán tử | Mô tả | Tương đương |
|---------|-------|------------|
| = | Gán giá trị | a = 5 |
| += | Cộng rồi gán | a += 2 ⟹ a = a + 2 |
| -= | Trừ rồi gán | a -= 2 ⟹ a = a - 2 |
| *= | Nhân rồi gán | a *= 2 ⟹ a = a * 2 |
| /= | Chia rồi gán | a /= 2 ⟹ a = a / 2 |
| %= | Chia lấy dư rồi gán | a %= 2 ⟹ a = a % 2 |
| &= | AND bit rồi gán | a &= 2 ⟹ a = a & 2 |
| |= | OR bit rồi gán | a |= 2 ⟹ a = a | 2 |
| ^= | XOR bit rồi gán | a ^= 2 ⟹ a = a ^ 2 |
| <<= | Dịch trái rồi gán | a <<= 2 ⟹ a = a << 2 |
| >>= | Dịch phải rồi gán | a >>= 2 ⟹ a = a >> 2 |

### 6. Quy tắc ưu tiên
Các toán tử được thực hiện theo thứ tự ưu tiên từ cao xuống thấp:
1. **Cao nhất**: (), [], . (nhóm, truy cập phần tử, truy cập trường)
2. *, /, %, <<, >>, &, &^
3. +, -, |, ^
4. ==, !=, <, <=, >, >=
5. &&
6. **Thấp nhất**: ||

```go
// Ví dụ về thứ tự ưu tiên
result := 5 + 3*2        // 5 + (3*2) = 11, không phải (5+3)*2 = 16
check := x > 10 && y < 5 // (x > 10) && (y < 5)
```
