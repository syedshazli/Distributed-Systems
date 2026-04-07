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
	jobID  types.JobID
	spec   types.JobSpec
	status types.JobStatus
}

// Coordinator is the RPC server.
// Keep shared state here and protect it with a mutex.
type Coordinator struct {
	mu           sync.Mutex
	jobs         map[types.JobID]*jobRecord   // all jobs , keyed by ID
	queue        []types.JobID                // pending job IDs in FIFO order
	workers      map[types.WorkerID]time.Time // registered workers
	nextJobID    int
	nextWorkerID int
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
	// we lock the region to add an new job so that there are not any jobs with the same ID
	// NOTE: We may be able to just retrieve the current ID and increment it and leave the critical region
	c.mu.Lock()
	defer c.mu.Unlock()


	// new job ID is the nextJobID
	ID := types.JobID("job-" + strconv.Itoa(c.nextJobID))

	// create a new jobRecord with the new ID and given spec
	jobRec := jobRecord{ID, spec, types.JobStatus{ID, types.StatePending, "", nil, ""}}

	// increment the nextJobID
	c.nextJobID++

	// add the new ID to new jobs
	c.jobs[ID] = &jobRec

	c.queue = append(c.queue, ID)

	*reply = ID

	return nil
}

// QueryJob returns the current status of a job.
func (c *Coordinator) QueryJob(id types.JobID, reply *types.JobStatus) error {
	c.mu.Lock()
	defer c.mu.Unlock()


	record, ok := c.jobs[id]
	if !ok {
		c.mu.Unlock()
		return fmt.Errorf("job not found: %s", id)
	}

	status := record.status

	*reply = status
	return nil
}

// ListJobs returns a summary of every known job.
func (c *Coordinator) ListJobs(_ struct{}, reply *[]types.JobSummary) error {
	// must place a lock here since the number of elements can change without a lock
	c.mu.Lock()
	defer c.mu.Unlock()

	numJobs := len(c.jobs)
	arr := make([]types.JobSummary, numJobs)
	idx:= 0
	for key := range c.jobs{
		arr[idx] = types.JobSummary{ID: key, Type: c.jobs[key].spec.Type, State: c.jobs[key].status.State}
		idx++
	}
	fmt.Println("Num elements = ", len(arr))
	*reply = arr
	return nil
}

// Register assigns a WorkerID to a new worker.
func (c *Coordinator) Register(_ types.WorkerInfo, reply *types.WorkerID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ID := types.WorkerID("worker-" + strconv.Itoa(c.nextWorkerID))
	c.nextWorkerID++
	c.workers[ID] = time.Time{}
	*reply = ID
	return nil
}

// RequestJob hands out the next pending job.
func (c *Coordinator) RequestJob(workerID types.WorkerID, reply *types.Job) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.queue) == 0{
		err := types.ErrNoWork
		if types.IsNoWork(err){
			return err
		}

	}
	pendingJob := c.queue[0]
	c.queue = c.queue[1:]
	ID := c.jobs[pendingJob] // returns a jobRecord
	ID.status.State = types.StateRunning 
	ID.status.WorkerID = workerID

	myJob := types.Job{ID.jobID, ID.spec}
	*reply = myJob
	return nil
}

// ReportResult stores the result of a finished job.
func (c *Coordinator) ReportResult(result types.JobResult, _ *struct{}) error {
	// TODO: implement
	c.mu.Lock()
	defer c.mu.Unlock()

	jobID := result.JobID
	if c.jobs[jobID] == nil{
		return fmt.Errorf("Could not find job")
	}

	if result.Err == ""{
		jobRecord := c.jobs[jobID] // is the corresponding jobRecord
		jobRecord.status.State = types.StateDone
		jobRecord.status.Output = result.Output
	}else{
		c.jobs[jobID].status.State = types.StateFailed
		c.jobs[jobID].status.Err = result.Err
	}


	return nil
}
