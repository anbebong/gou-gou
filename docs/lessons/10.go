package main

import (
	"fmt"
	"sync"
)

// User struct để demo map với struct
type User struct {
	Name  string
	Age   int
	Email string
}

func main() {
	// 1. Khai báo Map
	fmt.Println("=== Các cách khai báo Map ===")

	var scores map[string]int // Map rỗng (nil)
	fmt.Println("scores (nil map):", scores)

	ages := make(map[string]int) // Khởi tạo bằng make
	fmt.Println("ages (empty map):", ages)

	colors := map[string]string{ // Khai báo và khởi tạo
		"red":   "#ff0000",
		"green": "#00ff00",
		"blue":  "#0000ff",
	}
	fmt.Println("colors:", colors)

	// 2. Thêm và cập nhật phần tử
	fmt.Println("\n=== Thêm và cập nhật phần tử ===")

	ages = make(map[string]int) // Khởi tạo map rỗng

	// Thêm phần tử
	ages["Alice"] = 25
	ages["Bob"] = 30
	fmt.Println("Sau khi thêm:", ages)

	// Cập nhật phần tử
	ages["Alice"] = 26
	fmt.Println("Sau khi cập nhật:", ages)

	// 3. Truy cập và kiểm tra tồn tại
	fmt.Println("\n=== Truy cập và kiểm tra tồn tại ===")

	// Truy cập trực tiếp
	aliceAge := ages["Alice"]
	fmt.Println("Tuổi của Alice:", aliceAge)

	// Kiểm tra sự tồn tại
	charlieAge, exists := ages["Charlie"]
	if exists {
		fmt.Println("Tuổi của Charlie:", charlieAge)
	} else {
		fmt.Println("Không tìm thấy Charlie trong map")
	}

	// 4. Xóa phần tử
	fmt.Println("\n=== Xóa phần tử ===")

	fmt.Println("Trước khi xóa:", ages)
	delete(ages, "Bob")
	fmt.Println("Sau khi xóa Bob:", ages)

	// 5. Duyệt qua map
	fmt.Println("\n=== Duyệt qua map ===")

	// Thêm dữ liệu mẫu
	colors["white"] = "#ffffff"
	colors["black"] = "#000000"

	// Duyệt key và value
	fmt.Println("Duyệt key và value:")
	for key, value := range colors {
		fmt.Printf("- Màu %s có mã màu: %s\n", key, value)
	}

	// Chỉ duyệt key
	fmt.Println("\nDuyệt chỉ key:")
	for key := range colors {
		fmt.Printf("- %s\n", key)
	}

	// 6. Map với struct
	fmt.Println("\n=== Map với struct ===")

	users := map[string]User{
		"u1": {
			Name:  "Alice",
			Age:   25,
			Email: "alice@example.com",
		},
		"u2": {
			Name:  "Bob",
			Age:   30,
			Email: "bob@example.com",
		},
	}

	fmt.Println("Users:", users)

	// 7. Đếm tần suất (use case phổ biến)
	fmt.Println("\n=== Đếm tần suất ===")

	words := []string{"apple", "banana", "apple", "cherry", "date", "banana", "apple"}
	wordCount := make(map[string]int)

	for _, word := range words {
		wordCount[word]++
	}

	fmt.Println("Tần suất xuất hiện của các từ:")
	for word, count := range wordCount {
		fmt.Printf("- %s: %d lần\n", word, count)
	}

	// 8. Map với concurrent access
	fmt.Println("\n=== Concurrent Map ===")

	// Sử dụng sync.Map cho concurrent access an toàn
	var concurrentMap sync.Map

	// Store
	concurrentMap.Store("key1", "value1")
	concurrentMap.Store("key2", "value2")

	// Load
	value, ok := concurrentMap.Load("key1")
	if ok {
		fmt.Printf("Giá trị của key1: %v\n", value)
	}

	// Range
	concurrentMap.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true
	})
}
