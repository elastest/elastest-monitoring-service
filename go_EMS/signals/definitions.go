package definitions

type SampledSignalDefinition struct {
    Name SignalName
    InChannel Channel
    ValuePath JSONPath
}

type ConditionalSignalDefinition struct {
    Name SignalName
    SourceSignal SignalName
    Condition SignalName
}

type WriteValue struct {
    Valtype string
    Valpayload string
}

type SignalWriteDefinition struct {
    SourceSignal SignalName
    OutChannel Channel
    FieldsAndValues map[JSONPath]WriteValue
    // This is stringly typed. Actually, value could be either a Param, a
    // literal or "value". Maybe every Param should be present at least once.
    Triggerer SNameAndRebound // Those params not present can be any
}
