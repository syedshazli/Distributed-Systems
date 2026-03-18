package main

import (
	"fmt"
	"example.com/greetings"
)

func main(){
	Result, err := greetings.Hello("Syed")
	if err != nil{
		fmt.Println("BAD BAD ERROR")

	}
	fmt.Println(Result)

}