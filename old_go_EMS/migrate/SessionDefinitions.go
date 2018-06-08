package main

import "reflect"

type BaseSessionDefinition struct {
	Name SignalName
	ActivatorEvents []EventDefinition// cannot be empty
	InhibitorEvents []EventDefinition
}

type EventDefinition struct {
	InChannel Channel
	ParamsPaths map[Param]JSONPath // May be partial function for inhibitor evs
	// eventProperties func(e Event) bool <- add later, also for sampled signals
}

func (basedef *BaseSessionDefinition) getParams() []Param {
	paramsPaths := basedef.ActivatorEvents[0].ParamsPaths
	keys := reflect.ValueOf(paramsPaths).MapKeys()
	strkeys := make([]Param, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = Param(keys[i].String())
	}
	return strkeys
}

type CompositeSessionDefinition struct {
	/*
	TODO
	name SignalName
	funcName string
	signalsFamily SignalName
	params []Param
	quantifiedParams []Param // in signalsFamily's definition
	freeParams map[Param]Param
	// Every signalsFamily's not-quantified param must be a key of freeParams,
	// and every member of params must be a value
	*/
}

func (compositedef *CompositeSessionDefinition) getParams() []Param {
	return nil// TODO
}
