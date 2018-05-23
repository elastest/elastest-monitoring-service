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
	"github.com/elastest/elastest-monitoring-service/go_EMS/signals"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

func main() {
    fmt.Println("Serving server")
    go internalsv.Serve()
    fmt.Println("Server served. Starting scans")

    staticout := os.Args[1]
    dynout := os.Args[2]
    eventout.StartSender(staticout, dynout)

    // Remove this

    tagdef := `version 1.0
    when true do #EDS
    when e.path(sender) /\ e.path(sender) do #TJob
    when e.path(sender) /\ e.strcmp(sender,"tjob") do #TJob
    when e.tag(#TJob) do #TORM`
    et.DeployTaggerv01(tagdef)


    defs := []signals.SignalDefinition {
        signals.SampledSignalDefinition{"cpuload", "chan", "system.load.1"},
        signals.SampledSignalDefinition{"hostname", "chan", "beat.hostname"},
        signals.FuncSignalDefinition{"hostnameiselastest", []striverdt.StreamName{"hostname"}, signals.SignalEqualsLiteral{"host_elastest"}},
        signals.ConditionalAvgSignalDefinition{"condavg", "cpuload", "hostnameiselastest"},
        signals.FuncSignalDefinition{"increasing", []striverdt.StreamName{"condavg", "cpuload"}, signals.SignalsLT64{}},
    }
    eventproc.DeployRealSignals01(defs,444)
    // Up to here

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
