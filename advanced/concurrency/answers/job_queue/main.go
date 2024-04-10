package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	totalJobs      int = 100
	maxNumWorkers      = 100
	jobsPerRoutine     = totalJobs / maxNumWorkers
)

func main() {
	// Create a cancellation context to allow graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to receive OS signals
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	// Notify sigCh channel
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		// Wait until receiving the signal
		<-sigCh
		// Cancel the context to propagate cancellation through the context tree
		cancel()
	}()

	p := NewAPIWorker()
	d := NewDispatcher(p, 10, 1000)
	// Start the dispatcher with the cancellation context
	d.Start(ctx)

	// WaitGroup to wait for all goroutines to finish
	// var wg sync.WaitGroup
	// wg.Add(maxNumWorkers)
	// concurrentEnqueue(d, &wg)

	enqueue(d, 0, totalJobs)

	// Wait for all goroutiwnes to finish
	// wg.Wait()

	fmt.Println("All enqueue jobs completed.")

	// Wait for the dispatcher to finish processing all jobs
	d.Wait()
	fmt.Println("Finished!")
}

// enqueue enqueues jobs into the dispatcher within the specified range
func enqueue(d *Dispatcher, start, end int) {
	for i := start; i < end; i++ {
		payload := fmt.Sprintf("dummy_payload_%d", i)
		job := &Job{Payload: payload}
		d.Enqueue(job)
	}
}

// concurrentEnqueue enqueues jobs concurrently using multiple goroutines
func concurrentEnqueue(d *Dispatcher, wg *sync.WaitGroup) {
	for i := 0; i < maxNumWorkers; i++ {
		start := i * jobsPerRoutine
		end := start + jobsPerRoutine
		// Handle remaining jobs for the last goroutine
		if i == maxNumWorkers-1 {
			end = totalJobs
		}
		fmt.Printf("start: %d end: %d\n", start, end)
		// Launch a goroutine to enqueue jobs within the specified range
		go func() {
			defer wg.Done()
			enqueue(d, start, end)
		}()
	}
}
