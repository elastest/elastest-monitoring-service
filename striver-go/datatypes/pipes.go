package datatypes

import "fmt"

type InPipes struct {
    Pipes map[StreamName]Event
    OutChan chan FlowingEvent
}

// aux methods

func (inpipes *InPipes) Reset() {
    inpipes.Pipes = make(map[StreamName]Event)
}

func (inpipes InPipes) strictConsume(streamId StreamName) Event {
    ev,ok := inpipes.Pipes[streamId]
    if !ok {
        panic("Failed strict consume on pipe "+string(streamId))
    }
    return ev
}

// endof aux methods

func (inpipes InPipes) Put(streamId StreamName, ev Event) {
    fmt.Printf("%s[%d]: ", streamId, ev.Time)
    if !ev.Payload.IsSet {
        fmt.Println("NOTICK")
    } else {
        fmt.Println(ev.Payload.Val)
        if inpipes.OutChan != nil {
            inpipes.OutChan <- FlowingEvent{streamId, ev}
        }
    }
    inpipes.Pipes[streamId] = ev
}
