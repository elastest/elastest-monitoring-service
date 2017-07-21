package main

import "strings"
import "fmt"

type JSONPath string
type Channel string
type Param string
type SignalName string

type Signal interface {
	// GetValue(time uint64) interface{} 
	// Actually, Signal X should have
	// GetValue :: Time -> X 
	Sample() *interface{}
	// Actually, Signal X should have Sample :: () -> X.  It is GetValue(now)
}

type SampledSignal struct {
	latestValue *interface{} // SampledSignal X, latestValue :: X
}

func (ssignal SampledSignal) Sample() *interface{} {
	return ssignal.latestValue
}

type AggregatedSignal struct {
	signalsFamily []Signal
	aggregationFun func(vals []interface{}) *interface{} // combinator
}

func (aggSignal *AggregatedSignal) AddSource(srcSignal Signal) {
	aggSignal.signalsFamily = append(aggSignal.signalsFamily, srcSignal)
}

func (aggSignal AggregatedSignal) Sample() *interface{} {
	signals := aggSignal.signalsFamily
	vals := make([]interface{}, len(signals))
	i:=0
	for _, signal := range signals {
		val:=signal.Sample()
		if (val !=nil) {
			vals[i] = *val
			i++
		}
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

var UNDEFINED *interface{} = nil

func (csignal ConditionalSignal) Sample() *interface{} {

	fmt.Printf("csignal: %p\n", csignal.conditionSignal)
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

var theGlobalBaseSessionDefs = []BaseSessionDefinition {
	BaseSessionDefinition {
		"timeIsEven",
		[]EventDefinition {
			EventDefinition {
				"in_condition_true",
				nil,
			},
		},
		[]EventDefinition {
			EventDefinition {
				"in_condition_false",
				nil,
			},
		},
	},
}

func checkSamples(evt Event) {
	for _, ssignaldef := range theGlobalSampledSignalDefs {
		if (evt.Channel == ssignaldef.inChannel) {
			sampledSignal := SignalNameAndPars{ssignaldef.name, extractParamsMap(evt, ssignaldef.paramsPaths)}
			value := extractFromMap(evt.Payload, ssignaldef.valuePath)
			sampledSignals = append(sampledSignals, sampledSignal)
			reportSample(sampledSignal, value)
		}
	}
	for _, bsessiondef := range theGlobalBaseSessionDefs {
		// first inhibitors, so activators can override them if needed
		for _, inhEvent := range bsessiondef.inhibitorEvents {
			if (evt.Channel == inhEvent.inChannel) {
				inhibitedSession := SignalNameAndPars{bsessiondef.name, extractParamsMap(evt, inhEvent.paramsPaths)}
				updateBaseSession(inhibitedSession, false)
			}
		}
		for _, actEvent := range bsessiondef.activatorEvents {
			if (evt.Channel == actEvent.inChannel) {
				activatedSession := SignalNameAndPars{bsessiondef.name, extractParamsMap(evt, actEvent.paramsPaths)}
				updateBaseSession(activatedSession, true)
			}
		}
	}
}

func extractParamsMap(evt Event, paramsPaths map[Param]JSONPath) map[Param]string {
	paramsMap := map[Param]string{}
	for param, path := range paramsPaths {
		paramsMap[param] = extractFromMap(evt.Payload, path).(string)
	}
	return paramsMap
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
