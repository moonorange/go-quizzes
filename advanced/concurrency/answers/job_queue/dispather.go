package main

import (
	"context"
	"sync"
)

// Job represents a interface of job that can be enqueued into a dispatcher.
type Job struct {
	Payload string
}

type Worker interface {
	Work(j *Job)
}

// Dispatcher represents a job dispatcher.
type Dispatcher struct {
	sem       chan struct{} // semaphore
	jobBuffer chan *Job
	worker    Worker
	wg        sync.WaitGroup
}

// NewDispatcher will create a new instance of job dispatcher.
// maxWorkers means the maximum number of goroutines that can work concurrently.
// buffers means the maximum size of the queue.
func NewDispatcher(worker Worker, maxWorkers int, buffers int) *Dispatcher {
	return &Dispatcher{
		// Restrict the number of goroutine using buffered channel (as counting semaphor)
		sem:       make(chan struct{}, maxWorkers),
		jobBuffer: make(chan *Job, buffers),
		worker:    worker,
	}
}

// Start starts a dispatcher.
// This dispatcher will stops when it receive a value from `ctx.Done`.
func (d *Dispatcher) Start(ctx context.Context) {
	d.wg.Add(1)
	go d.loop(ctx)
}

// Wait blocks until the dispatcher stops.
func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

// Enqueue puts a job into the queue.
// If the number of enqueued jobs has already reached the maximum size,
// this will block until space becomes available in the queue to accept a new job.
func (d *Dispatcher) Enqueue(job *Job) {
	d.jobBuffer <- job
}

func (d *Dispatcher) stop() {
	d.wg.Done()
}

func (d *Dispatcher) loop(ctx context.Context) {
	var wg sync.WaitGroup
Loop:
	for {
		select {
		case <-ctx.Done():
			// block until all the jobs finishes
			wg.Wait()
			break Loop
		case job := <-d.jobBuffer:
			// Increment the waitgroup
			wg.Add(1)
			// Decrement a semaphore count
			d.sem <- struct{}{}
			go func(job *Job) {
				defer wg.Done()
				// After the job finished, increment a semaphore count
				defer func() { <-d.sem }()
				d.worker.Work(job)
			}(job)
		}
	}
	d.stop()
}
