package greetings

import (
	"fmt"
	"errors"
)

// function Hello returns a greeting for the named person. Takes in a string called name and outputs a string
// func starts with capitla letter: can be used outside of pkg
func Hello(name string) (string, error) {
	 // If no name was given, return an error with a message.
    if name == "" {
        return name, errors.New("empty name")
    }
    // Return a greeting that embeds the name in a message.
    message := fmt.Sprintf("Hi, %v. Welcome!", name) // declare and initialize var in one line
    return message, nil
}