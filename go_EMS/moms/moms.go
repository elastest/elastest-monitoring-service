package moms

import (
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    defs "github.com/elastest/elastest-monitoring-service/go_EMS/signals"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    strivercp "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "strings"
	"time"
)

type Sampler struct {
    Def defs.SampledSignalDefinition
    OutChan chan striverdt.Event
}

func (s Sampler) processEvent(evt dt.Event) {
    payload := striverdt.NothingPayload
    if sets.SetIn(s.Def.InChannel, evt.Channels) {
        payload = striverdt.Some(extractFromMap(evt.Payload, s.Def.ValuePath))
    }
    t, err := time.Parse(time.RFC3339Nano,evt.Timestamp)
    if err != nil {
        panic(err)
    }
    striverEvent := striverdt.Event{striverdt.Time(t.Unix()), payload}
    s.OutChan <- striverEvent
}

var samplers []Sampler

func StartEngine(sendchan chan dt.Event) {
	signalchan := make(chan striverdt.Event)
	writechan := make(chan striverdt.FlowingEvent)
	samplers = []Sampler{Sampler{defs.SampledSignalDefinition{"signalname", "chan", "system.load.1"}, signalchan}}
	theSampler := striverdt.InStream{"signalname", &striverdt.InFromChannel{signalchan, nil}}

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
	go strivercp.Start([]striverdt.InStream{theSampler}, []striverdt.OutStream{}, writechan)
}

func ProcessEvent(evt dt.Event) {
    for _,sampler := range samplers {
        sampler.processEvent(evt)
    }
}

func extractFromMap(themap map[string]interface{}, strpath dt.JSONPath) interface{} {
	path := strings.Split(string(strpath), ".")
	var ok bool
	if (len(path) == 0) {
		panic("empty path")
	}
	for _,key := range path[:len(path)-1] {
		themap, ok = themap[key].(map[string]interface{})
		if (!ok) {
			panic("incorrect path")
		}
	}
	ret, ok := themap[path[len(path)-1]]
	if (!ok) {
		panic("incorrect path")
	}
	return ret
}
