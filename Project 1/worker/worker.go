// Package worker implements the worker process.
//
// A worker connects over RPC, registers once, then loops:
// request a job, run it locally, report the result.
package worker

import (
	"fmt"
	// "cs4513/project1/coordinator"
	"cs4513/project1/types"
	"cs4513/project1/jobs"

	"time"
	"net/rpc"
)

// Run connects to the coordinator and processes jobs until it hits an error.
func Run(addr string) error {

	client,err := rpc.Dial("tcp", addr)

	if err != nil {

	return fmt.Errorf("dial %s: %w", addr , err )
	}
	defer client.Close()

	var WorkerID types.WorkerID
	var t types.WorkerInfo
	rpcErr := client.Call("Coordinator.Register", t, &WorkerID)
	if rpcErr != nil {
		return fmt.Errorf("Error registering worker %s", rpcErr)
	
	}

	var jobStats types.Job
	for{
		Err := client.Call("Coordinator.RequestJob", WorkerID, &jobStats )
		// fmt.Println("Job received!")
		if types.IsNoWork(Err){
			time.Sleep(1 * time.Second)
			continue
		}
		if Err != nil {
		return fmt.Errorf("Error requesting job %s:", Err)
		}


		// execute job locally
		var resultingJob types.JobResult

		resultingJob = jobs.ExecuteJob(jobStats)

		resultingJob.WorkerID = WorkerID

		ReportErr := client.Call("Coordinator.ReportResult", resultingJob, &struct{}{})
		if ReportErr != nil {

		return fmt.Errorf("Encountered an error reporting the result %s:", ReportErr)
		}
		
	}

}
