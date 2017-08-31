package main

import "reflect"

type SampledSignalDefinition struct {
	Name SignalName
	ParamsPaths map[Param]JSONPath
	InChannel Channel
	ValuePath JSONPath
	// add predicate over events
}

func (ssignaldef *SampledSignalDefinition) getParams() []Param {
	keys := reflect.ValueOf(ssignaldef.ParamsPaths).MapKeys()
	strkeys := make([]Param, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = Param(keys[i].String())
	}
	return strkeys
}

type AggregatedSignalDefinition struct {
	Name SignalName
	Params []Param
	FuncName string
	QuantifiedParams []Param // in signalsFamily's definition
	SignalsFamily SignalName
	SignalParams []Param
}

func (aggsignaldef *AggregatedSignalDefinition) getParams() []Param {
	return aggsignaldef.Params
}

type ConditionalSignalDefinition struct {
	Name SignalName
	Params []Param
	SourceSignal SignalName
	SrcSignalParams []Param
	Condition SignalName
	ConditionParams []Param
}


type WriteValue struct {
	Valtype string
	Valpayload string
}

type SignalWriteDefinition struct {
	SourceSignal SignalName
	OutChannel Channel
	FieldsAndValues map[JSONPath]WriteValue
	// This is stringly typed. Actually, value could be either a Param, a
	// literal or "value". Maybe every Param should be present at least once.
	Triggerer SNameAndRebound // Those params not present can be any
}
