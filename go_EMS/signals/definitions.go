package signals

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    //"fmt"
)

type SampledSignalDefinition struct {
    Name dt.SignalName
    InChannel dt.Channel
    ValuePath dt.JSONPath
}

type ConditionalSignalDefinition struct {
    Name dt.SignalName
    SourceSignal dt.SignalName
    Condition dt.SignalName
}

type WriteValue struct {
    Valtype string
    Valpayload string
}

type SignalWriteDefinition struct {
    SourceSignal dt.SignalName
    OutChannel dt.Channel
    //FieldsAndValues map[dt.JSONPath]dt.WriteValue
    // This is stringly typed. Actually, value could be either a Param, a
    // literal or "value". Maybe every Param should be present at least once.
}
