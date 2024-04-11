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

// TODO: Implement Start method to process jobs from the queue and manage worker routines.
// TODO: Limit the number of concurrent worker goroutines to prevent resource exhaustion.
// TODO: Limit the number of jobs in the Job Queue to avoid excessive memory usage.
// TODO: Stop the job processing when the context is canceled, such as by command 'C'.
func (d *Dispatcher) Start(ctx context.Context) {
	panic("Implement me")
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
