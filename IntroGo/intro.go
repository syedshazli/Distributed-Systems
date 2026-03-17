package main // belongs to executable package
import "rsc.io/quote" // run 'go mod tidy'

// build with go build

// bring formatting package
import(
	"fmt"
)
// main function just like C
func main(){
	fmt.Println("Hello World! :) ")
	fmt.Println(quote.Go())

}

// run with go run .