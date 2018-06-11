// Note that this version is completely sequential and doesn't use goroutines at all
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "bufio"
    "io"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    et "github.com/elastest/elastest-monitoring-service/go_EMS/eventtag"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
    internalsv "github.com/elastest/elastest-monitoring-service/go_EMS/internalapiserver"
	pe "github.com/elastest/elastest-monitoring-service/go_EMS/eventscounter"
	"github.com/elastest/elastest-monitoring-service/go_EMS/eventout"
	"github.com/elastest/elastest-monitoring-service/go_EMS/eventproc"
)

func main() {
    fmt.Println("Serving server")
    go internalsv.Serve()
    fmt.Println("Server served. Starting scans")

    staticout := os.Args[1]
    dynout := os.Args[2]
    eventout.StartSender(staticout, dynout)

    pipename := "/usr/share/logstash/pipes/leftpipe"
	file, err := os.Open(pipename)
    if err != nil {
        panic(err)
    }
    defer file.Close()
	for {
		scanStdIn(file)
        fmt.Println("RELOADING " + pipename)
	}
	panic("leaving!")
}

func scanStdIn(file io.Reader) {
	scanner := bufio.NewScanner(file)
    var rawEvent map[string]interface{}
    sendchan := eventout.GetSendChannel()
	for scanner.Scan() {
        // Remove this
        i := pe.GetProcessedEvents()
        if i==5 {
            eventproc.UndeploySignals01(444)
        }
        // Up to here
		rawEvent = nil
		thetextbytes := []byte(scanner.Text())
        fmt.Println("Read event ",i)

		if err := json.Unmarshal(thetextbytes, &rawEvent); err != nil {
			fmt.Println("No JSON. Error: " + err.Error())
		} else {
            var evt dt.Event = jsonrw.ReadEvent(rawEvent)
            et.TagEvent(&evt)
            eventproc.ProcessEvent(evt)
            sendchan <- evt
		}
        pe.IncrementProcessedEvents()
	}
}
