package main

import (
	"fmt"
	"math"
)

// Hàm không có tham số và không trả về giá trị
func sayHello() {
	fmt.Println("Hello, Go!")
}

// Hàm có tham số
func greet(name string) {
	fmt.Println("Hello,", name)
}

// Hàm trả về một giá trị
func add(a, b int) int {
	return a + b
}

// Hàm trả về nhiều giá trị
func calculate(a, b int) (int, int, int, int) {
	sum := a + b
	diff := a - b
	product := a * b
	quotient := a / b
	return sum, diff, product, quotient
}

// Hàm với kết quả trả về được đặt tên
func divide(a, b float64) (quotient float64, err error) {
	if b == 0 {
		err = fmt.Errorf("không thể chia cho 0")
		return
	}
	quotient = a / b
	return // Tự động trả về quotient và err
}

// Hàm đệ quy
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// Hàm nhận tham số là một hàm khác (function as parameter)
func applyOperation(a, b int, operation func(int, int) int) int {
	return operation(a, b)
}

// Hàm trả về một hàm khác (function as return value)
func getOperation(opName string) func(int, int) int {
	switch opName {
	case "add":
		return func(a, b int) int { return a + b }
	case "subtract":
		return func(a, b int) int { return a - b }
	case "multiply":
		return func(a, b int) int { return a * b }
	default:
		return func(a, b int) int { return 0 }
	}
}

// Hàm với số lượng tham số không xác định (variadic function)
func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

func main() {
	// Gọi hàm không tham số
	sayHello()

	// Gọi hàm có tham số
	greet("Gopher")

	// Gọi hàm trả về một giá trị
	result := add(5, 3)
	fmt.Println("5 + 3 =", result)

	// Gọi hàm trả về nhiều giá trị
	sumResult, diff, product, quotient := calculate(10, 5)
	fmt.Println("10 + 5 =", sumResult)
	fmt.Println("10 - 5 =", diff)
	fmt.Println("10 * 5 =", product)
	fmt.Println("10 / 5 =", quotient)

	// Chỉ lấy một số giá trị trả về, bỏ qua các giá trị khác
	sum2, _, _, _ := calculate(20, 10)
	fmt.Println("20 + 10 =", sum2)

	// Gọi hàm với kết quả trả về được đặt tên
	result2, err := divide(10, 2)
	if err != nil {
		fmt.Println("Lỗi:", err)
	} else {
		fmt.Println("10 / 2 =", result2)
	}

	// Thử với phép chia cho 0
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Lỗi:", err)
	}

	// Gọi hàm đệ quy
	fmt.Println("5! =", factorial(5))

	// Sử dụng hàm như tham số
	multiply := func(a, b int) int { return a * b }
	result3 := applyOperation(6, 7, multiply)
	fmt.Println("6 * 7 =", result3)

	// Sử dụng hàm ẩn danh (anonymous function) làm tham số
	result4 := applyOperation(6, 7, func(a, b int) int {
		return int(math.Pow(float64(a), float64(b)))
	})
	fmt.Println("6^7 =", result4)

	// Lấy hàm từ hàm khác trả về
	subtractFunc := getOperation("subtract")
	fmt.Println("15 - 7 =", subtractFunc(15, 7))
	// Sử dụng variadic function
	fmt.Println("Tổng của 1, 2, 3, 4, 5 =", sum(1, 2, 3, 4, 5))
}
