package main

import (
	"testing"
	"os"
)

func TestRest(t *testing.T) {
	//main()
	file, err := os.Open("testinputs/testdefs.json")
    if err != nil {
        panic(err)
    }
	file, err = os.Open("testinputs/testEvents.txt")
    if err != nil {
        panic(err)
    }
    os.Args = []string{"goEMS", "/dev/null", "/dev/null"}
	scanStdIn(file)
}
