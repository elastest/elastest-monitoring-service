package datatypes

type Time int64

type StreamName string

type Event struct {
    Time Time;
    Payload EvPayload
}

type EvPayload struct {
    IsSet bool;
    Val interface{}
}

type MaybeTime struct {
    IsSet bool;
    Val Time
}

func SomeTime(val Time) MaybeTime {
    return MaybeTime{true, val}
}

var NothingTime MaybeTime = MaybeTime{false, -100}

func Some(val interface{}) EvPayload {
    return EvPayload{true, val}
}

var NothingPayload EvPayload = EvPayload{false, nil}

type EpsVal struct {
    Eps Time
    Val interface{}
}

type OutStream struct {
    Name StreamName
    TicksDef TickerNode
    ValDef ValNode
}

type InStream struct {
    Name StreamName
    StreamDef InStreamDef
}

type InStreamDef interface {
    PeekNextTime() MaybeTime
    Exec(t Time) EvPayload
}

type FlowingEvent struct {
    Name StreamName
    Event Event
}
