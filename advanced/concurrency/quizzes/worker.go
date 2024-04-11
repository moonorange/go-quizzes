package main

import "fmt"

// WorkerImpl sends post requests
type WorkerImpl struct{}

func NewWorker() *WorkerImpl {
	return &WorkerImpl{}
}

func (p *WorkerImpl) Work(job *Job) {
	fmt.Println(job.Payload)
}
