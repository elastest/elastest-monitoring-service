package main

import "reflect"

type SampledSignalDefinition struct {
	name SignalName
	paramsPaths map[Param]JSONPath
	inChannel Channel
	valuePath JSONPath
	// add predicate over events
}

func (ssignaldef *SampledSignalDefinition) getParams() []Param {
	keys := reflect.ValueOf(ssignaldef.paramsPaths).MapKeys()
	strkeys := make([]Param, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = Param(keys[i].String())
	}
	return strkeys
}

type AggregatedSignalDefinition struct {
	name SignalName
	params []Param
	funcName string
	quantifiedParams []Param // in signalsFamily's definition
	signalsFamily SignalName
	signalParams []Param
}

func (aggsignaldef *AggregatedSignalDefinition) getParams() []Param {
	return aggsignaldef.params
}

type ConditionalSignalDefinition struct {
	name SignalName
	params []Param
	sourceSignal SignalName
	srcSignalParams []Param
	condition SessionName
	conditionParams []Param
}


type WriteValue struct {
	valtype string
	valpayload string
}

type SignalWriteDefinition struct {
	sourceSignal SignalName
	outChannel Channel
	fieldsAndValues map[JSONPath]WriteValue
	// This is stringly typed. Actually, value could be either a Param, a
	// literal or "value". Maybe every Param should be present at least once.
	triggerer SNameAndRebound // Those params not present can be any
}
