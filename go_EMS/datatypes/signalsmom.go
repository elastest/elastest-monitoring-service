package data

import (
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type SignalsDefinitions struct {
    Type string `json:"type"`
    Def string `json:"def"`
}

type Sampler struct {
    InChannel Channel
    ValuePath JSONPath
    OutChan chan striverdt.Event
}

type MoMEngine01 struct {
    Samplers []Sampler
    Striverkillchan chan bool
}

func (engine MoMEngine01) Kill() {
    for _,sampler := range engine.Samplers {
        close(sampler.OutChan)
    }
    close(engine.Striverkillchan)
}
