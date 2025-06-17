package main

import (
	"context"
	"fmt"
	"time"
)

// 1. Basic Select Example
func basicSelect(ch1, ch2 chan string) {
	select {
	case msg1 := <-ch1:
		fmt.Println("Received from ch1:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received from ch2:", msg2)
	default:
		fmt.Println("No data available")
	}
}

// 2. Timeout Pattern
func timeoutPattern(ch chan string) {
	select {
	case result := <-ch:
		fmt.Println("Received result:", result)
	case <-time.After(2 * time.Second):
		fmt.Println("Operation timed out")
	}
}

// 3. Context Cancellation
func contextCancellation(ctx context.Context, ch chan string) {
	select {
	case <-ctx.Done():
		fmt.Println("Operation cancelled:", ctx.Err())
	case result := <-ch:
		fmt.Println("Received result:", result)
	}
}

// 4. Fan-in Pattern
func fanIn(ch1, ch2 <-chan string) <-chan string {
	merged := make(chan string)
	go func() {
		defer close(merged)
		for {
			select {
			case v1 := <-ch1:
				merged <- v1
			case v2 := <-ch2:
				merged <- v2
			}
		}
	}()
	return merged
}

// 5. Rate Limiting
func rateLimiter() {
	requests := make(chan int, 5)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	// Producer
	go func() {
		for i := 1; i <= 10; i++ {
			requests <- i
		}
		close(requests)
	}()

	// Consumer with rate limiting
	for req := range requests {
		<-ticker.C // Wait for tick
		fmt.Printf("Processing request %d\n", req)
	}
}

// 6. Graceful Shutdown
func gracefulShutdown() {
	work := make(chan string)
	quit := make(chan bool)

	// Worker
	go func() {
		for {
			select {
			case <-quit:
				fmt.Println("Worker: Cleaning up...")
				fmt.Println("Worker: Shutting down")
				return
			case w := <-work:
				fmt.Println("Working on:", w)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Simulate work
	work <- "Task 1"
	work <- "Task 2"

	// Signal shutdown
	quit <- true
	time.Sleep(time.Second)
}

func demoSelect() {
	// 1. Basic Select
	fmt.Println("=== Basic Select ===")
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "Hello from ch1"
	}()
	basicSelect(ch1, ch2)

	// 2. Timeout Pattern
	fmt.Println("\n=== Timeout Pattern ===")
	ch := make(chan string)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "Delayed response"
	}()
	timeoutPattern(ch)

	// 3. Context Cancellation
	fmt.Println("\n=== Context Cancellation ===")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	contextCancellation(ctx, ch)

	// 4. Fan-in Pattern
	fmt.Println("\n=== Fan-in Pattern ===")
	input1 := make(chan string)
	input2 := make(chan string)
	merged := fanIn(input1, input2)

	go func() {
		input1 <- "Input 1"
		input2 <- "Input 2"
	}()

	for i := 0; i < 2; i++ {
		fmt.Println("Merged:", <-merged)
	}

	// 5. Rate Limiting
	fmt.Println("\n=== Rate Limiting ===")
	rateLimiter()

	// 6. Graceful Shutdown
	fmt.Println("\n=== Graceful Shutdown ===")
	gracefulShutdown()
}

func main() {
	demoSelect()
}
