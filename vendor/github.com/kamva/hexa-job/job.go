// Package hjob contains interfaces to
// define workers to process background
// jobs.
//--------------------------------
package hjob

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kamva/hexa"
)

// JobHandlerFunc is the handler of each job in the worker.
type JobHandlerFunc func(context.Context, Payload) error

type Payload interface {
	Decode(payload interface{}) error
}

// Job is a new instance of job to push to the queue by Jobs interface
type Job struct {
	Name  string // required
	Queue string
	// Retry specify retry counts of the job.
	// 0: means that throw job away (and dont push to dead queue) on first fail.
	// -1: means that push job to the dead queue on first fail.
	Retry   int
	Timeout time.Duration
	Payload interface{} // It can be any struct.
}

// Worker is the background jobs worker
type Worker interface {
	// Register handler for new job
	Register(name string, handlerFunc JobHandlerFunc) error
	hexa.Runnable
}

// Jobs pushes jobs to process by worker.
type Jobs interface {
	// Push push job to the default queue
	Push(context.Context, *Job) error
}

// NewJob returns new job instance
func NewJob(name string, payload interface{}) *Job {
	return NewJobWithQueue(name, "default", payload)
}

// NewJobWithQueue returns new job instance
func NewJobWithQueue(name string, queue string, p interface{}) *Job {
	return &Job{
		Name:    name,
		Queue:   queue,
		Retry:   4,
		Timeout: time.Minute * 30,
		Payload: p,
	}
}

//--------------------------------
// Json payload implementation
//--------------------------------

type jsonPayload struct {
	p []byte
}

func (j *jsonPayload) Decode(val interface{}) error {
	return json.Unmarshal(j.p, val)
}

func NewJsonPayload(p []byte) Payload {
	return &jsonPayload{p: p}
}

var _ Payload = &jsonPayload{}
