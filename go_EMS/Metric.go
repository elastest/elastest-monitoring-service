package main

import "strings"

type JSONPath string
type Channel string
type Param string
type SignalName string

type Signal interface {
	// GetValue(time uint64) interface{} 
	// Actually, Signal X should have
	// GetValue :: Time -> X 
	Sample() interface{}
	// Actually, Signal X should have Sample :: () -> X.  It is GetValue(now)
}

type SampledSignal struct {
	latestValue interface{} // SampledSignal X, latestValue :: X
}

func (ssignal SampledSignal) Sample() interface{} {
	return ssignal.latestValue
}

type AggregatedSignal struct {
	signalsFamily []Signal
	aggregationFun func(vals []interface{}) interface{} // combinator
}

func (aggSignal *AggregatedSignal) AddSource(srcSignal Signal) {
	aggSignal.signalsFamily = append(aggSignal.signalsFamily, srcSignal)
}

func (aggSignal AggregatedSignal) Sample() interface{} {
	signals := aggSignal.signalsFamily
	vals := make([]interface{}, len(signals))
	for i, signal := range signals {
		vals[i] = signal.Sample()
	}
	ret := aggSignal.aggregationFun(vals)
	/*
	if len(signals) == 0 {
		return nil
	}
	ret := signals[0].Sample()
	for _, signal := range signals[1:] {
		ret = aggSignal.aggregationFun(ret, signal.Sample())
	}
	*/
	return ret
}

type ConditionalSignal struct {
	srcSignal *Signal
	conditionSignal *SessionSignal
}

var UNDEFINED interface {} = nil

func (csignal ConditionalSignal) Sample() interface{} {
	if (*csignal.conditionSignal).getState() {
		return (*csignal.srcSignal).Sample()
	}
	return UNDEFINED
}

var sampledSignals []SignalNameAndPars

func checkWriteDefs (timestamp string) {
	for _, sigid := range sampledSignals {
		for _, writer := range getWriters(sigid) {
			writer(timestamp)
		}
	}
	sampledSignals = nil
}

var theGlobalSampledSignalDefs = []SampledSignalDefinition {
	SampledSignalDefinition {
			"cpuload",
			map[Param]JSONPath {
				"x": "beat.hostname",
			},
			"in",
			"system.load.1",
		},
}

func checkSamples(evt Event) {
	for _, ssdef := range theGlobalSampledSignalDefs {
		if (evt.Channel == ssdef.inChannel) {
			sampledSignal, value := extractSignalIdAndValue(ssdef, evt)
			sampledSignals = append(sampledSignals, sampledSignal)
			reportSample(sampledSignal, value)
		}
	}
}

func extractSignalIdAndValue(ssdef SampledSignalDefinition, evt Event) (SignalNameAndPars, interface{}) {
	paramsMap := map[Param]string{}
	for param, path := range ssdef.paramsPaths {
		paramsMap[param] = extractFromMap(evt.Payload, path).(string)
	}
	value := extractFromMap(evt.Payload, ssdef.valuePath)
	return SignalNameAndPars{ssdef.name, paramsMap}, value
}

func extractFromMap(themap map[string]interface{}, strpath JSONPath) interface{} {
	path := strings.Split(string(strpath), ".")
	var ok bool
	if (len(path) == 0) {
		panic("empty path")
	}
	for _,key := range path[:len(path)-1] {
		themap, ok = themap[key].(map[string]interface{})
		if (!ok) {
			panic("incorrect path")
		}
	}
	ret, ok := themap[path[len(path)-1]]
	if (!ok) {
		panic("incorrect path")
	}
	return ret
}
