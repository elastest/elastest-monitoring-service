package signals

import (
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
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
    Samplers []dt.Sampler
    OutStreams []striverdt.OutStream
    InStreams []striverdt.InStream
}


func (visitor *SignalToStriverVisitor) visitSampled(sampledsignal SampledSignalDefinition) {
	signalchan := make(chan striverdt.Event)
    visitor.Samplers = append(visitor.Samplers, dt.Sampler{signalchan})
	inStream := striverdt.InStream{striverdt.StreamName(sampledsignal.Name), &striverdt.InFromChannel{signalchan, nil, 0, false}}
    visitor.InStreams = append(visitor.InStreams, inStream)
}

func (visitor *SignalToStriverVisitor) visitFuncSignal(funcsignal FuncSignalDefinition) {
    generalfun := func (args...striverdt.EvPayload) striverdt.EvPayload {

        castedargs := make([]interface{}, len(args))
        for i,arg := range args {
            castedargs[i] = arg.Val.(striverdt.EvPayload).Val
        }
        return striverdt.Some(funcsignal.FuncDef.getFunction()(castedargs...))
    }
    sourcesValNodes := make([]striverdt.ValNode, len(funcsignal.SourcesNames))
    for i,name := range funcsignal.SourcesNames {
        sourcesValNodes[i] = &striverdt.PrevEqValNode{striverdt.TNode{}, name, []striverdt.Event{}}
    }
    valnode := striverdt.FuncNode{sourcesValNodes, generalfun}
    outStream := striverdt.OutStream{funcsignal.Name, striverdt.SrcTickerNode{funcsignal.SourcesNames[0]}, valnode}

    visitor.OutStreams = append(visitor.OutStreams, outStream)
}

func (visitor *SignalToStriverVisitor) visitConditionalAvgSignal(funcsignal ConditionalAvgSignalDefinition) {
    condCounterName := "condcounter::"+funcsignal.Name
    condCounterFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        cond := args[0]
        if !cond.IsSet || !cond.Val.(striverdt.EvPayload).Val.(bool) {
            return striverdt.NothingPayload
        }
        myprev := args[1]
        prev := 0
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(int)
        }
        return striverdt.Some(prev+1)
    }
    condCounterVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, funcsignal.Condition, []striverdt.Event{}},
        &striverdt.PrevValNode{striverdt.TNode{}, condCounterName, []striverdt.Event{}},
    }, condCounterFun}
    condCounterStream := striverdt.OutStream{condCounterName, striverdt.SrcTickerNode{funcsignal.SourceSignal}, condCounterVal}
    visitor.OutStreams = append(visitor.OutStreams, condCounterStream)

    condAvgFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        cond := args[0]
        if !cond.IsSet || !cond.Val.(striverdt.EvPayload).Val.(bool) {
            return striverdt.NothingPayload
        }
        myprev := args[1]
        cpuval := args[2].Val.(striverdt.EvPayload).Val.(float32)
        kplusone := float32(args[3].Val.(striverdt.EvPayload).Val.(int))
        prev := float32(0.0)
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(float32)
        }
        res := (prev*(kplusone-1)+cpuval)/kplusone
        return striverdt.Some(res)
    }
    condAvgVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, funcsignal.Condition, []striverdt.Event{}},
        &striverdt.PrevValNode{striverdt.TNode{}, funcsignal.Name, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, funcsignal.SourceSignal, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, condCounterName, []striverdt.Event{}},
    }, condAvgFun}
    condAvgStream := striverdt.OutStream{funcsignal.Name, striverdt.SrcTickerNode{funcsignal.SourceSignal}, condAvgVal}
    visitor.OutStreams = append(visitor.OutStreams, condAvgStream)
}
