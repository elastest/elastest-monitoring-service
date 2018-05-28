package main

import (
    "log"
    "os"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
)


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
	parsed, err := session.ParseReader("",in)
	if err != nil {
		log.Fatal(err)
	}
	decls := common.ToSlice(parsed)
	the_monitor := session.ProcessDeclarations(decls)
	session.Print(the_monitor)
}
