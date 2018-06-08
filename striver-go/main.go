package main

import (
    "fmt"
    "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

func main() {

    inStreams, outStreams := shiftExample()
    // inStreams, outStreams := changePointsExample()
    //inStreams, outStreams := clockExample()
    kchan := make (chan bool)
    outchan := make (chan dt.FlowingEvent)
    controlplane.Start(inStreams, outStreams, outchan , kchan)

    fmt.Println("End of execution")
}
