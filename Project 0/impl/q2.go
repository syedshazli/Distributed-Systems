package cs4513_go_impl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	// "strings"
)

/*
Do NOT modify function signature.

Sum numbers from channel `nums` and output sum to `out`.
You should only output to `out` once.
*/
func sumWorker(nums chan int, out chan int) {

	sum := 0
	// This loop runs until the 'nums' channel is closed
	// Both workers will pull from the same 'nums' pipe simultaneously
	for n := range nums {
		sum += n
	}
	out <- sum
}

/*
Do NOT modify function signature.

Read integers from the file `fileName` and return sum of all values.
This function must launch `num` go routines running `sumWorker` to find the sum of the values concurrently.

You should use `checkError` to handle potential errors.
*/
func Sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorker`
	// HINT: use buffered channels for splitting numbers between workers

	// we should probably know how many integers are in the file to find an equal split (found by getting size of resulting array)
	file, err := os.Open(fileName) // open the file given by the string
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // close file when function exits

	result, err := readInts(file)

	if err != nil {
		fmt.Println("ERROR OCCURED: SEE BELOW")
		log.Fatal(err)
	} else {
		fmt.Println(len(result))
	}

	// we now know how many results we have to process and how many goroutines need to be launched (telling us how many elements each goroutine needs)
	numInts := len(result)

	// create output channel for results and buffered input channel to send the results
	out := make(chan int)

	// Buffer size equal to the number of ints
	// ensures the sender never blocks.
	nums := make(chan int, numInts)

	// launch 'num' workers. number of elements they take each is undefined. whoever receives the value first processes and adds to their sum
	for i := 0; i < num; i++ {
		go sumWorker(nums, out)
	}

	for _, v := range result {
		nums <- v
	}
	close(nums)

	collectiveSum := 0
	for i := 0; i < num; i++ {
		collectiveSum += <-out // receive output value from each worker and accumulate for each worker
	}

	return collectiveSum
}

/*
Do NOT modify this function.
Read a list of integers separated by whitespace from `r`.
Return the integers successfully read with no error, or
an empty slice of integers and the error that occurred.
*/
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
