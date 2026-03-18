package main
import(
	"fmt"
	"sync"
	"time"
)


func printStr( c chan string) {
	s := <- c 
	fmt.Println(s);
	
}

func main() {

	// task a goroutine to print out a string
	go fmt.Println("This is a String!")

	// task 5 goroutines to print strings fed into a channel by the main thread

	myChannel := make(chan string)
	go printStr( myChannel)
	go printStr( myChannel)
	go printStr( myChannel)
	go printStr( myChannel)
	go printStr(myChannel)
	myChannel <- "Hello"
	myChannel <- "My Name Is"
	myChannel <- "Syed"
	myChannel <- "I hope youre doing well"
	myChannel <- "Bye!"





	// using a waitgroup, halt the main thread for a set amount of time using 3 different goroutines
	var wg sync.WaitGroup
	wg.Add(3) // 3 routines need to decrement the counter/semaphore
	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Sleeping 1")
		wg.Done()

	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Sleeping 2")
		wg.Done()

	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Sleeping 3")
		wg.Done()

	}()

	wg.Wait() // reaches here and waits for wg to be decremented to 0

	return
}

func Buffered_channel(limit int) chan string {
	// create and return a buffered channel for a passed amount of strings
	newChannel := make(chan string, limit)
	return newChannel
}