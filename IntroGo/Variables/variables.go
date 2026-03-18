package variables
import "fmt"

func main() {

	d, e := 1, 2 // Walrus Operator in action

	fmt.Println(d + e)

}

func Strings() (string, string, string) {
	var a string = "Apples" // Define types for both A and B

	var b string = " and Oranges"

	// Combine a and b and set implicitly to c
	var c string = a+b
	return a, b, c
}

func Boolean(first bool, second bool) bool {
	// define and return a boolean OR
	var myAnswer bool =  first || second
	return myAnswer
}