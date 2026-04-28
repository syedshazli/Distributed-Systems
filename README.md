CS4513: The MapReduce Framework
=======================================

Note, this document includes a number of design questions that can help your implementation. We highly recommend that you answer each design question **before** attempting the corresponding implementation.
These questions will help you design and plan your implementation and guide you towards the resources you need.
Finally, if you are unsure how to start the project, we recommend you visit office hours for some guidance on these questions before attempting to implement this project.


Design Questions
------------------

(1 point) 1. If there are n input files, and nReduce number of reduce tasks, how does the MapReduce Framework uniquely name the intermediate files?

The MapReduce Framwork for the nReduceth reduce task are named after the map task as followed: mrtmp.<jobName>-<mapTask>-<0> through mrtmp.<jobName>-<mapTask>-<nReduce-1>. This ensures a consistent naming scheme that is not only consice but also descriptive and uniquely identifiable.


(1 point) 2. Following the previous question, for reduce task r, what are the names of the files it will work on?

For reduce task r, we must first refer to every single map task from partition r that was executed. The files that reduce task r will work on are mrtmp.<jobName>-<0>-r, mrtmp.<jobName>-<1>-r, ... mrtmp.<jobName>-<n-1>-r


(1 point) 3. If the submitted mapreduce job name is "test", what will be the final output file's name?

The final output file's name will be mrtmp.test.


(1 point) 4. Based on `mapreduce/test_test.go`, when you run `TestDistributedBasic()`, how many masters and workers are started? What is the naming convention used for their Unix socket addresses?

When running this test, 2 workers are started and 1 master is started. The naming convention for the master Unix socket address encodes the current user ID and process ID to the path, as well as the suffix (in this case, 'master'). For the workers, the naming convention is 'worker' and the worker ID.


(1 point) 5. In real-world deployments, when giving a mapreduce job, we often start master and workers on different machines (physical or virtual). Describe briefly the protocol that allows the master and workers to become aware of each other's existence and subsequently start working together on completing the mapreduce job. Your description should be grounded in the RPC communications.

We must use a mechanism that tells the master that we exist and are ready to work. This is seen with our register function in worker.go. In register, we then call a function through RPC called Master.Register, which essentially tells the master through RPC that a worker is ready to work for the master. The master then can use the worker for map or reduce based tasks.



(1 point) 6. List all the RPC methods in this implementation and their signatures. Briefly describe the criteria a Go method must satisfy to be callable as an RPC method. (Hint: https://golang.org/pkg/net/rpc/)

There are multiple RPC methods used by the master and worker used to communicate with each other. 

Worker.DoTask is called by the master when it wants to assign a worker to do a task. The master must pass in is  JobName (string), File (string), Phase (jobPhase), TaskNumber (int) 

Worker.Shutdown is called by the master once all work is completed, and the number of tasks completed is the output. 

Master.Register is called by the worker to tell the master the worker exists and is ready to work, taking the master as a string and passing in the worker name as well.

Master.Shutdown is called by a master to shut down its own registration server, the signature is func (mr *Master) Shutdown(_, _ *struct{}) error

In order to be callable as an RPC method, a Go method must follow the signature of func (t *T) MethodName(argType T1, replyType *T2) error. To specify, the method must be exported by starting with a capital letter. The type must also be exported. The method must have 2 arguments, where the second argument is a pointer where the RPC server can write a result back to. The method should also return an error.


(1 point) 7. In your Inverted Index implementation (Part C), a word may appear multiple times in the same file, so reduceF may receive the same filename more than once in the `values` slice. How does your implementation handle this? What does the output line look like for a word that appears in exactly two distinct files?

First, the list of filenames in a file is sorted alphabetically. Then, if the current filename is the same as the filename before it, then we do not append the current filename, since the previous filename was already appended to the slice. For a word that appears in exactly two distinc files, the output would be 'word: file1, file2'

(1 point) 8. In schedule() (Part D/E), after a task completes or fails, the worker's address is returned to registerChannel so it can be reused. Why must this happen regardless of whether the task succeeded or failed? What would go wrong if the address were only returned on success?

If the address were only returned upon success, then other tasks could not use the worker after the task failed, rendering that worker as useless since it cannot pick up any more work.

(1 point) 9. In your Part F straggler implementation, both an original goroutine and a backup goroutine may successfully complete the same task. How does your implementation ensure the task is counted as finished exactly once? Describe the synchronization mechanism you used and explain why it is correct.

We know a task is not finished by checking if their status is set to 0 or 1. If 0, then the task is not finished, and a backup is launched if ntasks-1 are remaining. We find the status by using a atomic.LoadInt32 primitive, which ensures the correct value is read. 
wg.Done() indicates the task is finished, however, if we were to let both the straggler and worker call wg.Done() freely, we would have counted the task as finished twice instead of once. As a result, we utilize the compare and exchange operator in Go that ensures clean atomic variable writing. 
atomic.CompareandSwap() will take the status of the current task, and check if it's equal to 0. If so, we set the value of the status to 1, and return true. We check if the value returned true, and if so, we call wg.Done(). If the value did not return true, we do not call wg.Done(), indicating the task was already finished and is counted once.



Errata
------

Describe any known errors, bugs, or deviations from the requirements.

---
