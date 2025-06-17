package main

import (
	"fmt"
)

func main() {
	// 1. Khai báo và khởi tạo con trỏ
	var a int = 42  // Biến thông thường
	var p *int = &a // Con trỏ đến biến a

	fmt.Println("Giá trị của a:", a)
	fmt.Println("Địa chỉ của a:", &a)
	fmt.Println("Giá trị của p (địa chỉ mà p trỏ tới):", p)
	fmt.Println("Giá trị mà p đang trỏ tới:", *p)

	// 2. Thay đổi giá trị qua con trỏ
	*p = 100 // Thay đổi giá trị của a thông qua p
	fmt.Println("Giá trị của a sau khi thay đổi qua con trỏ:", a)

	// 3. Con trỏ nil
	var nilPtr *int // Con trỏ chưa trỏ tới đâu
	fmt.Println("Con trỏ nil:", nilPtr)

	// 4. Cấp phát bộ nhớ cho con trỏ với new()
	ptr := new(int) // Cấp phát bộ nhớ cho một biến int và trả về con trỏ
	*ptr = 50       // Gán giá trị
	fmt.Println("Giá trị mà ptr đang trỏ tới:", *ptr)

	// 5. Con trỏ đến con trỏ
	var pp **int = &p // Con trỏ đến con trỏ p
	fmt.Println("Giá trị của pp (địa chỉ mà pp trỏ tới):", pp)
	fmt.Println("Giá trị mà pp đang trỏ tới (địa chỉ mà p trỏ tới):", *pp)
	fmt.Println("Giá trị mà *pp đang trỏ tới (giá trị của a):", **pp)

	// 6. Sử dụng con trỏ với hàm
	x := 10
	fmt.Println("Trước khi gọi hàm increment:", x)
	increment(&x)
	fmt.Println("Sau khi gọi hàm increment:", x)

	// 7. Sử dụng con trỏ với struct
	person := Person{Name: "Alice", Age: 30}
	personPtr := &person

	fmt.Println("Person ban đầu:", person.Name, person.Age)

	// Có thể truy cập trực tiếp mà không cần dereference
	personPtr.Age = 31 // Tương đương với (*personPtr).Age = 31
	fmt.Println("Person sau khi thay đổi:", person.Name, person.Age)

	// 8. Con trỏ đến array vs slice
	// Slice đã là tham chiếu ngầm định, không cần dùng con trỏ
	numbers := []int{1, 2, 3, 4, 5}
	modifySlice(numbers)
	fmt.Println("Slice sau khi thay đổi:", numbers)

	// Array cần dùng con trỏ nếu muốn thay đổi
	arr := [5]int{10, 20, 30, 40, 50}
	modifyArray(&arr) // Truyền con trỏ đến array
	fmt.Println("Array sau khi thay đổi:", arr)

	// 9. So sánh con trỏ
	b := 42
	ptrA := &a
	ptrB := &b
	ptrC := &a

	fmt.Println("ptrA == ptrC:", ptrA == ptrC)     // true (cùng trỏ đến a)
	fmt.Println("ptrA == ptrB:", ptrA == ptrB)     // false (trỏ đến biến khác)
	fmt.Println("*ptrA == *ptrB:", *ptrA == *ptrB) // true (giá trị giống nhau)

	// 10. Sử dụng hàm trả về con trỏ
	newPerson := createPerson("Bob", 25)
	fmt.Println("Người mới tạo:", newPerson.Name, newPerson.Age)

	// 11. Sử dụng phương thức nhận con trỏ receiver
	newPerson.Birthday()
	fmt.Println("Sau sinh nhật:", newPerson.Name, newPerson.Age)

	// 12. Hoán đổi giá trị
	c, d := 10, 20
	fmt.Println("Trước khi swap:", c, d)
	swap(&c, &d)
	fmt.Println("Sau khi swap:", c, d)
}

// Hàm sử dụng con trỏ làm tham số
func increment(n *int) {
	*n = *n + 1 // Tăng giá trị mà n đang trỏ tới
}

// Struct biểu diễn thông tin người
type Person struct {
	Name string
	Age  int
}

// Hàm nhận slice (tham chiếu ngầm định)
func modifySlice(s []int) {
	if len(s) > 0 {
		s[0] = 999 // Thay đổi phần tử đầu tiên
	}
}

// Hàm nhận con trỏ đến array
func modifyArray(a *[5]int) {
	if len(a) > 0 {
		a[0] = 999 // Thay đổi phần tử đầu tiên
	}
}

// Ví dụ hàm trả về con trỏ
func createPerson(name string, age int) *Person {
	p := Person{
		Name: name,
		Age:  age,
	}
	return &p // Trả về con trỏ đến p (Go tự động xử lý, không lo lắng về vùng nhớ stack)
}

// Thêm ví dụ sử dụng con trỏ làm tham số receiver
func (p *Person) Birthday() {
	p.Age++ // Tăng tuổi
}

// Hàm swap dùng con trỏ để hoán đổi giá trị
func swap(a, b *int) {
	temp := *a
	*a = *b
	*b = temp
}
