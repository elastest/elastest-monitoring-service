package moms

import (
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    strivercp "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
    "github.com/elastest/elastest-monitoring-service/go_EMS/eventout"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
    parserimpl "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/impl"
    "strings"
)

func StartEngine(signaldefs []session.MoM) dt.MoMEngine01 {
	writechan := make(chan striverdt.FlowingEvent)
    startWriter(writechan)

    inSignalName := striverdt.StreamName("elastest::in_events")
	signalchan := make(chan striverdt.Event)
    sampler := dt.Sampler{signalchan}
	inStream := striverdt.InStream{inSignalName, &striverdt.InFromChannel{signalchan, nil, 0, false}}
    inStreams := []striverdt.InStream{inStream}

    momtostrivervisitor := parserimpl.MoMToStriverVisitor{[]striverdt.OutStream{}, inSignalName}

    for _,signaldef := range signaldefs {
        signaldef.Accept(&momtostrivervisitor)
    }

    kchan := make (chan bool)
    go strivercp.Start(inStreams, momtostrivervisitor.OutStreams, writechan, kchan)
    return dt.MoMEngine01{sampler, kchan}
}

func startWriter(writechan chan striverdt.FlowingEvent) {
    sendchan := eventout.GetSendChannel()

    go func () {
        for {
            flowev, open := <-writechan
            if !open {
                break
            }
            if strings.HasPrefix(string(flowev.Name), string(parserimpl.TRIGGER_PREFIX)) {
                theEv := flowev.Event.Payload.Val.(dt.Event)
                sendchan <- theEv
            }
        }
    }()
}
