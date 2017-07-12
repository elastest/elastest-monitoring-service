package main

import "reflect"

type ParamAndValue struct {
	param Param
	value JSONPath
}

type SampledSignalDefinition struct {
	name SignalName
	paramsPaths map[Param]JSONPath
	inChannel Channel
	valuePath JSONPath
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
	funcName string
	signalsFamily SignalName
	params []Param
	quantifiedParams []Param // in signalsFamily's definition
	freeParams map[Param]Param
	// Every signalsFamily's not-quantified param must be a key of freeParams,
	// and every member of params must be a value
}

func (aggsignaldef *AggregatedSignalDefinition) getParams() []Param {
	return aggsignaldef.params
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
