package signals

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    //"fmt"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
	"time"
)

func SamplerProcessEvent(s dt.Sampler, evt dt.Event) {
    payload := striverdt.NothingPayload
    if sets.SetIn(s.InChannel, evt.Channels) {
        pl,_ := jsonrw.ExtractFromMap(evt.Payload, s.ValuePath)
        payload = striverdt.Some(pl)
    }
    t, err := time.Parse(time.RFC3339Nano,evt.Timestamp)
    if err != nil {
        panic(err)
    }
    striverEvent := striverdt.Event{striverdt.Time(t.Unix()), payload}
    s.OutChan <- striverEvent
}
