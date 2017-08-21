package main

import "errors"
import "encoding/json"
import "fmt"

type SNameAndRebound struct {
	signalName SignalName
	reboundParameters map[Param]Param
}

type SNameAndBiRebound struct {
	signalName SignalName
	reboundMetricParameters map[Param]Param
	reboundSessionParameters map[Param]Param
}

type SignalNameAndPars struct {
	signalName SignalName
	parameters map[Param]string
}

type SignalIdToSampledSignal struct {
	sigid SignalNameAndPars
	signal *SampledSignal
}

type SignalIdToAggregatedSignal struct {
	sigid SignalNameAndPars
	signal *AggregatedSignal
}

type SignalIdToConditionalSignal struct {
	sigid SignalNameAndPars
	signal *ConditionalSignal
}

type SignalIdToSignal struct {
	sigid SignalNameAndPars
	signal Signal
}

type MetricParsAndSignal struct {
	params map[Param]string
	signal Signal
}

func (sigida SignalNameAndPars) equals(sigidb SignalNameAndPars) bool {
	if (sigida.signalName == sigidb.signalName) {
		paramsa := sigida.parameters
		paramsb := sigidb.parameters
		if (len(paramsa) == len(paramsb)) {
			for k, v := range paramsa {
				if (v != paramsb[k]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

var sampledSignalMan []*SignalIdToSampledSignal
var aggregatedSignalMan []*SignalIdToAggregatedSignal
var conditionalSignalMan []*SignalIdToConditionalSignal

func getSignals(signalName SignalName, boundParams map[Param]string) []MetricParsAndSignal {
	var ret []MetricParsAndSignal = nil
	sigidsandsignals := make([]SignalIdToSignal, len(sampledSignalMan) + len(aggregatedSignalMan) + len(conditionalSignalMan))
	for _, sigidandsignal := range sampledSignalMan {
		sigidsandsignals = append(sigidsandsignals, SignalIdToSignal{sigidandsignal.sigid, sigidandsignal.signal})
	}
	for _, sigidandsignal := range aggregatedSignalMan {
		sigidsandsignals = append(sigidsandsignals, SignalIdToSignal{sigidandsignal.sigid, sigidandsignal.signal})
	}
	for _, sigidandsignal := range conditionalSignalMan {
		sigidsandsignals = append(sigidsandsignals, SignalIdToSignal{sigidandsignal.sigid, sigidandsignal.signal})
	}
	for _, sigidandsignal := range sigidsandsignals {
		sigid := sigidandsignal.sigid
		add := false
		if (sigid.signalName == signalName) {
			add = true
			for p,v := range boundParams {
				if (sigid.parameters[p] != v) {
					add = false
				}
			}
		}
		if (add) {
			ret = append(ret, MetricParsAndSignal{sigid.parameters, sigidandsignal.signal})
		}
	}
	return ret
}

func registerSampledSignal(signalid SignalNameAndPars, signal *SampledSignal) error {
	for _, entry := range sampledSignalMan {
		if (signalid.equals(entry.sigid)) {
			return errors.New("entry already exists")
		}
	}
	sampledSignalMan = append(sampledSignalMan, &SignalIdToSampledSignal{signalid, signal})
	return nil
}

func getSampledSignal(signalid SignalNameAndPars) (*SampledSignal, error) {
	for _, entry := range sampledSignalMan {
		if (signalid.equals(entry.sigid)) {
			return entry.signal, nil
		}
	}
	return &SampledSignal{}, errors.New("no such entry")
}

func registerAggregatedSignal(signalid SignalNameAndPars, signal *AggregatedSignal) error {
	for _, entry := range aggregatedSignalMan {
		if (signalid.equals(entry.sigid)) {
			return errors.New("entry already exists")
		}
	}
	aggregatedSignalMan = append(aggregatedSignalMan, &SignalIdToAggregatedSignal{signalid, signal})
	return nil
}

func getAggregatedSignal(signalid SignalNameAndPars) (*AggregatedSignal, error) {
	for _, entry := range aggregatedSignalMan {
		if (signalid.equals(entry.sigid)) {
			return entry.signal, nil
		}
	}
	return &AggregatedSignal{}, errors.New("no such entry")
}

func registerConditionalSignal(signalid SignalNameAndPars, signal *ConditionalSignal) error {
	for _, entry := range conditionalSignalMan {
		if (signalid.equals(entry.sigid)) {
			return errors.New("entry already exists")
		}
	}
	conditionalSignalMan = append(conditionalSignalMan, &SignalIdToConditionalSignal{signalid, signal})
	return nil
}

var aggregatedSignalCreationMap = map[SignalName][]SNameAndRebound {
	/*"cpuload": []SNameAndRebound {
		SNameAndRebound{"avgcpuload", map[Param]Param{}},
	},*/
}

var conditionalSignalCreationMap = map[SignalName][]SNameAndBiRebound {
	/*"cpuload" : []SNameAndBiRebound {
		SNameAndBiRebound{"condcpuload", map[Param]Param{"x":"x"}, nil},
	},
	"timeIsEven" : []SNameAndBiRebound {
		SNameAndBiRebound{"condcpuload", map[Param]Param{"x":"x"}, nil},
	},*/
}

func reportSample(signalpars SignalNameAndPars, value interface{}) {
	theSignal, err := getSampledSignal(signalpars)
	if (err != nil) {
		theSignal = createSampledSignal(signalpars)
	}
	theSignal.latestValue = &value
}

// TODO change to array maybe (to avoid repetition)?
var theGlobalAggregatedSignalDefs = map[SignalName]AggregatedSignalDefinition {
	/*
	"avgcpuload" : AggregatedSignalDefinition {
			"avgcpuload",
			[]Param{},
			"avg",
			[]Param{"x"},
			"cpuload",
			[]Param{"x"},
		},
	*/
}

var theGlobalConditionalSignalDefs = map[SignalName]ConditionalSignalDefinition {
	/*
	"condcpuload" : ConditionalSignalDefinition {
			"condcpuload",
			[]Param{"x"},
			"cpuload",
			[]Param{"x"},
			"timeIsEven",
			nil,
	},
	*/
}


var aggregatorsMap = map[string] (func(vals []interface{}) *interface{}) {
	"avg": func (vals []interface{}) *interface{} {
		if len(vals) == 0 {
			return UNDEFINED
		}
		res := 0.0
		for _, val := range vals {
			res += val.(float64)
		}
		var ret interface{} = res / float64(len(vals))
		return &ret
	},
}

func createSampledSignal(signalpars SignalNameAndPars) *SampledSignal {
	ret := &SampledSignal{nil}
	err := registerSampledSignal(signalpars, ret)
	if (err!= nil) {
		panic(err)
	}
	reportSignalCreation(signalpars, ret)
	return ret
}

func createAggregatedSignal(signalpars SignalNameAndPars) *AggregatedSignal {
	aggDef, ok := theGlobalAggregatedSignalDefs[signalpars.signalName]
	// assert ok
	if (!ok) {
		// error
		panic("nosuchaggsig")
	}
	aggfun, ok := aggregatorsMap[aggDef.funcName]
	// assert ok
	if (!ok) {
		panic("nosuchaggfun")
	}
	ret := &AggregatedSignal{make([]Signal,0), aggfun}
	err := registerAggregatedSignal(signalpars, ret)
	if (err!= nil) {
		panic(err)
	}
	reportSignalCreation(signalpars, ret)
	return ret
}

func createConditionalSignal(signalpars SignalNameAndPars, sessionSignal SessionSignal, metricSignal Signal) {
	csignal := &ConditionalSignal{metricSignal, sessionSignal}
	err := registerConditionalSignal(signalpars, csignal)
	if (err!= nil) {
		panic(err)
	}
	reportSignalCreation(signalpars, csignal)
}

func reportSignalCreation(srcSignalId SignalNameAndPars, srcSignal Signal) {
	// it also triggers write def creations
	registerWriteDefs(srcSignalId, srcSignal)

	sName := srcSignalId.signalName
	sPars := srcSignalId.parameters

	// create conditional signals
	arr, ok := conditionalSignalCreationMap[sName]
	if (ok) {
		for _, inducedSignal := range arr {
			theDefinition,ok := theGlobalConditionalSignalDefs[inducedSignal.signalName]
			// assert ok
			if (!ok) {
				// error
				panic("nosuchcondsig")
			}

			signalBoundParams := make(map[Param]string)
			for srcParam, myParam := range inducedSignal.reboundMetricParameters {
				signalBoundParams[myParam] = sPars[srcParam]
			}

			conditionBoundParams := make(map[Param]string)
			for sessParam, myParam := range inducedSignal.reboundSessionParameters {
				val, ok := signalBoundParams[myParam]
				if (ok) {
					conditionBoundParams[sessParam] = val
				}
			}

			sessionParsAndSignals := getSessionSignals(theDefinition.condition, conditionBoundParams)

			for _, sessParsAndSignal := range sessionParsAndSignals {
				paramvals := make(map[Param]string)
				for k,v := range signalBoundParams {
					paramvals[k] = v
				}
				for k,v := range sessParsAndSignal.params {
					paramvals[inducedSignal.reboundSessionParameters[k]] = v
				}
				// assert paramvals are all the parameters
				nameAndPars := SignalNameAndPars{inducedSignal.signalName, paramvals}
				createConditionalSignal(nameAndPars, sessParsAndSignal.signal, srcSignal)
			}
		}
	}

	// create aggregated signals and add source
	aggSignals, ok := aggregatedSignalCreationMap[sName]
	if (ok) {
		for _, inducedSignal := range aggSignals {
			theDefinition,ok := theGlobalAggregatedSignalDefs[inducedSignal.signalName]
			// assert ok
			if (!ok) {
				// error
				panic("nosuchaggsig")
			}
			defParams := theDefinition.getParams()
			paramvals := make(map[Param]string)
			for _, param := range defParams {
				paramvals[param] = sPars[inducedSignal.reboundParameters[param]]
			}
			signalid := SignalNameAndPars{inducedSignal.signalName, paramvals}
			theAggSignal, err:= getAggregatedSignal(signalid)
			if (err != nil) {
				theAggSignal = createAggregatedSignal(signalid)
			}
			theAggSignal.AddSource(srcSignal)
		}
	}
}

var theGlobalWriteDefs = []SignalWriteDefinition {
	/*
	// write def of avg
	SignalWriteDefinition{
		"avgcpuload",
		"out",
		map[JSONPath]WriteValue{
			"load" : WriteValue{"value", ""},
			"hostname" : WriteValue{"literal", "average"},
		},
		SNameAndRebound{
			"cpuload",
			map[Param]Param{},
		},
	},
	// write def of cpuload
	SignalWriteDefinition{
		"cpuload",
		"out",
		map[JSONPath]WriteValue{
			"metric" : WriteValue{"literal", "cpuload"},
			"load" : WriteValue{"value", ""},
			"hostname" : WriteValue{"param", "x"},
		},
		SNameAndRebound{
			"cpuload",
			map[Param]Param{
				"x" : "x",
			},
		},
	},
	// write def of condcpuload
	SignalWriteDefinition{
		"condcpuload",
		"out",
		map[JSONPath]WriteValue{
			"metric" : WriteValue{"literal", "condcpuload"},
			"load" : WriteValue{"value", ""},
			"hostname" : WriteValue{"param", "x"},
		},
		SNameAndRebound{
			"cpuload",
			map[Param]Param{
				"x" : "x",
			},
		},
	},
	*/
}

func registerWriteDefs(srcSignalId SignalNameAndPars, srcSignal Signal) {
	for _, wd := range theGlobalWriteDefs {
		if (wd.sourceSignal == srcSignalId.signalName) {
			createWriter(srcSignal, srcSignalId.parameters, wd)
		}
	}
}

type SignalParsToWriter struct {
	triggerer SignalNameAndPars
	writer func(timestamp string)
}

var theGlobalWriters []SignalParsToWriter

func createWriter(srcSignal Signal, parameters map[Param]string, writedef SignalWriteDefinition) {
	thefun := func(timestamp string) {
			value := srcSignal.Sample()
			dasmap := map[JSONPath]interface{} {}
			for f, v := range writedef.fieldsAndValues {
				switch v.valtype {
					case "value":
						if (value == nil) {
							dasmap[f] = 0 // should get the "zero value" from the type of the signal
						} else {
							dasmap[f] = *value // should do a proper structure for dasmap <- I no longer know what I meant with this comment
						}
					case "param":
						dasmap[f] = parameters[Param(v.valpayload)]
					case "literal":
						dasmap[f] = v.valpayload
				}
			}
			dasmap["@timestamp"] = timestamp
			newJSON, _ := json.Marshal(dasmap)
			fmt.Println(string(newJSON))
		}
		// aca
		reboundMap := map[Param]string {}
		for k,v := range writedef.triggerer.reboundParameters {
			reboundMap[k] = parameters[v]
		}
		snameAndPars := SignalNameAndPars{writedef.triggerer.signalName, reboundMap}
		theGlobalWriters = append(theGlobalWriters, SignalParsToWriter{snameAndPars, thefun})
}

func getWriters(signalid SignalNameAndPars) []func(timestamp string) {
	var ret []func(timestamp string)
	for _, sigxwriter := range theGlobalWriters {
		triggerer := sigxwriter.triggerer
		if (signalid.signalName == triggerer.signalName) {
			matches := true
			for param, tval := range triggerer.parameters {
				matches = matches && tval == signalid.parameters[param] // could be shortcircuited
			}
			if (matches) {
				ret = append(ret, sigxwriter.writer)
			}
		}
	}
	return ret
}
