package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Ví dụ 1: Goroutine cơ bản
func sayHello(s string, times int) {
	for i := 0; i < times; i++ {
		fmt.Printf("Hello from %s\n", s)
		time.Sleep(100 * time.Millisecond)
	}
}

// Ví dụ 2: Counter với race condition
var counter = 0

func incrementWithRace() {
	temp := counter
	time.Sleep(1 * time.Microsecond) // Giả lập xử lý
	counter = temp + 1
}

// Ví dụ 3: Counter với mutex
var (
	counterSafe = 0
	mutex       sync.Mutex
)

func incrementSafe() {
	mutex.Lock()
	defer mutex.Unlock()
	counterSafe++
}

// Ví dụ 4: Worker Pool
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
		time.Sleep(100 * time.Millisecond) // Giả lập xử lý
		fmt.Printf("Worker %d finished job %d\n", id, j)
		results <- j * 2 // Gửi kết quả
	}
}

// Ví dụ 5: Fan-out pattern
func process(item int) int {
	time.Sleep(50 * time.Millisecond) // Giả lập xử lý
	return item * item
}

func main() {
	fmt.Println("=== 1. Goroutine cơ bản ===")
	go sayHello("Goroutine 1", 3)
	go sayHello("Goroutine 2", 3)
	time.Sleep(1 * time.Second) // Đợi goroutines hoàn thành

	fmt.Println("\n=== 2. Race Condition Demo ===")
	for i := 0; i < 5; i++ {
		go incrementWithRace()
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Counter với race: %d\n", counter)

	fmt.Println("\n=== 3. Mutex Demo ===")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incrementSafe()
		}()
	}
	wg.Wait()
	fmt.Printf("Counter với mutex: %d\n", counterSafe)

	fmt.Println("\n=== 4. Worker Pool Demo ===")
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// Khởi tạo 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Gửi 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Nhận kết quả
	for a := 1; a <= 5; a++ {
		<-results
	}

	fmt.Println("\n=== 5. Fan-out Pattern Demo ===")
	items := []int{1, 2, 3, 4, 5}
	numWorkers := 3
	itemsChan := make(chan int)
	resultsChan := make(chan int)

	// Khởi tạo workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for item := range itemsChan {
				resultsChan <- process(item)
			}
		}()
	}

	// Gửi items
	go func() {
		for _, item := range items {
			itemsChan <- item
		}
		close(itemsChan)
	}()

	// Nhận kết quả
	var results []int
	for i := 0; i < len(items); i++ {
		result := <-resultsChan
		results = append(results, result)
	}
	fmt.Printf("Kết quả xử lý: %v\n", results)

	fmt.Println("\n=== 6. Goroutine Cleanup Demo ===")
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Goroutine cleanup và thoát")
				return
			default:
				fmt.Println("Goroutine đang chạy...")
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// Để goroutine chạy một lúc
	time.Sleep(600 * time.Millisecond)

	// Signal goroutine dừng lại
	close(done)

	// Đợi goroutine cleanup
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\nSố goroutines còn đang chạy:", runtime.NumGoroutine())
}
