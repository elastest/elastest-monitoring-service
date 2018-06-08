package main

import (
    "log"
	"os"
	"flag"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
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

	session.Verbose =  *verbosePtr

	parsed, err := session.ParseReader("",in)
	if err != nil {
		log.Fatal(err)
	}
	decls := common.ToSlice(parsed)
	the_monitor := session.ProcessDeclarations(decls)
	session.Print(the_monitor)
}
