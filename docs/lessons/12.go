package main

import (
	"fmt"
	"strings"
)

// Shape interface để demo method với interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle struct và các method của nó
type Rectangle struct {
	width  float64
	height float64
}

// Value receiver - không thay đổi receiver
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

// Pointer receiver - thay đổi receiver
func (r *Rectangle) Scale(factor float64) {
	r.width *= factor
	r.height *= factor
}

// Person struct để demo nhiều loại method
type Person struct {
	Name     string
	Age      int
	Email    string
	Password string
}

// Value receiver method
func (p Person) GetInfo() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

// Pointer receiver method
func (p *Person) SetAge(age int) {
	p.Age = age
}

// Validation methods
func (p Person) ValidateEmail() error {
	if !strings.Contains(p.Email, "@") {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func (p Person) ValidatePassword() error {
	if len(p.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

// String method (implementing Stringer interface)
func (p Person) String() string {
	return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

// Method với non-struct type
type MyInt int

func (m MyInt) Double() MyInt {
	return m * 2
}

func (m MyInt) IsPositive() bool {
	return m > 0
}

// Builder pattern với method chaining
type StringBuilder struct {
	str string
}

func (b *StringBuilder) Append(s string) *StringBuilder {
	b.str += s
	return b
}

func (b *StringBuilder) AppendLine(s string) *StringBuilder {
	b.str += s + "\n"
	return b
}

func (b StringBuilder) String() string {
	return b.str
}

// Query builder để demo fluent interface
type Query struct {
	sql string
}

func (q *Query) Select(columns ...string) *Query {
	q.sql = "SELECT " + strings.Join(columns, ", ")
	return q
}

func (q *Query) From(table string) *Query {
	q.sql += " FROM " + table
	return q
}

func (q *Query) Where(condition string) *Query {
	q.sql += " WHERE " + condition
	return q
}

func (q Query) String() string {
	return q.sql
}

func main() {
	// 1. Rectangle với value và pointer receivers
	fmt.Println("=== Rectangle Methods ===")
	rect := Rectangle{width: 10, height: 5}

	fmt.Printf("Original: width=%.2f, height=%.2f\n", rect.width, rect.height)
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())

	rect.Scale(2)
	fmt.Printf("\nAfter scaling: width=%.2f, height=%.2f\n", rect.width, rect.height)
	fmt.Printf("New area: %.2f\n", rect.Area())

	// 2. Person với nhiều loại method
	fmt.Println("\n=== Person Methods ===")
	person := Person{
		Name:     "Alice",
		Age:      25,
		Email:    "alice@example.com",
		Password: "secretpass123",
	}

	fmt.Println("Info:", person.GetInfo())
	person.SetAge(26)
	fmt.Println("After SetAge:", person)

	// Validation
	if err := person.ValidateEmail(); err != nil {
		fmt.Println("Email error:", err)
	}
	if err := person.ValidatePassword(); err != nil {
		fmt.Println("Password error:", err)
	}

	// 3. Method với non-struct type
	fmt.Println("\n=== MyInt Methods ===")
	num := MyInt(5)
	fmt.Printf("Original: %d\n", num)
	fmt.Printf("Doubled: %d\n", num.Double())
	fmt.Printf("Is positive? %v\n", num.IsPositive())

	// 4. String builder với method chaining
	fmt.Println("\n=== StringBuilder with Method Chaining ===")
	builder := StringBuilder{}
	result := builder.
		Append("Hello").
		Append(" ").
		Append("World").
		AppendLine("!").
		String()

	fmt.Println(result)

	// 5. Query builder với fluent interface
	fmt.Println("\n=== Query Builder with Fluent Interface ===")
	query := Query{}
	sql := query.
		Select("id", "name", "email").
		From("users").
		Where("age >= 18").
		String()

	fmt.Println(sql)
}
