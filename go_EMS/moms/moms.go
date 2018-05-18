package moms

import (
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    strivercp "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "github.com/elastest/elastest-monitoring-service/go_EMS/signals"
	"time"
)

var samplers []signals.Sampler

func StartEngine(sendchan chan dt.Event, signaldefs []signals.SignalDefinition) {
	writechan := make(chan striverdt.FlowingEvent)
    startWriter(writechan, sendchan)

    signaltostrivervisitor := signals.SignalToStriverVisitor{[]signals.Sampler{}, []striverdt.OutStream{}, []striverdt.InStream{}}
    for _,signaldef := range signaldefs {
        signaldef.Accept(&signaltostrivervisitor)
    }

    samplers = signaltostrivervisitor.Samplers
    go strivercp.Start(signaltostrivervisitor.InStreams, signaltostrivervisitor.OutStreams, writechan)
}

func ProcessEvent(evt dt.Event) {
    for _,sampler := range samplers {
        sampler.ProcessEvent(evt)
    }
}


func startWriter(writechan chan striverdt.FlowingEvent, sendchan chan dt.Event) {
    loc := time.FixedZone("fakeplace", 0)

    go func () {
        for {
            flowev := <-writechan
            sendchan <- dt.Event{
                sets.SetFromList([]string{string(flowev.Name), "striver"}),
                map[string]interface{}{"value": flowev.Event.Payload.Val},
                time.Unix(int64(flowev.Event.Time),0).In(loc).Format("2006-01-02T15:04:05.000Z"),
            }
        }
    }()
}
