package main

import (
	"context"
	"sync"
)

// Job represents an interface of a job that can be enqueued into a dispatcher.
type Job struct {
	Payload string
}

type Worker interface {
	Work(j *Job) // Work method defines the behavior of processing a job.
}

// Dispatcher represents a job dispatcher.
type Dispatcher struct {
	workerPool chan struct{}  // Semaphore for limiting concurrent worker goroutines.
	jobQueue   chan *Job      // Channel for queuing incoming jobs.
	worker     Worker         // Worker interface for processing jobs.
	globalWg   sync.WaitGroup // WaitGroup for tracking the termination of the dispatcher.
}

// NewDispatcher creates a new instance of a job dispatcher with the given parameters.
func NewDispatcher(worker Worker, maxWorkers int, maxQueueSize int) *Dispatcher {
	return &Dispatcher{
		workerPool: make(chan struct{}, maxWorkers), // Buffered channel acting as a workerPool. Use empty struct to minimize the memory allocation
		jobQueue:   make(chan *Job, maxQueueSize),
		worker:     worker,
	}
}

// Start initiates the dispatcher to begin processing jobs.
// The dispatcher stops when it receives a value from `ctx.Done`.
func (d *Dispatcher) Start(ctx context.Context) {
	// Increment the wait group counter to indicate that the dispatcher has started processing jobs.
	d.globalWg.Add(1)

	var wg sync.WaitGroup

	// Main loop for processing jobs.
	for {
		select {
		case <-ctx.Done():
			// Block until all currently processing jobs have finished.
			wg.Wait()
			// When the loop exits (due to context cancellation), stop decrements the wait group counter to indicate that the dispatcher has stopped processing jobs
			d.globalWg.Done()
			return
		case job := <-d.jobQueue:
			// Increment the local wait group to track the processing of this job.
			wg.Add(1)
			// Push to the workerPool to control the number of concurrent workers.
			d.workerPool <- struct{}{}
			// Process the job concurrently.
			go func(job *Job) {
				defer wg.Done()
				// After the job finishes, release the slot in the workerPool.
				defer func() { <-d.workerPool }()
				d.worker.Work(job)
			}(job)
		}
	}
}

// Wait blocks until the dispatcher stops.
func (d *Dispatcher) Wait() {
	d.globalWg.Wait()
}

// Enqueue puts a job into the queue.
// If the number of enqueued jobs has already reached the maximum size,
// This will block until space becomes available in the queue to accept a new job.
func (d *Dispatcher) Enqueue(job *Job) {
	d.jobQueue <- job
}
