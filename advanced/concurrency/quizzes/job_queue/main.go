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

	p := NewWorker()
	d := NewDispatcher(p, 10, 100)
	// Start the dispatcher with the cancellation context
	go d.Start(ctx)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(maxNumWorkers)
	enqueue(d, 0, totalJobs)

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All enqueue jobs completed.")

	// Wait for the dispatcher to finish processing all jobs
	d.Wait()
	fmt.Println("Finished!")
}

func enqueue(d *Dispatcher, start, end int) {
	for i := start; i < end; i++ {
		payload := fmt.Sprintf("dummy_payload_%d", i)
		job := &Job{Payload: payload}
		d.Enqueue(job)
	}
}
