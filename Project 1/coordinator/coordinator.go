// Package coordinator implements the RPC coordinator server.
package coordinator

// Add what you need here.
import (
	"cs4513/project1/types"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

// jobRecord is the coordinator's private copy of a job and its status.
type jobRecord struct {
	jobID types.JobID
	spec  types.JobSpec
	state types.JobState
}

// Coordinator is the RPC server.
// Keep shared state here and protect it with a mutex.
type Coordinator struct {
	mu           sync.Mutex
	jobs         map[types.JobID]*jobRecord   // all jobs , keyed by ID
	queue        []types.JobID                // pending job IDs in FIFO order
	workers      map[types.WorkerID]time.Time // registered workers
	nextJobID    int
	nextWorkerID int64
}

// New returns an initialized Coordinator.
// we initialize an empty map of jobs and their records, workers and their IDs, and a queue of job IDs to ex
func New() *Coordinator {
	return &Coordinator{
		jobs:    make(map[types.JobID]*jobRecord),
		workers: make(map[types.WorkerID]time.Time),
		queue:   []types.JobID{},
	}

}

// Start creates a Coordinator, registers it, and starts listening on addr.
func Start(addr string) error {
	myCoord := New() // create the coordinator

	srv := rpc.NewServer()
	srv.Register(myCoord) // register the coordinator

	ln, err := net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("Error listening on address %w: ", err)
	}

	for { // forever loop that accepts connections and sends out go routines
		conn, err := ln.Accept()
		// handle err ...
		if err != nil {
			return fmt.Errorf("Error connecting to server: %w", err)
		}
		go srv.ServeConn(conn)
	}

}

// SubmitJob adds a new job and returns its ID.
func (c *Coordinator) SubmitJob(spec types.JobSpec, reply *types.JobID) error {
	// TODO: implement
	c.mu.Lock()

	// new job ID is the nextJobID
	ID := types.JobID("job-" + strconv.Itoa(c.nextJobID))

	// create a new jobRecord with the new ID and given spec
	jobRec := jobRecord{ID, spec, types.StatePending}

	// increment the nextJobID
	c.nextJobID++

	// add the new ID to new jobs
	c.jobs[ID] = &jobRec

	c.queue = append(c.queue, ID)

	c.mu.Unlock()

	*reply = ID

	return nil
}

// QueryJob returns the current status of a job.
func (c *Coordinator) QueryJob(id types.JobID, reply *types.JobStatus) error {
	// TODO: implement
	return nil
}

// ListJobs returns a summary of every known job.
func (c *Coordinator) ListJobs(_ struct{}, reply *[]types.JobSummary) error {
	// TODO: implement
	return nil
}

// Register assigns a WorkerID to a new worker.
func (c *Coordinator) Register(_ types.WorkerInfo, reply *types.WorkerID) error {
	// TODO: implement
	return nil
}

// RequestJob hands out the next pending job.
func (c *Coordinator) RequestJob(workerID types.WorkerID, reply *types.Job) error {
	// TODO: implement
	return nil
}

// ReportResult stores the result of a finished job.
func (c *Coordinator) ReportResult(result types.JobResult, _ *struct{}) error {
	// TODO: implement
	return nil
}
