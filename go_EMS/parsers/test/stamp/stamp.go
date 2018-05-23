package main

import "log"
import "os"
import . "ems/stamp"


func main() {
	in := os.Stdin
	if len(os.Args) >1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}
	parsed, err := ParseReader("foo",in)
	if err != nil {
		log.Fatal(err)
	}
	the_monitor := Monitor{parsed.(Filters).Defs}
	Print(the_monitor)
}


