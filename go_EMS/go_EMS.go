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

    // Remove this

    /* This is broken now
    tagdef := `version 1.0
    when true do #EDS
    when e.tag(#TJob) do #TORM`
    et.DeployTaggerv01(tagdef)
    */


    defs := `
    pred istjobmark := e.path(TJobMark)
    pred isnet := e.strcmp(system.network.name,"eth0")
    stream bool truestream := true
    stream num inbytes := if isnet then e.getnum(system.network.in.bytes)

    stream bool low_is_running := if istjobmark then e.strcmp(TJobMark, "LOW_START")
    stream num gradlow := gradient(inbytes within low_is_running)
    stream num avggradlow := avg(gradlow within truestream)

    stream bool high_is_running := if istjobmark then e.strcmp(TJobMark, "HIGH_START")
    stream num gradhigh := gradient(inbytes within high_is_running)
    stream num avggradhigh := avg(gradhigh within truestream)

    stream bool testcorrect := avggradhigh * 0.7 < avggradlow

    trigger isnet do emit inbytes on #bytesval
    trigger isnet do emit avggradlow on #bytesgradlow
    trigger isnet do emit avggradhigh on #bytesgradhigh
    trigger isnet do emit testcorrect on #testresult
    `
    /*stream num load := if otrohost then e.getnum(system.load.1)
    stream bool high_load := load > 0.4
    stream num avgcond := avg(load within pred)
    trigger e.strcmp(beat.hostname,"otrohost") do emit load on #outchannel
    trigger true do emit high_load on #outhighload`*/
    eventproc.DeploySignals01(defs)
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
