package signals

import (
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

/*type SampledSignalDefinition struct {
    Name dt.SignalName
    InChannel dt.Channel
    ValuePath dt.JSONPath
}

type ConditionalSignalDefinition struct {
    Name dt.SignalName
    SourceSignal dt.SignalName
    Condition dt.SignalName
}

type WriteValue struct {
    Valtype string
    Valpayload string
}

type SignalWriteDefinition struct {
    SourceSignal dt.SignalName
    OutChannel dt.Channel
    //FieldsAndValues map[dt.JSONPath]dt.WriteValue
    // This is stringly typed. Actually, value could be either a Param, a
    // literal or "value". Maybe every Param should be present at least once.
}*/

type SignalToStriverVisitor struct {
    Samplers []Sampler
    OutStreams []striverdt.OutStream
    InStreams []striverdt.InStream
}


func (visitor *SignalToStriverVisitor) visitSampled(sampledsignal SampledSignalDefinition) {
	signalchan := make(chan striverdt.Event)
    visitor.Samplers = append(visitor.Samplers, Sampler{sampledsignal.InChannel, sampledsignal.ValuePath, signalchan})
	inStream := striverdt.InStream{striverdt.StreamName(sampledsignal.Name), &striverdt.InFromChannel{signalchan, nil, 0}}
    visitor.InStreams = append(visitor.InStreams, inStream)
}

func (visitor *SignalToStriverVisitor) visitFuncSignal(funcsignal FuncSignalDefinition) {
    generalfun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        sourceval := args[0].Val.(striverdt.EvPayload).Val
        return striverdt.Some(funcsignal.FuncDef.getFunction()(sourceval))
    }
    valnode := striverdt.FuncNode{[]striverdt.ValNode{&striverdt.PrevEqValNode{striverdt.TNode{}, funcsignal.SourceName, []striverdt.Event{}}}, generalfun}
    outStream := striverdt.OutStream{funcsignal.Name, striverdt.SrcTickerNode{funcsignal.SourceName}, valnode}

    visitor.OutStreams = append(visitor.OutStreams, outStream)
}
