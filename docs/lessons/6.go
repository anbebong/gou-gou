package main

import "fmt"

func main() {
	// For với 3 thành phần
	fmt.Println("For với 3 thành phần:")
	for i := 1; i <= 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// For với điều kiện (như while trong các ngôn ngữ khác)
	fmt.Println("For với điều kiện:")
	count := 1
	for count <= 5 {
		fmt.Printf("%d ", count)
		count++
	}
	fmt.Println()

	// For vô hạn với break
	fmt.Println("For vô hạn với break:")
	sum := 0
	for {
		sum++
		if sum > 5 {
			break
		}
		fmt.Printf("%d ", sum)
	}
	fmt.Println()

	// For với range cho mảng/slice
	fmt.Println("For với range cho slice:")
	numbers := []int{1, 2, 3, 4, 5}
	for i, num := range numbers {
		fmt.Printf("numbers[%d] = %d\n", i, num)
	}

	// For với range cho map
	fmt.Println("For với range cho map:")
	person := map[string]string{
		"name": "Gopher",
		"age":  "10",
		"lang": "Go",
	}
	for key, value := range person {
		fmt.Printf("%s: %s\n", key, value)
	}

	// Continue: bỏ qua số chẵn
	fmt.Println("For với continue (chỉ in số lẻ):")
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// For lồng nhau với label và break
	fmt.Println("For lồng nhau với label:")
OuterLoop:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if i*j > 5 {
				break OuterLoop
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}
