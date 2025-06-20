package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Worker xử lý số với timeout 2 giây
func worker(id int, jobs <-chan int, results chan<- string) {
	for job := range jobs {
		fmt.Printf("Worker %d started job %d at %s\n",
			id, job, time.Now().Format("15:04:05"))
		// Giả lập công việc mất thời gian ngẫu nhiên từ 1 đến 3 giây
		time.Sleep(1 * time.Second)

		// Sử dụng select để handle timeout
		select {
		case <-done:
			results <- fmt.Sprintf("Worker %d completed job %d successfully", id, job)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	jobs := make(chan int, 10)
	results := make(chan string, 1)

	// Khởi tạo 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Gửi 10 jobs
	for j := 1; j <= 10; j++ {
		jobs <- j
	} // Gửi tín hiệu kết thúc cho workers
	close(jobs)

	// biết rằng có 10 job
	for range 10 {
		result := <-results
		fmt.Println(result)
	}
	// todo nếu như không biết có bao nhiêu job thì làm sao để biết khi nào kết thúc?
}
