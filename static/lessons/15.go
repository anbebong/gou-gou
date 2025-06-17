package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// 1. Custom Error Types
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// 2. Sentinel Errors
var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid input")
)

// 3. Multiple Errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var msgs []string
	for _, err := range v {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// 4. Utility Functions
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("không thể chia cho 0")
	}
	return a / b, nil
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("đọc file %s: %w", path, err)
	}
	return string(data), nil
}

// 5. Validation
type Validator struct {
	errors ValidationErrors
}

func (v *Validator) Check(valid bool, field, msg string) *Validator {
	if !valid {
		v.errors = append(v.errors, ValidationError{field, msg})
	}
	return v
}

func (v *Validator) Errors() error {
	if len(v.errors) == 0 {
		return nil
	}
	return v.errors
}

func main() {
	fmt.Println("=== Error Handling Demo ===")

	// 1. Basic error
	fmt.Println("\n1. Basic error:")
	if result, err := divide(10, 0); err != nil {
		fmt.Printf("Lỗi: %v\n", err)
	} else {
		fmt.Printf("Kết quả: %v\n", result)
	}

	// 2. Error wrapping
	fmt.Println("\n2. Error wrapping:")
	if _, err := readFile("config.txt"); err != nil {
		fmt.Printf("Original error: %v\n", err)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File không tồn tại")
		}
	}

	// 3. Custom error
	fmt.Println("\n3. Validation errors:")
	user := struct {
		Name  string
		Email string
		Age   int
	}{"", "invalid", -1}

	v := &Validator{}
	err := v.Check(user.Name != "", "name", "required").
		Check(strings.Contains(user.Email, "@"), "email", "invalid format").
		Check(user.Age >= 0, "age", "must be positive").
		Errors()

	if err != nil {
		fmt.Printf("Lỗi validation: %v\n", err)
		// Type assertion
		if valErrs, ok := err.(ValidationErrors); ok {
			fmt.Printf("Số lỗi: %d\n", len(valErrs))
		}
	}

	// 4. Panic recovery
	fmt.Println("\n4. Panic recovery:")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from: %v\n", r)
			}
		}()
		panic("something bad happened")
	}()
}
