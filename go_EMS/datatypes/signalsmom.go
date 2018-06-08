package data

import (
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
	"time"
)

type SignalsDefinitions struct {
    Type string `json:"type"`
    Def string `json:"def"`
}

type Sampler struct {
    OutChan chan striverdt.Event
}

type MoMEngine01 struct {
    Sampler Sampler
    Striverkillchan chan bool
}

func (engine MoMEngine01) Kill() {
    close(engine.Sampler.OutChan)
    close(engine.Striverkillchan)
}

func (s Sampler) ProcessEvent(evt Event) {
    payload := striverdt.Some(evt)
    t, err := time.Parse(time.RFC3339Nano,evt.Timestamp)
    if err != nil {
        panic(err)
    }
    striverEvent := striverdt.Event{striverdt.Time(t.Unix()), payload}
    s.OutChan <- striverEvent
}
