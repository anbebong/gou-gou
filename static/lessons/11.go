package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Address struct để demo struct lồng nhau
type Address struct {
	Street  string
	City    string
	Country string
}

// Person struct với nhiều loại trường
type Person struct {
	Name     string
	Age      int
	Email    string
	Address  Address   // Struct lồng nhau
	JoinedAt time.Time // Sử dụng package time
}

// User struct với tags
type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Sẽ không xuất hiện trong JSON
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Animal struct để demo embedding
type Animal struct {
	Name string
	Age  int
}

// Method cho Animal
func (a Animal) Description() string {
	return fmt.Sprintf("%s is %d years old", a.Name, a.Age)
}

// Dog struct với embedding
type Dog struct {
	Animal        // Embedding Animal
	Breed  string // Thêm trường mới
}

func main() {
	// 1. Các cách khởi tạo struct
	fmt.Println("=== Khởi tạo Struct ===")

	// Khởi tạo trống
	var person1 Person
	fmt.Printf("Person1 (zero value): %+v\n", person1)

	// Khởi tạo với giá trị
	person2 := Person{
		Name:  "Alice",
		Age:   25,
		Email: "alice@example.com",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
		JoinedAt: time.Now(),
	}
	fmt.Printf("\nPerson2: %+v\n", person2)

	// 2. Struct lồng nhau
	fmt.Println("\n=== Struct lồng nhau ===")

	// Truy cập trường của struct lồng nhau
	fmt.Printf("Address: %s, %s, %s\n",
		person2.Address.Street,
		person2.Address.City,
		person2.Address.Country)

	// 3. Struct với JSON
	fmt.Println("\n=== Struct với JSON ===")

	user := User{
		Username:  "johndoe",
		Password:  "secret123", // Sẽ không xuất hiện trong JSON
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	// Chuyển struct thành JSON
	jsonData, _ := json.MarshalIndent(user, "", "    ")
	fmt.Printf("User as JSON:\n%s\n", string(jsonData))

	// 4. Embedding
	fmt.Println("\n=== Embedding ===")

	dog := Dog{
		Animal: Animal{
			Name: "Max",
			Age:  3,
		},
		Breed: "Golden Retriever",
	}

	// Truy cập trường và method được embedded
	fmt.Println("Dog name:", dog.Name)             // Từ Animal
	fmt.Println("Dog breed:", dog.Breed)           // Từ Dog
	fmt.Println("Description:", dog.Description()) // Method từ Animal

	// 5. Con trỏ tới struct
	fmt.Println("\n=== Con trỏ tới Struct ===")

	personPtr := &person2
	fmt.Printf("Name (through pointer): %s\n", personPtr.Name)

	// Thay đổi giá trị qua con trỏ
	personPtr.Age = 26
	fmt.Printf("Updated age: %d\n", person2.Age)

	// 6. Anonymous struct
	fmt.Println("\n=== Anonymous Struct ===")

	point := struct {
		X int
		Y int
	}{
		X: 10,
		Y: 20,
	}
	fmt.Printf("Point: %+v\n", point)

	// 7. Slice của struct
	fmt.Println("\n=== Slice của Struct ===")

	people := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}

	fmt.Println("People:")
	for _, p := range people {
		fmt.Printf("- %s (%d years old)\n", p.Name, p.Age)
	}

	// 8. Map với struct
	fmt.Println("\n=== Map với Struct ===")

	employees := map[string]Person{
		"emp1": {Name: "David", Age: 28},
		"emp2": {Name: "Eve", Age: 32},
	}

	fmt.Println("Employees:")
	for id, emp := range employees {
		fmt.Printf("- ID: %s, Name: %s, Age: %d\n", id, emp.Name, emp.Age)
	}
}
