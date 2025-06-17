package main

import "fmt"

func main() {
	// 1. Khai báo biến
	var a int = 10        // Khai báo đầy đủ
	b := 20               // Khai báo ngắn gọn (:=)
	var c, d int = 30, 40 // Nhiều biến cùng kiểu
	var e int             // Biến chưa có giá trị (= 0)

	// 2. Các kiểu dữ liệu cơ bản
	var myString string = "Hello, Go!" // Chuỗi
	var myBool bool = true             // Boolean
	var myFloat float64 = 3.14         // Số thực
	var myRune rune = '国'              // Ký tự unicode (hiển thị mã số)

	// 3. Hằng số - giá trị không thể thay đổi
	const PI = 3.14159
	// PI = 3.15 // Lỗi: không thể thay đổi giá trị hằng số

	// In ra các giá trị
	fmt.Println("Biến số nguyên a =", a)
	fmt.Println("Biến số nguyên b =", b)
	fmt.Println("Biến c và d =", c, d)
	fmt.Println("Biến e (zero value) =", e)
	fmt.Println("Chuỗi:", myString)
	fmt.Println("Boolean:", myBool)
	fmt.Println("Số thực:", myFloat)
	fmt.Println("Ký tự Unicode:", myRune, "biểu diễn cho ký tự:", string(myRune))
	fmt.Println("Hằng số PI =", PI)
}
