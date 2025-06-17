# Struct (Cấu trúc) trong Go

Struct là một kiểu dữ liệu tự định nghĩa cho phép gom nhóm các trường dữ liệu có liên quan với nhau. Struct trong Go tương tự như class trong các ngôn ngữ OOP khác, nhưng đơn giản hơn.

## Định nghĩa Struct

```go
type Person struct {
    Name    string
    Age     int
    Email   string
    Address string
}
```

## Khởi tạo Struct

Có nhiều cách để khởi tạo một struct:

```go
// Cách 1: Khởi tạo trống
var person1 Person

// Cách 2: Khởi tạo với các giá trị theo thứ tự
person2 := Person{"Alice", 25, "alice@example.com", "123 Street"}

// Cách 3: Khởi tạo với tên trường (recommended)
person3 := Person{
    Name:    "Bob",
    Age:     30,
    Email:   "bob@example.com",
    Address: "456 Avenue",
}

// Cách 4: Sử dụng new (trả về con trỏ)
person4 := new(Person)
```

## Truy cập và thay đổi trường dữ liệu

```go
person := Person{Name: "Alice", Age: 25}

// Truy cập trường
fmt.Println(person.Name)  // "Alice"

// Thay đổi giá trị trường
person.Age = 26
```

## Struct lồng nhau

Struct có thể chứa struct khác làm trường dữ liệu:

```go
type Address struct {
    Street  string
    City    string
    Country string
}

type Person struct {
    Name    string
    Age     int
    Address Address    // Struct lồng nhau
}
```

## Con trỏ tới Struct

```go
person := &Person{Name: "Alice"}
// Hai cách sau đều đúng:
fmt.Println((*person).Name)  // Dereference rồi truy cập
fmt.Println(person.Name)     // Go tự động dereference
```

## Struct với Tag

Tags cho phép thêm metadata cho các trường của struct:

```go
type User struct {
    Name     string `json:"name"`
    Age      int    `json:"age" validate:"gte=0,lte=130"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"-"`  // Dấu - nghĩa là bỏ qua trường này khi serialize
}
```

## Anonymous Struct

Go cho phép tạo struct mà không cần định nghĩa type:

```go
point := struct {
    X int
    Y int
}{
    X: 10,
    Y: 20,
}
```

## Anonymous Fields (Embedding)

Go hỗ trợ embedding để thực hiện tính kế thừa đơn giản:

```go
type Animal struct {
    Name string
    Age  int
}

type Dog struct {
    Animal      // Anonymous field (embedding)
    Breed string
}
```

## Zero Value của Struct

- Các trường số: 0
- Các trường string: ""
- Các trường con trỏ: nil
- Các trường struct lồng nhau: zero value của struct đó

## Best Practices

1. Quy ước đặt tên:
   - Tên struct và trường public: viết hoa (exported)
   - Tên trường private: viết thường (unexported)

```go
type Company struct {
    Name      string  // public
    location  string  // private
    Employees int
}
```

2. Sử dụng named fields khi khởi tạo:
```go
person := Person{
    Name: "Alice",
    Age:  25,
}
```

3. Sử dụng con trỏ cho struct lớn:
```go
func UpdatePerson(p *Person) {
    p.Age++
}
```

4. Tổ chức code:
```go
// person.go
type Person struct {
    // fields
}

// Các hàm liên quan đến Person
func (p Person) GetFullName() string {
    // ...
}
```

5. Validation:
```go
func (p Person) Validate() error {
    if p.Name == "" {
        return errors.New("name is required")
    }
    if p.Age < 0 {
        return errors.New("age must be positive")
    }
    return nil
}
```

## Ứng dụng thực tế

1. Mô hình hóa dữ liệu:
```go
type Product struct {
    ID          int
    Name        string
    Price       float64
    Description string
    CreatedAt   time.Time
}
```

2. Configuration:
```go
type Config struct {
    Server struct {
        Host string
        Port int
    }
    Database struct {
        URL      string
        Username string
        Password string
    }
}
```

3. API Request/Response:
```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token   string `json:"token"`
    Expires string `json:"expires"`
}
```
