package main

import (
	"testing"
	"encoding/json"
)

var theEvent Event

func TestChannelInference(t *testing.T) {

	tables := []struct {
        json string
        channel string
    }{
        {"{\"channel\":\"algo\"}", "algo"},
        {"{\"otherfields\":\"somevals\"}", "undefined"},
    }

    for _, table := range tables {
		var rawEvent map[string]interface{} = nil
		thejson := []byte(table.json)
		json.Unmarshal(thejson, &rawEvent)
		theEvent = getEvent(rawEvent)
		inferredchan := string(theEvent.Channel)
		if inferredchan != table.channel {
			t.Errorf("Wrong inferred channel, got: %s, want: %s.", inferredchan, table.channel)
		}
	}
}

func TestSumAggSignal(t *testing.T) {
	aggSignal := AggregatedSignal{[]Signal{}, aggregatorsMap["avg"]}
	sample := aggSignal.Sample()
	if (sample != UNDEFINED) {
			t.Errorf("Wrong aggregated value, got: %s, want: %s.", sample, UNDEFINED)
	}
}

var five interface{} = 5
var sampledSignal SampledSignal = SampledSignal{ &five }

func TestSampleSampled(t *testing.T) {
	sample := sampledSignal.Sample()
	if ((*sample).(int) != 5) {
		t.Errorf("Wrong sampledsignal sample value, got: %v, want: %v.", sample, 5)
	}
}

var aggSignal AggregatedSignal = AggregatedSignal{ nil, func(vals []interface{}) *interface{} {return nil} }

func TestAddSource(t *testing.T) {
	aggSignal.AddSource(sampledSignal)
	aggSignal.Sample()
}

var baseSession BaseSessionSignal = BaseSessionSignal{true}
var condSignal ConditionalSignal = ConditionalSignal{aggSignal, baseSession}
var signalNameAndPars SignalNameAndPars = SignalNameAndPars{"signal", map[Param]string{}}

func TestCreate(t *testing.T) {
	theGlobalAggregatedSignalDefs = map[SignalName]AggregatedSignalDefinition {
	"signal" : AggregatedSignalDefinition {
				"signal",
				[]Param{},
				"avg",
				[]Param{"x"},
				"cpuload",
				[]Param{"x"},
			},
	}
	createBaseSession(signalNameAndPars)
	createSampledSignal(signalNameAndPars)
	createConditionalSignal(signalNameAndPars, baseSession, aggSignal)
	createAggregatedSignal(signalNameAndPars)
}

func TestRest(t *testing.T) {
	condSignal.Sample()
	checkWriteDefs("ts")
	checkSamples(theEvent)
	extractParamsMap(theEvent, map[Param]JSONPath{})
	extractFromMap(map[string]interface{} {"a":5,}, "a")
	//readAndRegister(map[string]interface{}{})
	baseSession.getState()
	getSessionSignals("sessionName", map[Param]string{})
	registerBaseSessionSignal(signalNameAndPars, &baseSession)
	getBaseSession(signalNameAndPars)
	updateBaseSession(signalNameAndPars, false)
	reportSessionSignalCreation(signalNameAndPars, baseSession)
	signalNameAndPars.equals(signalNameAndPars)
	getSignals("signal", map[Param]string{})
	registerSampledSignal(signalNameAndPars, &sampledSignal)
	getSampledSignal(signalNameAndPars)
	registerAggregatedSignal(signalNameAndPars, &aggSignal)
	getAggregatedSignal(signalNameAndPars)
	registerConditionalSignal(signalNameAndPars, &condSignal)
	reportSample(signalNameAndPars, 8)
	reportSignalCreation(signalNameAndPars, aggSignal)


	theGlobalWriteDefs = []SignalWriteDefinition {
		SignalWriteDefinition{
			"signal",
			"out",
			map[JSONPath]WriteValue{
			},
			SNameAndRebound{
				"signal",
				map[Param]Param{
				},
			},
		},
	}
	registerWriteDefs(signalNameAndPars, aggSignal)
	getWriters(signalNameAndPars)
	main()
}

var ssdef SampledSignalDefinition = SampledSignalDefinition {
            "cpuload",
            map[Param]JSONPath {
                "x": "beat.hostname",
            },
            "in",
            "system.load.1",
        }

var bsdef BaseSessionDefinition = BaseSessionDefinition {
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
    }

func TestGetParams(t *testing.T) {
	ssdef.getParams()
	bsdef.getParams()
}

/*
*/
