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

    // Ad hoc

	/*signalchan := make(chan striverdt.Event)
	samplers = []signals.Sampler{signals.Sampler{"chan", "system.load.1", signalchan}}
	inStream := striverdt.InStream{"cpuload", &striverdt.InFromChannel{signalchan, nil, 0}}

    plusone := func (args...striverdt.EvPayload) striverdt.EvPayload{
        myprev := args[0]
        prev := 0
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(int)
        }
        return striverdt.Some(prev+1)
    }
    plusOneVal := striverdt.FuncNode{[]striverdt.ValNode{&striverdt.PrevValNode{striverdt.TNode{}, "cpuevcount", []striverdt.Event{}}}, plusone}

    cpuloadcounter := striverdt.OutStream{"cpuevcount", striverdt.SrcTickerNode{"cpuload"}, plusOneVal}

    acc := func (args...striverdt.EvPayload) striverdt.EvPayload{
        myprev := args[0]
        accval := args[1]
        prev := 0.0
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(float64)
        }
        return striverdt.Some(prev+accval.Val.(striverdt.EvPayload).Val.(float64))
    }
    accVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevValNode{striverdt.TNode{}, "cpuacc", []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, "cpuload", []striverdt.Event{}},
    }, acc}

    cpuacc := striverdt.OutStream{"cpuacc", striverdt.SrcTickerNode{"cpuload"}, accVal}

    // endof ad hoc
	go strivercp.Start([]striverdt.InStream{inStream}, []striverdt.OutStream{cpuloadcounter, cpuacc}, writechan)*/
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
                sets.SetFromList([]string{string(flowev.Name)}),
                map[string]interface{}{"value": flowev.Event.Payload.Val},
                time.Unix(int64(flowev.Event.Time),0).In(loc).Format("2006-01-02T15:04:05.000Z"),
            }
        }
    }()
}
