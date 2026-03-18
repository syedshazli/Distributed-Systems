package conditionals

import "fmt"

func main() {

	var a bool = true
	var str string = "Print if true"

	if a {
		fmt.Println(str)
	}

	var b bool = false
	if !b{
		fmt.Println("Not true")
	}
	
}

func Switch_statement(digit int) string {
	// create a switch statement that deals with all possible digits (0-9) and returns character representation
	var myString = ""
	switch digit {
    case 0:
            myString = "0"
    case 1:
            myString = "1"
    case 2:
            myString = "2"
	case 3:
            myString = "3"
	case 4:
		    myString = "4"
	case 5:
		    myString = "5"
	case 6:
			myString = "6"
	case 7:
		    myString = "7"
	case 8:
		    myString = "8"
	case 9:
		    myString = "9"


}
	return myString;
}