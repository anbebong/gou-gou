package main

import (
	"fmt"
	"time"
)

func main() {
	// Câu lệnh if-else
	x := 8

	if x > 10 {
		fmt.Println("x lớn hơn 10")
	} else if x > 5 {
		fmt.Println("x lớn hơn 5 nhưng nhỏ hơn hoặc bằng 10")
	} else {
		fmt.Println("x nhỏ hơn hoặc bằng 5")
	}

	// If với khai báo biến ngắn gọn
	if y := 2 * x; y > 10 {
		fmt.Printf("y (= %d) lớn hơn 10\n", y)
	} else {
		fmt.Printf("y (= %d) nhỏ hơn hoặc bằng 10\n", y)
	}
	// Biến y chỉ tồn tại trong phạm vi if-else

	// Câu lệnh switch
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Saturday, time.Sunday:
		fmt.Println("Cuối tuần rồi!")
	default:
		fmt.Println("Ngày trong tuần :(")
	}

	// Switch không cần điều kiện
	hour := time.Now().Hour()
	switch {
	case hour < 12:
		fmt.Println("Chào buổi sáng!")
	case hour < 18:
		fmt.Println("Chào buổi chiều!")
	default:
		fmt.Println("Chào buổi tối!")
	}

	// Switch với fallthrough
	num := 6
	switch {
	case num < 10:
		fmt.Println("Số nhỏ hơn 10")
		fallthrough
	case num < 20:
		fmt.Println("Số nhỏ hơn 20")
	case num < 30:
		fmt.Println("Số nhỏ hơn 30")
	}
}
