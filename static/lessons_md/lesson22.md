# Templates trong Go

Go cung cấp package `html/template` để tạo và xử lý HTML templates một cách an toàn, tránh các vấn đề như XSS (Cross-Site Scripting).

## 1. Template Cơ bản

### Tạo và sử dụng template đơn giản:
```go
const tmpl = `
<h1>Xin chào {{.Name}}!</h1>
<p>Tuổi: {{.Age}}</p>
`

type Person struct {
    Name string
    Age  int
}

func main() {
    // Parse template
    t := template.Must(template.New("example").Parse(tmpl))
    
    // Data để render
    person := Person{
        Name: "Alice",
        Age:  25,
    }
    
    // Render template
    t.Execute(os.Stdout, person)
}
```

## 2. Template từ File

### Cấu trúc thư mục:
```
myapp/
  ├── templates/
  │   ├── header.html
  │   ├── footer.html
  │   └── home.html
  └── main.go
```

### File templates:
```html
<!-- header.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>

<!-- footer.html -->
</body>
</html>

<!-- home.html -->
{{template "header.html" .}}
<h1>Xin chào {{.Name}}!</h1>
<ul>
    {{range .Hobbies}}
    <li>{{.}}</li>
    {{end}}
</ul>
{{template "footer.html"}}
```

### Sử dụng trong code:
```go
type PageData struct {
    Title   string
    Name    string
    Hobbies []string
}

func main() {
    // Load tất cả templates
    tmpl := template.Must(template.ParseGlob("templates/*.html"))
    
    // Data cho template
    data := PageData{
        Title:   "Trang chủ",
        Name:    "Alice",
        Hobbies: []string{"Đọc sách", "Chơi game", "Code"},
    }
    
    // Render template "home.html"
    tmpl.ExecuteTemplate(os.Stdout, "home.html", data)
}
```

## 3. Template Functions

### Functions có sẵn:
- and, or, not: Phép logic
- eq, ne, lt, le, gt, ge: So sánh
- index: Truy cập phần tử của array/slice
- len: Độ dài của array/slice/map
- print, printf, println: In ra output

### Custom functions:
```go
funcMap := template.FuncMap{
    "upper": strings.ToUpper,
    "formatDate": func(t time.Time) string {
        return t.Format("02-01-2006")
    },
}

tmpl := template.New("test").Funcs(funcMap)
```

Sử dụng trong template:
```html
<h1>{{upper .Title}}</h1>
<p>Ngày: {{formatDate .Date}}</p>
```

## 4. Điều kiện và Vòng lặp

### If-else:
```html
{{if .IsAdmin}}
    <a href="/admin">Admin Panel</a>
{{else if .IsUser}}
    <a href="/profile">Profile</a>
{{else}}
    <a href="/login">Login</a>
{{end}}
```

### Range:
```html
<ul>
{{range $index, $element := .Items}}
    <li>{{$index}}: {{$element}}</li>
{{else}}
    <li>Không có phần tử nào</li>
{{end}}
</ul>
```

## Bài tập

1. Blog Template:
   - Tạo template cho blog với header, content, footer
   - Hiển thị danh sách bài viết
   - Template riêng cho trang chi tiết bài viết
   - Thêm phân trang

2. Admin Dashboard:
   - Layout chung với sidebar
   - Form thêm/sửa bài viết
   - Bảng quản lý users
   - Hiển thị thông báo lỗi/thành công

3. Email Template:
   - Template cho email xác nhận đăng ký
   - Template cho email reset password
   - Sử dụng partial templates cho header/footer
   - Đảm bảo responsive
