package signals

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    //"fmt"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
	"time"
    "strings"
)

type Sampler struct {
    InChannel dt.Channel
    ValuePath dt.JSONPath
    OutChan chan striverdt.Event
}

func (s Sampler) ProcessEvent(evt dt.Event) {
    payload := striverdt.NothingPayload
    if sets.SetIn(s.InChannel, evt.Channels) {
        payload = striverdt.Some(extractFromMap(evt.Payload, s.ValuePath))
    }
    t, err := time.Parse(time.RFC3339Nano,evt.Timestamp)
    if err != nil {
        panic(err)
    }
    striverEvent := striverdt.Event{striverdt.Time(t.Unix()), payload}
    s.OutChan <- striverEvent
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
