package main

import (
	"fmt"
	"time"
)

func timerExample(t time.Duration, done chan bool) {
	time.Sleep(t)
	done <- true
}

func doTask(name string, duration time.Duration) {
	fmt.Printf("Task %s starting at: %s\n", name, time.Now().Format("15:04:05"))
	time.Sleep(duration)
	fmt.Printf("Task %s completed at: %s\n", name, time.Now().Format("15:04:05"))
}

func main() {
	done := make(chan bool)
	fmt.Println("Starting timer at:", time.Now().Format("15:04:05"))

	// Chạy timer 3s
	go timerExample(3*time.Second, done)

	// Chạy task mất 2s
	go doTask("Processing", 4*time.Second)

	// Đợi signal từ timerExample
	<-done
	fmt.Println("Timer complete at:", time.Now().Format("15:04:05"))
	// Tổng thời gian thực hiện sẽ là 3s
}
