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
        codePath: "lessons/3.go",
        tasks: [
            {
                title: "Toán tử so sánh",
                desc: "Viết chương trình sử dụng các toán tử so sánh (==, !=, <, >, <=, >=) với hai số. In kết quả (true/false) của mỗi phép so sánh.",
                hint: "Kết quả của toán tử so sánh là một giá trị boolean.",
                level: "easy"
            },
            {
                title: "Toán tử logic",
                desc: "Viết chương trình sử dụng các toán tử logic (&&, ||, !) với các biểu thức boolean. In kết quả.",
                hint: "Tạo các biến boolean hoặc sử dụng kết quả từ các phép so sánh.",
                level: "medium"
            },
            {
                title: "Toán tử gán",
                desc: "Thực hành sử dụng các toán tử gán kết hợp (+=, -=, *=, /=, %=).",
                hint: "Ví dụ: x += 5 tương đương với x = x + 5.",
                level: "medium"
            },
            {
                title: "Độ ưu tiên toán tử",
                desc: "Viết một biểu thức phức tạp kết hợp nhiều loại toán tử và dự đoán kết quả. Sau đó, chạy chương trình để kiểm tra.",
                hint: "Sử dụng dấu ngoặc đơn () để thay đổi độ ưu tiên nếu cần.",
                level: "hard"
            }
        ]
    },
    {
        title: "Câu lệnh điều kiện (if, switch)",
        desc: "Sử dụng if, else, switch để điều khiển luồng chương trình.",
        mdPath: "lessons_md/lesson4.md",
        codePath: "lessons/4.go",
        tasks: [
            {
                title: "Kiểm tra số chẵn/lẻ",
                desc: "Viết chương trình nhận một số nguyên và kiểm tra xem đó là số chẵn hay số lẻ sử dụng câu lệnh if-else. In kết quả ra màn hình.",
                hint: "Sử dụng toán tử modulo (%) để kiểm tra tính chẵn lẻ.",
                level: "easy"
            },
            {
                title: "Phân loại tuổi",
                desc: "Viết chương trình nhận vào một độ tuổi và phân loại: 'Trẻ em' (dưới 13), 'Thiếu niên' (13-17), 'Người lớn' (18-64), 'Người cao tuổi' (65 trở lên) sử dụng if-else if-else.",
                hint: "Sử dụng nhiều điều kiện else if.",
                level: "easy"
            },
            {
                title: "Switch với ngày trong tuần",
                desc: "Viết chương trình nhận một số từ 1 đến 7 và in ra ngày tương ứng trong tuần (1 là Chủ Nhật, 2 là Thứ Hai,...). Sử dụng câu lệnh switch.",
                hint: "Mỗi case sẽ tương ứng với một số.",
                level: "medium"
            },
            {
                title: "Switch với fallthrough",
                desc: "Sử dụng switch với từ khóa fallthrough để xử lý trường hợp nhiều case có cùng một hành động. Ví dụ: kiểm tra một ký tự có phải là nguyên âm hay không (a, e, i, o, u).",
                hint: "Sau một case, sử dụng fallthrough để thực thi tiếp case tiếp theo.",
                level: "medium"
            },
            {
                title: "If với khởi tạo biến",
                desc: "Thực hành câu lệnh if có phần khởi tạo biến ngắn gọn (ví dụ: if err := someFunc(); err != nil { ... }). Viết một hàm trả về lỗi và kiểm tra lỗi đó.",
                hint: "Biến được khởi tạo trong if chỉ có phạm vi trong khối if/else đó.",
                level: "hard"
            }
        ]
    },
    {
        title: "Vòng lặp (for, break, continue)",
        desc: "Cách sử dụng vòng lặp for, break, continue trong Go.",
        mdPath: "lessons_md/lesson6.md",
        codePath: "lessons/6.go",
        tasks: [
            {
                title: "Tính tổng các số",
                desc: "Viết chương trình tính tổng các số từ 1 đến N (N được nhập vào).",
                hint: "Khởi tạo một biến tổng và cộng dồn trong vòng lặp.",
                level: "easy"
            },
            {
                title: "Vẽ hình chữ nhật bằng dấu *",
                desc: "Viết chương trình nhận vào chiều rộng và chiều cao, sau đó vẽ một hình chữ nhật bằng các dấu '*' sử dụng vòng lặp lồng nhau.",
                hint: "Sử dụng một vòng lặp for cho hàng và một vòng lặp for lồng bên trong cho cột.",
                level: "hard"
            }
        ]
    },
    {
        title: "Hàm (function), tham số, giá trị trả về", 
        desc: "Định nghĩa và sử dụng hàm, truyền tham số, giá trị trả về.", 
        mdPath: "lessons_md/lesson7.md", 
        codePath: "lessons/7.go", 
        tasks: [
            {
                title: "Hàm tính diện tích hình chữ nhật",
                desc: "Viết một hàm nhận vào chiều dài và chiều rộng, sau đó trả về diện tích của hình chữ nhật.",
                hint: "Hàm sẽ có hai tham số kiểu số và trả về một giá trị kiểu số.",
                level: "easy"
            },
            {
                title: "Hàm với nhiều giá trị trả về",
                desc: "Viết một hàm nhận vào hai số nguyên a và b, sau đó trả về cả tổng và hiệu của chúng (a+b và a-b).",
                hint: "Go cho phép hàm trả về nhiều giá trị.",
                level: "medium"
            },
            {
                title: "Hàm variadic (tham số thay đổi)",
                desc: "Viết một hàm tính tổng của một danh sách các số nguyên. Hàm này có thể nhận vào một số lượng tham số bất kỳ (variadic parameter).",
                hint: "Sử dụng ...int cho tham số variadic và duyệt qua nó như một slice.",
                level: "medium"
            },
            {
                title: "Hàm đệ quy: Giai thừa",
                desc: "Viết một hàm đệ quy để tính giai thừa của một số nguyên không âm N (N!).",
                hint: "Giai thừa của 0 là 1. Giai thừa của N là N * (N-1)!. Cẩn thận với trường hợp cơ sở.",
                level: "hard"
            }
        ]
    },
    {
        title: "Con trỏ (pointer)",
        desc: "Khái niệm và cách sử dụng con trỏ trong Go.",
        mdPath: "lessons_md/lesson8.md",
        codePath: "lessons/8.go",
        tasks: [
            {
                title: "Thay đổi giá trị qua con trỏ",
                desc: "Khai báo một biến. Tạo một con trỏ trỏ đến biến đó. Thay đổi giá trị của biến gốc thông qua con trỏ. In giá trị của biến gốc để xác nhận.",
                hint: "Sử dụng *pointer = newValue.",
                level: "easy"
            },
            {
                title: "Con trỏ và hàm",
                desc: "Viết một hàm nhận vào một con trỏ đến một số nguyên và thay đổi giá trị của số nguyên đó bên trong hàm. Gọi hàm và kiểm tra xem giá trị có thay đổi bên ngoài hàm không.",
                hint: "Hàm có tham số kiểu *int.",
                level: "medium"
            },
            {
                title: "Con trỏ nil",
                desc: "Khai báo một con trỏ nhưng không khởi tạo cho nó trỏ đến đâu cả (con trỏ nil). Thử giải tham chiếu con trỏ nil và quan sát lỗi (panic). Sau đó, thêm kiểm tra xem con trỏ có nil không trước khi giải tham chiếu.",
                hint: "Một con trỏ chưa được gán sẽ có giá trị nil. Kiểm tra if p != nil.",
                level: "medium"
            },
            {
                title: "Con trỏ tới con trỏ",
                desc: "Khai báo một biến, một con trỏ trỏ tới biến đó, và một con trỏ thứ hai trỏ tới con trỏ đầu tiên. In ra giá trị của biến gốc thông qua con trỏ thứ hai.",
                hint: "Sử dụng **p để giải tham chiếu hai lần.",
                level: "hard"
            }
        ]
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
    }, {
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
    {
        title: "Đọc/ghi file", desc: "Cách đọc và ghi file trong Go.", mdPath: "lessons_md/lesson20.md", codePath: "lessons/20.go",
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
    { title: "JSON & encoding", desc: "Xử lý dữ liệu JSON và encoding trong Go.", mdPath: "lessons_md/lesson22.md", codePath: "lessons/22.go" }
];
