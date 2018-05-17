// Note that this version is completely sequential and doesn't use goroutines at all
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "bufio"
    "io"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    et "github.com/elastest/elastest-monitoring-service/go_EMS/eventproc"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
    internalsv "github.com/elastest/elastest-monitoring-service/go_EMS/internalapiserver"
	pe "github.com/elastest/elastest-monitoring-service/go_EMS/eventscounter"
	"github.com/elastest/elastest-monitoring-service/go_EMS/moms"
	"github.com/elastest/elastest-monitoring-service/go_EMS/eventout"
)

func main() {
    fmt.Println("Serving server")
    go internalsv.Serve()
    fmt.Println("Server served. Starting scans")

    sendchan := make(chan dt.Event)
    moms.StartEngine(sendchan)

    // Opening staticout
    staticout := os.Args[1]
    // Opening dynout
    dynout := os.Args[2]
    go eventout.StartSender(sendchan, staticout, dynout)

    pipename := "/usr/share/logstash/pipes/leftpipe"
	file, err := os.Open(pipename)
    if err != nil {
        panic(err)
    }
    defer file.Close()
	for {
		scanStdIn(file, sendchan)
        fmt.Println("RELOADING " + pipename)
	}
	panic("leaving!")
}

func scanStdIn(file io.Reader, sendchan chan dt.Event) {
	scanner := bufio.NewScanner(file)
    var rawEvent map[string]interface{}
	for scanner.Scan() {
		rawEvent = nil
		thetextbytes := []byte(scanner.Text())
        fmt.Println("Read event")

		if err := json.Unmarshal(thetextbytes, &rawEvent); err != nil {
			fmt.Println("No JSON. Error: " + err.Error())
		} else {
            var evt dt.Event = jsonrw.ReadEvent(rawEvent)
            et.TagEvent(&evt)
            moms.ProcessEvent(evt)
            sendchan <- evt
		}
        pe.IncrementProcessedEvents()
	}
}
