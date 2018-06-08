package main

import ( "log"
	"os"
	"flag"
	"fmt"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/stamp"
//    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
)

var Verbose bool

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

	stamp.Verbose = *verbosePtr

	if Verbose { fmt.Printf("verbosity on\n") }
	parsed, err := stamp.ParseReader("",in)
	if err != nil {
		log.Fatal(err)
	}
	the_monitor := stamp.Monitor{parsed.(stamp.Filters).Defs}
	stamp.Print(the_monitor)
}


