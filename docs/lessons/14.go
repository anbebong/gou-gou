package main

import (
	"fmt"
	"strings"
	"time"
)

// Khai báo một package level variable
var greeting = "Hello World"

// Khai báo constants
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// User struct để demo về exported fields
type User struct {
	Name   string // Exported field
	email  string // Unexported field
	Status string
}

// NewUser là một exported function
func NewUser(name, email string) *User {
	return &User{
		Name:   name,
		email:  email,
		Status: StatusActive,
	}
}

// GetEmail là method để truy cập unexported field
func (u *User) GetEmail() string {
	return u.email
}

// private function, chỉ dùng trong package này
func validateEmail(email string) bool {
	return strings.Contains(email, "@")
}

// Exported function sử dụng private function
func ValidateUser(u *User) bool {
	return validateEmail(u.email)
}

func main() {
	// 1. Demo package level variables và constants
	fmt.Println("=== Package Level Variables & Constants ===")
	fmt.Println("Greeting:", greeting)
	fmt.Println("Active status:", StatusActive)

	// 2. Demo exported vs unexported fields
	fmt.Println("\n=== Exported vs Unexported Fields ===")
	user := NewUser("John", "john@example.com")

	fmt.Println("Name (exported):", user.Name)
	// fmt.Println(user.email)  // Error: không thể truy cập unexported field
	fmt.Println("Email (through getter):", user.GetEmail())

	// 3. Demo các package phổ biến trong standard library
	fmt.Println("\n=== Standard Library Packages ===")

	// strings package
	message := "  hello, world  "
	fmt.Println("Original:", message)
	fmt.Println("Trimmed:", strings.TrimSpace(message))
	fmt.Println("Uppercase:", strings.ToUpper(message))
	fmt.Println("Contains 'world':", strings.Contains(message, "world"))

	// time package
	fmt.Println("\n=== Time Package ===")
	now := time.Now()
	fmt.Println("Current time:", now)
	fmt.Println("Formatted:", now.Format("2006-01-02 15:04:05"))

	// Thêm 2 ngày
	future := now.Add(48 * time.Hour)
	fmt.Println("After 2 days:", future.Format("2006-01-02"))

	// Tính khoảng thời gian
	duration := future.Sub(now)
	fmt.Printf("Duration: %.2f hours\n", duration.Hours())

	// 4. Demo validation
	fmt.Println("\n=== Validation ===")
	if ValidateUser(user) {
		fmt.Println("User email is valid")
	} else {
		fmt.Println("User email is invalid")
	}

	// 5. Demo multiple return values (phổ biến trong Go packages)
	fmt.Println("\n=== Multiple Return Values ===")
	if name, ok := getUserName(1); ok {
		fmt.Println("Found user:", name)
	}
	if name, ok := getUserName(2); ok {
		fmt.Println("Found user:", name)
	} else {
		fmt.Println("User not found")
	}
}

// Function trả về multiple values
func getUserName(id int) (string, bool) {
	users := map[int]string{
		1: "Alice",
	}
	name, exists := users[id]
	return name, exists
}
