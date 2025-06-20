package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Data int
}

type Result struct {
	JobID     int
	WorkerID  int
	Result    int
	StartTime time.Time
	EndTime   time.Time
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Ghi nhận thời gian bắt đầu
		startTime := time.Now()

		// Giả lập xử lý job
		time.Sleep(time.Duration(job.Data%3+1) * time.Second)

		// Tạo result
		result := Result{
			JobID:     job.ID,
			WorkerID:  id,
			Result:    job.Data * 2, // giả sử công việc là nhân 2
			StartTime: startTime,
			EndTime:   time.Now(),
		}

		// Gửi kết quả (có thể block nếu result buffer đầy)
		results <- result
		fmt.Printf("Worker %d completed job %d at %s\n",
			id, job.ID, time.Now().Format("15:04:05"))
	}
}

func main() {
	N := 20 // Số lượng jobs
	X := 5  // Số lượng workers
	Y := 2  // Kích thước buffer của result channel (Y < X)

	// Tạo channels
	jobs := make(chan Job, N)
	results := make(chan Result, Y) // Buffer giới hạn cho results

	var wg sync.WaitGroup

	// Khởi tạo workers
	fmt.Printf("Starting %d workers...\n", X)
	for i := 1; i <= X; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// Gửi jobs
	fmt.Printf("Sending %d jobs...\n", N)
	go func() {
		for i := 1; i <= N; i++ {
			jobs <- Job{
				ID:   i,
				Data: i,
			}
		}
		close(jobs)
	}()

	// Goroutine để đợi tất cả workers hoàn thành và đóng results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Nhận và xử lý results
	completedJobs := 0
	startTime := time.Now()

	for result := range results {
		completedJobs++
		processTime := result.EndTime.Sub(result.StartTime)

		fmt.Printf("Result: Job %d processed by Worker %d in %v\n",
			result.JobID, result.WorkerID, processTime)

		if completedJobs == N {
			break
		}
	}

	totalTime := time.Since(startTime)
	fmt.Printf("\nAll %d jobs completed in %v\n", N, totalTime)
	fmt.Printf("Average time per job: %v\n", totalTime/time.Duration(N))
}
