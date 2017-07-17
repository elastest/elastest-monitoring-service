package main

import "reflect"

type SessionName string

type BaseSessionDefinition struct {
	name SignalName
	activatorEvents []EventDefinition// cannot be empty
	inhibitorEvents []EventDefinition
}

type EventDefinition struct {
	inChannel Channel
	paramsPaths map[Param]JSONPath // May be partial function for inhibitor evs
	// add predicate over events
}

func (basedef *BaseSessionDefinition) getParams() []Param {
	paramsPaths := basedef.activatorEvents[0].paramsPaths
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
