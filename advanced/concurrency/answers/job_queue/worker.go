package main

import (
	"fmt"
)

// APIWorker sends post requests
type APIWorker struct{}

func NewAPIWorker() *APIWorker {
	return &APIWorker{}
}

func (p *APIWorker) Work(job *Job) {
	fmt.Println(job.Payload)
}
