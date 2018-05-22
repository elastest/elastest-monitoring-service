package eventout

import (
    "encoding/json"
    "fmt"
    "os"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    et "github.com/elastest/elastest-monitoring-service/go_EMS/eventtag"
)

var sendchan chan dt.Event

func StartSender(staticout string, dynout string) {
    sendchan = make(chan dt.Event)

    go func() {
        // Opening staticout
        fstatic, err := os.OpenFile(staticout, os.O_APPEND|os.O_WRONLY, 0600)
        if err != nil {
            panic(err)
        }
        defer fstatic.Close()

        // Opening dynout
        fdyn, err := os.OpenFile(dynout, os.O_APPEND|os.O_WRONLY, 0600)
        if err != nil {
            panic(err)
        }
        defer fdyn.Close()

        for {
            evt := <-sendchan
            et.TagEvent(&evt)
            evt.Payload["@timestamp"] = evt.Timestamp
            evt.Payload["channels"] = sets.SetToList(evt.Channels)
            newJSON, _ := json.Marshal(evt.Payload)
            evstring := string(newJSON)+"\n"
            if _, err = fstatic.WriteString(evstring); err != nil {
                panic(err)
            }
            if _, err = fdyn.WriteString(evstring); err != nil {
                fmt.Println("Broken dynamic output. Retrying...")
                fdyn, err = os.OpenFile(dynout, os.O_APPEND|os.O_WRONLY, 0600)
                if err != nil {
                    panic(err)
                }
                if _, err = fdyn.WriteString(evstring); err != nil {
                    fmt.Println("Broken retry, panicking")
                    panic(err)
                }
                fmt.Println("Recovered dyn output")
            }
        }
    }()
}

func GetSendChannel() chan dt.Event {
    return sendchan
}
