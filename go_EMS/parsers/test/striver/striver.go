package main

import (
	"log"
	"os"
	"flag"
	"fmt"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/striver"
)


func main() {
	in := os.Stdin

	verbosePtr := flag.Bool("v",false,"verbose output")
	flag.Parse()
	args := flag.Args()
	
	if len(args)>0 {
		f, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	striver.Verbose =  *verbosePtr

	parsed, err := striver.ParseReader("",in)
	if err != nil {
		log.Fatal(err)
	}
	decls := striver.ToSlice(parsed)
	spec,err  := striver.ProcessDeclarations(decls)
	if  err!=nil {
		log.Fatal(err)
	}
	fmt.Printf(striver.Sprint(*spec))
}
