// client is a small CLI for submit/query/list.
//
//	go run ./cmd/client -addr localhost:8080 submit <type> <payload-json>
//	go run ./cmd/client -addr localhost:8080 query <job-id>
//	go run ./cmd/client -addr localhost:8080 list
package main

import (
	"cs4513/project1/types"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "coordinator address (host:port)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: client -addr <host:port> <submit|query|list> [args...]")
		os.Exit(1)
	}

	client, err := rpc.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("client: connect to %s: %v", *addr, err)
	}
	defer client.Close()

	switch args[0] {
	case "submit":
		cmdSubmit(client, args[1:])
	case "query":
		cmdQuery(client, args[1:])
	case "list":
		cmdList(client)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n", args[0])
		os.Exit(1)
	}
}

// cmdSubmit submits a job and prints the returned ID.
func cmdSubmit(client *rpc.Client, args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: client submit <type> <payload-json>")
		os.Exit(1)
	}

	// Note: The job is processed and sent out  by the coordinator
	// the JobID is found by querying 'QueryJobs'

	var returnedID types.JobID
	var jobSpec types.JobSpec
	jobSpec.Type = args[0]
	jobSpec.Payload = []byte(args[1])
	err := client.Call("Coordinator.SubmitJob", jobSpec, &returnedID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "RPC call failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Submitted " + returnedID)

}

// cmdQuery prints the current state of one job.
func cmdQuery(client *rpc.Client, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: client query <job-id>")
		os.Exit(1)
	}
	// TODO: implement

	// need to call queryJob
	var currID types.JobID
	var queryResult types.JobStatus

	currID = types.JobID(args[0])
	err := client.Call("Coordinator.QueryJob", currID, &queryResult)

	if err != nil{
		fmt.Println(err)
	}
	fmt.Print(queryResult.ID," ", queryResult.State)

	if queryResult.WorkerID != ""{
		fmt.Print(" ("+queryResult.WorkerID+")")
	}

	if queryResult.State == 2{ // 2 corresponds to done in 0 index
		fmt.Println(string(queryResult.Output))
	} else if queryResult.State == 3{
		fmt.Println(queryResult.Err)
	}

}

// cmdList prints all jobs known to the coordinator.
func cmdList(client *rpc.Client) {
	// TODO: implement
}
