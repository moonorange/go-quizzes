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
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		// wait until receiving the signal
		<-sigCh
		cancel()
	}()

	p := NewAPIWorker()
	d := NewDispatcher(p, 10, 1000)
	d.Start(ctx)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(maxNumWorkers)
	// concurrentEnqueue(d, &wg)

	enqueue(d, 0, totalJobs)

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All enqueue jobs completed.")

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

func concurrentEnqueue(d *Dispatcher, wg *sync.WaitGroup) {
	// Launch goroutines
	for i := 0; i < maxNumWorkers; i++ {
		start := i * jobsPerRoutine
		end := start + jobsPerRoutine
		if i == maxNumWorkers-1 {
			end = totalJobs // Handle remaining jobs for the last goroutine
		}
		fmt.Printf("start: %d end: %d\n", start, end)
		go func() {
			defer wg.Done()
			enqueue(d, start, end)
		}()
	}
}
