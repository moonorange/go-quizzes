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
	sem      chan struct{}  // Semaphore for limiting concurrent worker goroutines.
	jobQueue chan *Job      // Channel for queuing incoming jobs.
	worker   Worker         // Worker interface for processing jobs.
	stopWg   sync.WaitGroup // WaitGroup for tracking the termination of the dispatcher.
}

// / NewDispatcher creates a new instance of a job dispatcher with the given parameters.
func NewDispatcher(worker Worker, maxWorkers int, maxQueueSize int) *Dispatcher {
	return &Dispatcher{
		sem:      make(chan struct{}, maxWorkers), // Buffered channel acting as a semaphore. Use empty struct to minimize the memory allocation
		jobQueue: make(chan *Job, maxQueueSize),
		worker:   worker,
	}
}

// Start initiates the dispatcher to begin processing jobs.
// The dispatcher stops when it receives a value from `ctx.Done`.
func (d *Dispatcher) Start(ctx context.Context) {
	// Increment the wait group counter to indicate that the dispatcher has started processing jobs.
	d.stopWg.Add(1)
	go func() {
		// Initialize a local wait group to track the processing of jobs.
		var wg sync.WaitGroup
	Loop:
		for {
			select {
			case <-ctx.Done():
				// Block until all currently processing jobs have finished.
				wg.Wait()
				break Loop
			case job := <-d.jobQueue:
				// Increment the local wait group to track the processing of this job.
				wg.Add(1)
				// Push to the semaphore to control the number of concurrent workers.
				d.sem <- struct{}{}
				go func(job *Job) {
					defer wg.Done()
					// After the job finishes, pop from the semaphore to release the slot for another job.
					defer func() { <-d.sem }()
					d.worker.Work(job)
				}(job)
			}
		}
		// When the loop exits (due to context cancellation), call the stop method to perform cleanup.
		d.stop()
	}()
}

// Wait blocks until the dispatcher stops.
func (d *Dispatcher) Wait() {
	d.stopWg.Wait()
}

// Enqueue puts a job into the queue.
// If the number of enqueued jobs has already reached the maximum size,
// This will block until space becomes available in the queue to accept a new job.
func (d *Dispatcher) Enqueue(job *Job) {
	d.jobQueue <- job
}

// stop decrements the wait group counter to indicate that the dispatcher has stopped processing jobs
func (d *Dispatcher) stop() {
	d.stopWg.Done()
}
