package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// PageData chứa dữ liệu cho template
type PageData struct {
	Title    string
	Name     string
	Hobbies  []string
	IsAdmin  bool
	Date     time.Time
	Articles []Article
}

// Article đại diện cho một bài viết
type Article struct {
	Title    string
	Content  string
	Author   string
	PostDate time.Time
}

func main() {
	// Custom template functions
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("02-01-2006")
		},
	}

	// Load và parse templates
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	// Sample data
	data := PageData{
		Title:   "Trang Blog",
		Name:    "Alice",
		IsAdmin: true,
		Date:    time.Now(),
		Hobbies: []string{"Đọc sách", "Code", "Du lịch"},
		Articles: []Article{
			{
				Title:    "Bài viết 1",
				Content:  "Nội dung bài viết 1...",
				Author:   "Alice",
				PostDate: time.Now().Add(-24 * time.Hour),
			},
			{
				Title:    "Bài viết 2",
				Content:  "Nội dung bài viết 2...",
				Author:   "Bob",
				PostDate: time.Now(),
			},
		},
	}

	// Handler cho trang chủ
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "home.html", data)
	})

	// Handler cho trang admin
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "admin.html", data)
	})

	// Start server
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
