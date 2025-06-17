package main

import (
	"fmt"
)

func main() {
	// 1. Arrays
	fmt.Println("--- Arrays ---")
	var nums [5]int
	fmt.Printf("Empty array: %v\n", nums)

	fruits := [3]string{"apple", "banana", "orange"}
	fmt.Printf("Fruit array: %v\n", fruits)

	// Array comparison
	arr1 := [3]int{1, 2, 3}
	arr2 := [3]int{1, 2, 3}
	fmt.Printf("arr1 == arr2: %v\n", arr1 == arr2)

	// 2. Slices
	fmt.Println("\n--- Slices ---")
	// Creating slices
	slice1 := []int{1, 2, 3, 4, 5}
	fmt.Printf("slice1: %v\n", slice1)

	slice2 := make([]int, 3, 5)
	fmt.Printf("slice2 length: %d, capacity: %d\n", len(slice2), cap(slice2))

	// Slicing
	fmt.Println("\n--- Slicing ---")
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("numbers[1:4]: %v\n", numbers[1:4])
	fmt.Printf("numbers[:3]: %v\n", numbers[:3])
	fmt.Printf("numbers[2:]: %v\n", numbers[2:])

	// Appending
	fmt.Println("\n--- Append ---")
	s := []int{1, 2}
	fmt.Printf("Original slice: %v\n", s)

	s = append(s, 3)
	fmt.Printf("After append(s, 3): %v\n", s)

	s = append(s, 4, 5, 6)
	fmt.Printf("After append(s, 4, 5, 6): %v\n", s)

	// Combining slices
	s1 := []int{1, 2}
	s2 := []int{3, 4}
	s1 = append(s1, s2...)
	fmt.Printf("Combined slices: %v\n", s1)

	// 3. Demonstrating nil vs empty slice
	fmt.Println("\n--- Nil vs Empty Slice ---")
	var nilSlice []int
	emptySlice := []int{}

	fmt.Printf("nilSlice == nil: %v\n", nilSlice == nil)
	fmt.Printf("emptySlice == nil: %v\n", emptySlice == nil)

	// 4. Copy slices
	fmt.Println("\n--- Copy Slices ---")
	original := []int{1, 2, 3}
	copied := make([]int, len(original))
	copy(copied, original)

	// Modify original to show they're independent
	original[0] = 99
	fmt.Printf("Original after modification: %v\n", original)
	fmt.Printf("Copied slice (unchanged): %v\n", copied)

	// 5. Demonstrating slice capacity growth
	fmt.Println("\n--- Capacity Growth ---")
	s = make([]int, 0)
	fmt.Printf("Initial cap: %d\n", cap(s))

	for i := 0; i < 10; i++ {
		s = append(s, i)
		fmt.Printf("len=%d\tcap=%d\n", len(s), cap(s))
	}
}
