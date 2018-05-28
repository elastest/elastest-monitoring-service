package moms

import (
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    strivercp "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "github.com/elastest/elastest-monitoring-service/go_EMS/eventout"
	"time"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
)

func StartEngine(signaldefs []session.MoM) dt.MoMEngine01 {
	writechan := make(chan striverdt.FlowingEvent)
    startWriter(writechan)

    momtostrivervisitor := session.MoMToStriverVisitor{[]dt.Sampler{}, []striverdt.OutStream{}, []striverdt.InStream{}}

    for _,signaldef := range signaldefs {
        signaldef.Accept(&momtostrivervisitor)
    }

    samplers := momtostrivervisitor.Samplers
    kchan := make (chan bool)
    go strivercp.Start(momtostrivervisitor.InStreams, momtostrivervisitor.OutStreams, writechan, kchan)
    return dt.MoMEngine01{samplers, kchan}
}

func startWriter(writechan chan striverdt.FlowingEvent) {
    loc := time.FixedZone("fakeplace", 0)
    sendchan := eventout.GetSendChannel()

    go func () {
        for {
            flowev, open := <-writechan
            if !open {
                break
            }
            theEv := dt.Event{
                sets.SetFromList([]string{string(flowev.Name), "striver"}),
                map[string]interface{}{"value": flowev.Event.Payload.Val},
                time.Unix(int64(flowev.Event.Time),0).In(loc).Format("2006-01-02T15:04:05.000Z"),
            }
            sendchan <- theEv
        }
    }()
}
