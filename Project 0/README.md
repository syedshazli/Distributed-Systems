CS4513: Warmup The Go Basic
==========================

Student
-----------------

[Syed Shazli] ([sashazli@wpi.edu])

Two Problems
------------------

1. What is the difference between unbuffered and buffered channels? And why do you choose one over the other for this assignment?

Buffered channels explicitly state how many values it expects to be sent into the channel, while unbuffered channels do not have such primitive. Buffered channels enables non-blocking parallelism, where all values are in the channel before workers process the values, allowing for parallel work, which could not be done in an unbuffered implementation. In the 2nd problem, I decided to use a buffered channel. This is because in a unbuffered channel, the main goroutine can send a worker, but would block until the worker receives the message, meaning another working waiting to get work could not run in this time. Using one buffered channel which is the length of the input array means workers don't have to wait for synchronous handshakes by the sender, and can process immedietley. If I used a unbuffered channel, I would have to send over each element and wait for it to be received before processing, which will not scale well when we have to deal with tens of thousands of numbers.

2. Briefly explain how you approached the two problems.

In the 1st problem, I broke down the top-K words approach into several steps. The 1st step was to get the whole file in one string, and make every character lowercase with strings.ToLower. I then made an array containing each token in the string, separated by whitespaces. This was done before removing any punctuation, because I first tried to remove punctuation all at once, and ended up overcounting words such as --Jonathan. 
I defined the regex and cleaned each word, removing any non alphanumeric characters from each token. If the new length of that token was greater than or equal to the charaacter threshold, I added it to a new array of filtered tokens that met the criteria. The next step was to make a hashmap of each word and their count, which was done by iterating through the list of cleaned words and incrementing the count by 1 if the string was in the hashmap, otherwise initializing with 1. 
After this, I then converted the hashmap to an array of WordCount structs, and sorted the struct so that the count is in descending order. I then sliced the array to only return the top numWords elements, succesfully passing all test cases.

In the 2nd problem, I broke it down to a couple high level steps
    i: Read the data from the file into an array of integers
    ii: Launch 'numWorkers' workers to process the array at the same time, with each worker having a local sum
    iii: Accumulate all the local sums of the workers to get the total sum of the file.

In terms of the details, I just used os.Open to get the file contents locally, calling the given readInts() function which takes in a file and gives back a array of integers in the file. 

I created one output channel for results, allowing multiple workers to send values as output which would then be processed as an accumulating sum. I then made a buffered channel which had a capacity of the number of elements in the array. I chose a buffered channel over an unbuffered channel to prevent constant communication overhead required by unbuffered channels, which would cause high runtimes when having thousands of elements to process. 

My approach does not give each worker a set amount of elements to process. Instead, I used a first-come first-serve pattern, also known as work-stealing. Whenever workers are ready, they take an element to add to their local sum from the shared buffered channel. The approach is first-come first-serve since the first worker available to process an element will grab it.

Finally, once there are no more elements left to process, signalled by the sender closing the channel, the workers will take their local sum and return it, which will be collected by a for loop which accumulates the sum of each local sum of the worker.

Measurement
------------------

Paste your measurement output here.
Timing cs4513_go_test/q2_test1.txt...
Average time for 1 worker in test1.txt was:  328.205µs  with a result of  499500
Average time for 2 worker in test1.txt was:  490.109µs  with a result of  499500
Average time for 3 worker in test1.txt was:  382.716µs  with a result of  499500
Average time for 5 worker in test1.txt was:  471.094µs  with a result of  499500
Average time for 10 worker in test1.txt was:  431.713µs  with a result of  499500
Average time for 100 worker in test1.txt was:  947.022µs  with a result of  499500
--------------------TIMES FOR TEST2.TXT-----------------
Average time for 1 worker in test2.txt was:  331.254µs  with a result of  117652
Average time for 2 worker in test2.txt was:  487.64µs  with a result of  117652
Average time for 3 worker in test2.txt was:  564.948µs  with a result of  117652
Average time for 5 worker in test2.txt was:  639.331µs  with a result of  117652
Average time for 10 worker in test2.txt was:  570.761µs  with a result of  117652
Average time for 100 worker in test2.txt was:  1.154692ms  with a result of  117652
--------------------TIMES FOR TEST3.TXT-----------------
Average time for 1 worker in test3.txt was:  858.555µs  with a result of  617152
Average time for 2 worker in test3.txt was:  915.213µs  with a result of  617152
Average time for 3 worker in test3.txt was:  928.986µs  with a result of  617152
Average time for 5 worker in test3.txt was:  1.017386ms  with a result of  617152
Average time for 10 worker in test3.txt was:  1.140654ms  with a result of  617152
Average time for 100 worker in test3.txt was:  2.024248ms  with a result of  617152
--------------------TIMES FOR TEST4.TXT-----------------
Average time for 1 worker in test4.txt was:  2.606758ms  with a result of  4995000
Average time for 2 worker in test4.txt was:  2.760523ms  with a result of  4995000
Average time for 3 worker in test4.txt was:  3.609355ms  with a result of  4995000
Average time for 5 worker in test4.txt was:  3.782855ms  with a result of  4995000
Average time for 10 worker in test4.txt was:  4.243487ms  with a result of  4995000
Average time for 100 worker in test4.txt was:  7.599749ms  with a result of  4995000
--------------------------TIMES FOR TEST5.TXT------------------------
Average time for 1 worker in test5.txt was:  23.505559ms  with a result of  49950000
Average time for 2 worker in test5.txt was:  29.006334ms  with a result of  49950000
Average time for 3 worker in test5.txt was:  32.472809ms  with a result of  49950000
Average time for 5 worker in test5.txt was:  35.689944ms  with a result of  49950000
Average time for 10 worker in test5.txt was:  41.898264ms  with a result of  49950000
Average time for 100 worker in test5.txt was:  84.692522ms  with a result of  49950000

Observations and Explanations
------------------

Briefly describe the relationship among the result, number of workers, and total time. Include possible explanations. You may also include additional measurement scenarios here.

In this work-stealing partition, the result always returns the correct answer. However, each text file has a different number of elements, which is important to account for scale. 

In the 1st text file, q2_test1.txt, the runtime increases as we add more workers, with an approximate 3x runtime increase adding 100x more workers. This is because there are only 1000 elements, which can be done relativley fast by 1 worker. Adding more workers requires more synchronization overhead, and as a result the parallelism gains from adding more workers are not very apparent.

For the 2nd text file, q2_test2.txt, the runtime also increases as we add more workers, even if we have the same amount of elements as test1.txt! The reason for this may be because test1.txt has a very predictable pattern, where each element in the array is the index of the array. As a result, the compiler may have done some very advanced optimizations due to it being very easy to predict the next element in test1.txt, while the numbers in test2.txt are more random, so compiler prediction based optimizers are of little help. This can explain why performs noticabley degrades compared to test1 runtime even with the same amount of workers.

The 3rd text file, q2_test3.txt, has twice as many elements (2000) but only takes approximately less than 2x the time compared to test2.txt. The results remain correct, but the average time is still high, presumabley because 2000 elements is still not that many elements for one worker to process, and having to acquire a lock between workers to process elements still takes some time. Interestingly, the runtimes decrease going from 1 to 2 workers, and only increase marginally (70 us) when scaling to 3 workers, while in test2 the increase was hundreds of microseconds.

The 4th text file, q2_test4.txt, has 10,000 elements to process. Now, seeing 1 worker doing the work has gone into the milliseconds, not seen before with 1000-2000 elements to process. Scaling up to more workers, the performance gains are still not necessarily a perfect equation. For example, one would expect by now, with 10,000 elements to process, that for 5x more workers, we should get a 5x performance improvement compared to 1 worker. However, the performance actually drops off, increasing by less than 2x but increasing nonetheless rather than decreasing, which is what we expected. Scaling up to more workers at this stage generally leads to higher runtimes. Again, this is likely due to channel synchronization overhead, where each of the workers compete to receive from the channel, which means a lock needs to be placed once a worker gains access to the channel to decrement. When reaching many goroutines, eventually the CPU will have to do rapid context switching instead of true parallelism, resulting in frequent cache misses and loads from memory.

Finally, the 5th text file, q2_test5.txt, with 100,000 elements to process, the runtimes are still higher compared to before, as expected. The reasoning would be the same as before, where the communication and locking overhead would cause higher runtimes compared to 1 worker. Additionally, cache coherence protocols would require frequent invalidating of cache lines since each core that's taking a worker needs to have a proper copy of the buffer. As a result, in this example, runtime still increases as the number of workers increases, but the runtime increase is usually always less than 2x the runtime with 1 worker.

For future improvements, it would be interesting to have each worker take a certain amount of elements, which wouldn't require as much synchronization between each worker, potentially reducing runtime. It becomes apparent that this communication becomes a bottleneck for runtime as the worker count increases.

Errata
------

List any known errors, bugs, or deviations from the requirements.

Although this is not a deviation from the requirement, each worker does not take an equal amount of elements. Instead, a first-come first-served functionality is employed, where whichever worker is available will grab the data and process it.
