package main

import "log"
import "os"
import . "ems/session"


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
	parsed, err := ParseReader("",in)
	if err != nil {
		log.Fatal(err)
	}
	decls := ToSlice(parsed)
	the_monitor := ProcessDeclarations(decls)
	Print(the_monitor)
}


