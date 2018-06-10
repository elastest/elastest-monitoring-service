package eventproc

/*
import (
	"testing"
	"encoding/json"
	// "os"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
	sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
)

var theEvent dt.Event

func TestChannelInference(t *testing.T) {
    correctdef0 := `{"inchannels":["a"], "filter":"true", "outchannel":"C"}`
    correctdef1 := `{"inchannels":["C","a"], "filter":"true", "outchannel":"D"}`
    if reply := DeployTaggerv01(correctdef0); reply.Deploymenterror != "" {
        t.Errorf("Expected no error for %s, got %s", correctdef0, reply.Deploymenterror)
    }
    if reply := DeployTaggerv01(correctdef1); reply.Deploymenterror != "" {
        t.Errorf("Expected no error for %s, got %s", correctdef1, reply.Deploymenterror)
    }
	tables := []struct {
        json string
        channels dt.ChannelSet
    }{
        {"{\"channels\":[\"a\"]}", sets.SetFromList([]string{"a","C","D"})},
        {"{\"otherfields\":\"somevals\"}", sets.SetFromList([]string{})},
    }

    for _, table := range tables {
		var rawEvent map[string]interface{} = nil
		thejson := []byte(table.json)
		json.Unmarshal(thejson, &rawEvent)
        theEvent = jsonrw.ReadEvent(rawEvent)
		TagEvent(&theEvent)
		inferredchan := theEvent.Channels
		if !sets.SetsAreEqual(inferredchan, table.channels) {
			t.Errorf("Wrong inferred channel, got: %s, want: %s.", inferredchan, table.channels)
		}
	}
}

/*
func TestSumAggSignal(t *testing.T) {
	aggSignal := AggregatedSignal{[]Signal{}, aggregatorsMap["avg"]}
	sample := aggSignal.Sample()
	if (sample != UNDEFINED) {
			t.Errorf("Wrong aggregated value, got: %v, wanted nil.", sample)
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
var otherSignalNameAndPars SignalNameAndPars = SignalNameAndPars{"othersignal", map[Param]string{}}

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
	reportSessionSignalCreation(otherSignalNameAndPars, baseSession)
fmt.Println("")
	signalNameAndPars.equals(signalNameAndPars)
	getSignals("signal", map[Param]string{})
	registerSampledSignal(signalNameAndPars, &sampledSignal)
	getSampledSignal(signalNameAndPars)
	registerAggregatedSignal(signalNameAndPars, &aggSignal)
	getAggregatedSignal(signalNameAndPars)
	registerConditionalSignal(signalNameAndPars, &condSignal)
	reportSample(signalNameAndPars, 8)


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
	getWriters(signalNameAndPars)[0]("ts")
	//main()
	file, err := os.Open("testinputs/testdefs.json")
    if err != nil {
        panic(err)
    }
	scanAPIPipe(file)
	file, err = os.Open("testinputs/testEvents.txt")
    if err != nil {
        panic(err)
    }
    os.Args = []string{"goEMS", "/dev/null", "/dev/null"}
	scanStdIn(file)
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

*/
