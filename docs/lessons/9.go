package main

import "fmt"

func main() {
	// 1. Arrays
	fmt.Println("=== Arrays ===")

	// Khai báo và khởi tạo mảng
	var numbers [5]int                          // Mảng số nguyên kích thước 5
	colors := [3]string{"red", "green", "blue"} // Khai báo và khởi tạo
	matrix := [2][3]int{{1, 2, 3}, {4, 5, 6}}   // Mảng 2 chiều

	// In ra các mảng
	fmt.Println("Numbers (default):", numbers)
	fmt.Println("Colors:", colors)
	fmt.Println("Matrix:", matrix)

	// Truy cập và thay đổi phần tử
	numbers[0] = 10
	numbers[1] = 20
	fmt.Println("Numbers (after modify):", numbers)

	// 2. Slices
	fmt.Println("\n=== Slices ===")

	// Các cách tạo slice
	var s1 []int               // Slice rỗng (nil)
	s2 := []int{1, 2, 3, 4, 5} // Khởi tạo trực tiếp
	s3 := make([]int, 3, 5)    // Make với length 3, capacity 5

	fmt.Println("s1 (nil slice):", s1, "length:", len(s1), "capacity:", cap(s1))
	fmt.Println("s2:", s2, "length:", len(s2), "capacity:", cap(s2))
	fmt.Println("s3:", s3, "length:", len(s3), "capacity:", cap(s3))

	// Slicing
	fmt.Println("\n=== Slicing ===")
	arr := [5]int{1, 2, 3, 4, 5}
	slice1 := arr[1:4] // Từ index 1 đến 3
	slice2 := arr[:3]  // Từ đầu đến index 2
	slice3 := arr[2:]  // Từ index 2 đến cuối
	slice4 := arr[:]   // Toàn bộ mảng

	fmt.Println("Original array:", arr)
	fmt.Println("slice1 (arr[1:4]):", slice1)
	fmt.Println("slice2 (arr[:3]):", slice2)
	fmt.Println("slice3 (arr[2:]):", slice3)
	fmt.Println("slice4 (arr[:]):", slice4)

	// append và copy
	fmt.Println("\n=== Append và Copy ===")
	s := []int{1, 2, 3}
	fmt.Println("Original slice:", s)

	// Append
	s = append(s, 4)
	fmt.Println("After append(s, 4):", s)
	s = append(s, 5, 6, 7)
	fmt.Println("After append(s, 5, 6, 7):", s)

	// Copy
	dest := make([]int, len(s))
	copied := copy(dest, s)
	fmt.Println("Copied slice:", dest, "number of elements copied:", copied)

	// Demo capacity growth
	fmt.Println("\n=== Capacity Growth ===")
	slice := make([]int, 0)
	fmt.Printf("Initial: len=%d cap=%d\n", len(slice), cap(slice))

	for i := 0; i < 10; i++ {
		slice = append(slice, i)
		fmt.Printf("After append %d: len=%d cap=%d\n", i, len(slice), cap(slice))
	}

	// Demo slice sharing underlying array
	fmt.Println("\n=== Slice Sharing ===")
	original := []int{1, 2, 3, 4, 5}
	shared := original[1:4]

	fmt.Println("Before modification:")
	fmt.Println("original:", original)
	fmt.Println("shared:", shared)

	shared[0] = 20 // Thay đổi phần tử đầu tiên của shared (index 1 của original)

	fmt.Println("\nAfter modification:")
	fmt.Println("original:", original)
	fmt.Println("shared:", shared)
}
