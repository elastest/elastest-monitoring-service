package datatypes

import "fmt"

// interface
type TickerNode interface {
    Vote (t Time) MaybeTime;
    Exec (t Time, inpipes InPipes) EvPayload;
    Rinse (inpipes InPipes)
}

type ValNode interface {
    Exec (t Time, w interface{}, inpipes InPipes) EvPayload;
    Rinse (inpipes InPipes)
}

// tickers
type ConstTickerNode struct {
    ConstT Time
    ConstW interface{}
}

type SrcTickerNode struct {
    SrcStream StreamName
}

type DelayTickerNode struct {
    SrcStream StreamName
    Combiner func(a EvPayload, b EvPayload) EvPayload;
    Alarms []Event
}

type UnionTickerNode struct {
    LeftTicker TickerNode
    RightTicker TickerNode
    Combiner func(a EvPayload, b EvPayload) EvPayload
}

// values

type TNode struct {
}

type WNode struct {
}

type PrevNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type PrevEqNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type PrevValNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type PrevEqValNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type FuncNode struct {
    ArgNodes []ValNode
    Innerfun func (args ...EvPayload) EvPayload
}

type InFromChannel struct {
    InChannel chan Event
    NextEvent *Event
    MinT Time
    Closed bool
}

func (ticker *InFromChannel) PeekNextTime () MaybeTime {
    if ticker.Closed {
        return NothingTime
    }
    if ticker.NextEvent == nil {
        nextEv, open := <-ticker.InChannel
        if !open {
            fmt.Println("Closing input from channel")
            ticker.Closed = true
            return NothingTime
        }
        if (nextEv.Time <= ticker.MinT) {
            // Ensure safety of input events
            nextEv.Time = ticker.MinT +1
        }
        ticker.NextEvent = &nextEv
        ticker.MinT = nextEv.Time
    }
    return SomeTime(ticker.NextEvent.Time)
}

func (ticker *InFromChannel) Exec (t Time) EvPayload {
    if ticker.Closed {
        return NothingPayload
    }
    if t == ticker.NextEvent.Time {
        ret := ticker.NextEvent.Payload
        ticker.NextEvent = nil
        return ret
    }
    return NothingPayload
}
