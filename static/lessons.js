const lessons = [
    {
        title: "Cấu trúc chương trình Go",
        desc: "Tìm hiểu cấu trúc cơ bản của một chương trình Go.",
        mdPath: "lessons_md/lesson1.md",
        codePath: "lessons/1.go",
        tasks: [
            {
                title: "Hello World",
                desc: "Viết chương trình Go đầu tiên in ra \"Hello, World!\"",
                hint: "Sử dụng package main và func main()",
                level: "easy"
            },
            {
                title: "Multiple Prints",
                desc: "Viết chương trình in nhiều dòng text khác nhau sử dụng Println và Printf",
                hint: "Thử các cách định dạng khác nhau với fmt.Printf",
                level: "easy"
            },
            {
                title: "Comment Types",
                desc: "Thực hành sử dụng các loại comment khác nhau trong code",
                hint: "Có 2 loại comment: // và /* */",
                level: "easy"
            }
        ]
    },
    {
        title: "Biến, hằng số và kiểu dữ liệu cơ bản",
        desc: "Khai báo biến, hằng số và các kiểu dữ liệu cơ bản trong Go.",
        mdPath: "lessons_md/lesson2.md",
        codePath: "lessons/2.go",
        tasks: [
            {
                title: "Khai báo biến",
                desc: "Thực hành các cách khai báo biến khác nhau",
                hint: "Thử var, :=, và khai báo nhiều biến",
                level: "easy"
            },
            {
                title: "Kiểu dữ liệu cơ bản",
                desc: "Tạo biến với các kiểu dữ liệu khác nhau và in ra giá trị",
                hint: "Dùng các kiểu int, float64, string, bool",
                level: "easy"
            },
            {
                title: "Hằng số",
                desc: "Khai báo và sử dụng các hằng số",
                hint: "Dùng const để khai báo hằng số",
                level: "easy"
            },
            {
                title: "Type Conversion",
                desc: "Thực hành chuyển đổi giữa các kiểu dữ liệu",
                hint: "Dùng int(), float64(), string()",
                level: "medium"
            }
        ]
    },
    {
        title: "Toán tử và biểu thức",
        desc: "Các loại toán tử và cách sử dụng biểu thức trong Go.",
        mdPath: "lessons_md/lesson3.md",
        codePath: "lessons/3.go"
    },
    {
        title: "Câu lệnh điều kiện (if, switch)",
        desc: "Sử dụng if, else, switch để điều khiển luồng chương trình.",
        mdPath: "lessons_md/lesson4.md",
        codePath: "lessons/4.go"
    },
    {
        title: "Vòng lặp (for, break, continue)",
        desc: "Cách sử dụng vòng lặp for, break, continue trong Go.",
        mdPath: "lessons_md/lesson6.md",
        codePath: "lessons/6.go"
    },
    { title: "Hàm (function), tham số, giá trị trả về", desc: "Định nghĩa và sử dụng hàm, truyền tham số, giá trị trả về.", mdPath: "lessons_md/lesson7.md", codePath: "lessons/7.go" },
    { 
        title: "Con trỏ (pointer)", 
        desc: "Khái niệm và cách sử dụng con trỏ trong Go.", 
        mdPath: "lessons_md/lesson8.md",
        codePath: "lessons/8.go" 
    },
    { title: "Mảng (array), slice", desc: "Làm việc với mảng và slice trong Go.", mdPath: "lessons_md/lesson9.md", codePath: "lessons/9.go" },
    { title: "Map (bản đồ)", desc: "Cách sử dụng map (bản đồ) trong Go.", mdPath: "lessons_md/lesson10.md", codePath: "lessons/10.go" },
    { title: "Struct (cấu trúc)", desc: "Định nghĩa và sử dụng struct trong Go.", mdPath: "lessons_md/lesson11.md", codePath: "lessons/11.go" },
    { title: "Method (phương thức)", desc: "Định nghĩa và sử dụng method cho struct.", mdPath: "lessons_md/lesson12.md", codePath: "lessons/12.go" },
    { title: "Interface", desc: "Khái niệm và cách sử dụng interface trong Go.", mdPath: "lessons_md/lesson13.md", codePath: "lessons/13.go" },
    { title: "Package & import", desc: "Tổ chức code với package và import.", mdPath: "lessons_md/lesson14.md", codePath: "lessons/14.go" },
    { 
        title: "Xử lý lỗi (error handling)", 
        desc: "Cách xử lý lỗi trong Go.", 
        mdPath: "lessons_md/lesson15.md", 
        codePath: "lessons/15.go",
        tasks: [
            {
                title: "Error Cơ bản",
                desc: "Tạo và xử lý error đơn giản",
                hint: "Dùng errors.New() hoặc fmt.Errorf()",
                level: "easy"
            },
            {
                title: "Custom Error Type",
                desc: "Tạo custom error type với nhiều thông tin hơn",
                hint: "Tạo struct implement Error interface",
                level: "medium"
            },
            {
                title: "Multiple Errors",
                desc: "Xử lý nhiều lỗi cùng lúc",
                hint: "Tạo struct chứa slice các error",
                level: "medium" 
            },
            {
                title: "Error Wrapping",
                desc: "Wrap error với context bổ sung",
                hint: "Dùng fmt.Errorf() với %w",
                level: "medium"
            },
            {
                title: "Clean Error Handling",
                desc: "Viết hàm xử lý lỗi sạch sẽ với defer",
                hint: "Kết hợp defer và named return values",
                level: "hard"
            }
        ]
    },
    { 
        title: "Goroutine (lập trình song song)", 
        desc: "Giới thiệu về goroutine và concurrency.", 
        mdPath: "lessons_md/lesson16.md", 
        codePath: "lessons/16.go",
        tasks: [
            {
                title: "Goroutine cơ bản",
                desc: "Tạo và chạy một goroutine đơn giản",
                hint: "Sử dụng từ khóa go và WaitGroup",
                level: "easy"
            },
            {
                title: "Multiple Goroutines",
                desc: "Chạy nhiều goroutines cùng lúc và đồng bộ kết quả",
                hint: "Dùng WaitGroup để đợi tất cả goroutines hoàn thành",
                level: "medium"
            }
        ]
    },
    { 
        title: "Channel", 
        desc: "Truyền thông tin giữa các goroutine bằng channel.", 
        mdPath: "lessons_md/lesson17.md", 
        codePath: "lessons/17.go",
        tasks: [
            {
                title: "Channel cơ bản",
                desc: "Tạo channel và truyền dữ liệu giữa goroutines",
                hint: "Dùng make(chan Type) và phép toán <-",
                level: "easy"
            },
            {
                title: "Buffer Channel",
                desc: "Thực hành với buffered channel",
                hint: "Tạo channel với buffer size > 0",
                level: "medium"
            },
            {
                title: "Worker Pool",
                desc: "Xây dựng worker pool sử dụng channels",
                hint: "Dùng 2 channels: jobs và results",
                level: "medium"
            },
            {
                title: "Fan-out Fan-in",
                desc: "Implement pattern fan-out/fan-in với channels",
                hint: "Nhiều workers xử lý từ một input channel",
                level: "hard"
            },
            {
                title: "Chat System",
                desc: "Tạo chat system đơn giản sử dụng channels",
                hint: "Broadcast messages cho nhiều clients",
                level: "hard"
            }
        ]
    },
    { 
        title: "Select", 
        desc: "Sử dụng select để xử lý nhiều channel.", 
        mdPath: "lessons_md/lesson18.md", 
        codePath: "lessons/18.go",
        tasks: [
            {
                title: "Select cơ bản",
                desc: "Sử dụng select để đọc từ nhiều channels",
                hint: "Dùng select với nhiều case channels",
                level: "easy"
            },
            {
                title: "Timeout Handler",
                desc: "Implement timeout cho operations bằng select",
                hint: "Dùng time.After() trong select",
                level: "medium"
            },
            {
                title: "Rate Limiter",
                desc: "Tạo rate limiter sử dụng select và ticker",
                hint: "Kết hợp time.Ticker với select",
                level: "medium"
            },
            {
                title: "Message Router",
                desc: "Tạo router điều hướng messages từ nhiều channels",
                hint: "Select với priority handling",
                level: "hard"
            },
            {
                title: "Pub/Sub System",
                desc: "Implement hệ thống publish/subscribe",
                hint: "Quản lý nhiều subscribers với select",
                level: "hard"
            }
        ]
    },    { 
        title: "Quản lý module (go mod)", 
        desc: "Giới thiệu về go mod và quản lý module.", 
        mdPath: "lessons_md/lesson19.md", 
        codePath: "lessons/19.go",
        tasks: [
            {
                title: "Basic Module",
                desc: "Tạo module đơn giản với một package",
                hint: "Tạo go.mod và package với 2-3 functions",
                level: "easy"
            },
            {
                title: "Import Package",
                desc: "Import và sử dụng package từ module khác",
                hint: "Dùng go get để thêm dependency",
                level: "easy"
            },
            {
                title: "Multiple Packages",
                desc: "Tạo module với nhiều packages",
                hint: "Chia code thành các package logic",
                level: "medium"
            },
            {
                title: "Version Migration",
                desc: "Nâng cấp/hạ cấp version của dependency",
                hint: "Dùng go get pkg@version",
                level: "medium"
            },
            {
                title: "Project Structure",
                desc: "Tổ chức project với internal và pkg",
                hint: "Phân chia code public/private",
                level: "hard"
            }
        ]
    },
    { title: "Đọc/ghi file", desc: "Cách đọc và ghi file trong Go.", mdPath: "lessons_md/lesson20.md", codePath: "lessons/20.go", 
        tasks: [
            {
                title: "Đọc file log",
                desc: "Tạo file log.txt với một số dòng log",
                hint: "Viết chương trình đọc và in ra màn hình các dòng có chữ ERROR",
                level: "easy"
            },
            {
                title: "Ghi file cấu hình",
                desc: "Tạo struct Config lưu các thiết lập",
                hint: "Viết hàm SaveConfig() lưu cấu hình xuống file, LoadConfig() đọc cấu hình từ file",
                level: "easy"
            },
            {
                title: "Sao chép file",
                desc: "Viết chương trình sao chép nội dung từ file nguồn sang file đích",
                hint: "Hiển thị tiến trình sao chép (phần trăm hoàn thành)",
                level: "medium"
            },
        ]
     },
    { title: "Làm việc với HTTP (client/server)", desc: "Tạo HTTP client và server cơ bản.", mdPath: "lessons_md/lesson21.md", codePath: "lessons/21.go" },
    { title: "JSON & encoding", desc: "Xử lý dữ liệu JSON và encoding trong Go.", mdPath: "lessons_md/lesson22.md", codePath: "lessons/22.go" },
    { title: "Unit test trong Go", desc: "Viết và chạy unit test trong Go.", mdPath: "lessons_md/lesson23.md", codePath: "lessons/23.go" },
    { title: "Context", desc: "Sử dụng context để kiểm soát goroutine.", mdPath: "lessons_md/lesson24.md", codePath: "lessons/24.go" },
    { title: "Reflection", desc: "Khái niệm và ứng dụng reflection trong Go.", mdPath: "lessons_md/lesson25.md", codePath: "lessons/25.go" },
    { title: "Xây dựng ứng dụng web đơn giản với net/http", desc: "Tạo web app đơn giản sử dụng net/http.", mdPath: "lessons_md/lesson26.md", codePath: "lessons/26.go" },
    { title: "Sử dụng thư viện ngoài (third-party)", desc: "Cài đặt và sử dụng thư viện ngoài.", mdPath: "lessons_md/lesson27.md", codePath: "lessons/27.go" },
    { title: "Tổng kết & tài nguyên học thêm", desc: "Tổng kết và gợi ý tài liệu học nâng cao.", mdPath: "lessons_md/lesson28.md", codePath: "lessons/28.go" }
];
