package main

import (
	"fmt"
	"time"
)

// 1. Worker Pool Example
func workerPoolDemo() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// Start 3 workers
	for w := 1; w <= 3; w++ {
		go poolWorker(w, jobs, results)
	}

	// Send 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for a := 1; a <= 5; a++ {
		fmt.Printf("Result: %d\n", <-results)
	}
}

func poolWorker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(100 * time.Millisecond) // Simulate work
		results <- j * 2
	}
}

// 2. Pipeline Example
func pipeline() {
	// Create channels
	numbers := make(chan int)
	doubles := make(chan int)
	done := make(chan bool)

	// Stage 1: Generate numbers
	go func() {
		for i := 1; i <= 3; i++ {
			numbers <- i
		}
		close(numbers)
	}()

	// Stage 2: Double numbers
	go func() {
		for n := range numbers {
			doubles <- n * 2
		}
		close(doubles)
	}()

	// Stage 3: Print results
	go func() {
		for d := range doubles {
			fmt.Printf("Pipeline result: %d\n", d)
		}
		done <- true
	}()

	<-done
}

// 3. Buffered vs Unbuffered Example
func bufferedExample() {
	unbuffered := make(chan int)
	buffered := make(chan int, 2)

	// With unbuffered channel
	go func() {
		fmt.Println("Sending to unbuffered channel...")
		unbuffered <- 1
		fmt.Println("Sent to unbuffered channel")
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Received from unbuffered: %d\n", <-unbuffered)

	// With buffered channel
	fmt.Println("Sending to buffered channel...")
	buffered <- 1
	buffered <- 2
	fmt.Println("Sent to buffered channel")

	fmt.Printf("Received from buffered: %d\n", <-buffered)
	fmt.Printf("Received from buffered: %d\n", <-buffered)
}

// 4. Timeout Example
func timeoutExample() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Done!"
	}()

	select {
	case result := <-ch:
		fmt.Println("Received:", result)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout!")
	}
}

// 5. Fan-out/Fan-in Example
func fanOutFanIn() {
	input := make(chan int, 100)
	output := make(chan int, 100)

	// Fan-out to 3 workers
	for i := 0; i < 3; i++ {
		go fanOutWorker(input, output)
	}

	// Send work
	go func() {
		for i := 1; i <= 9; i++ {
			input <- i
		}
		close(input)
	}()

	// Collect all results
	for i := 1; i <= 9; i++ {
		fmt.Printf("Fan-in result: %d\n", <-output)
	}
}

func fanOutWorker(input <-chan int, output chan<- int) {
	for n := range input {
		output <- n * n // Square the number
		time.Sleep(100 * time.Millisecond)
	}
}

func demoChannels() {
	fmt.Println("=== 1. Worker Pool Example ===")
	workerPoolDemo()

	fmt.Println("\n=== 2. Pipeline Example ===")
	pipeline()

	fmt.Println("\n=== 3. Buffered vs Unbuffered Example ===")
	bufferedExample()

	fmt.Println("\n=== 4. Timeout Example ===")
	timeoutExample()

	fmt.Println("\n=== 5. Fan-out/Fan-in Example ===")
	fanOutFanIn()
}

func main() {
	demoChannels()
}
