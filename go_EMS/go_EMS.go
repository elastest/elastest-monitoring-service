// Note that this version is completely sequential and doesn't use goroutines at all
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "bufio"
    "io"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    et "github.com/elastest/elastest-monitoring-service/go_EMS/eventproc"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
    internalsv "github.com/elastest/elastest-monitoring-service/go_EMS/internalapiserver"
)

func main() {
    fmt.Println("Serving server")
    go internalsv.Serve()
    fmt.Println("Server served. Starting scans")
	go openAndLoop("/usr/share/logstash/pipes/swagpipe",scanAPIPipe)
	openAndLoop("/usr/share/logstash/pipes/leftpipe",scanStdIn)
}


func openAndLoop(pipename string, callback func(reader io.Reader)) {
	file, err := os.Open(pipename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

	for {
		callback(file)
        fmt.Println("RELOADING " + pipename)
	}
	panic("leaving!")
}

func scanStdIn(file io.Reader) {
    // Opening staticout
    staticout := os.Args[1]
    fstatic, err := os.OpenFile(staticout, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer fstatic.Close()

    // Opening dynout
    dynout := os.Args[2]
    fdyn, err := os.OpenFile(dynout, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer fdyn.Close()

	scanner := bufio.NewScanner(file)
    var rawEvent map[string]interface{}
    i:=0
	for scanner.Scan() {
		rawEvent = nil
		thetextbytes := []byte(scanner.Text())
        fmt.Println("Read event ", i)

		if err := json.Unmarshal(thetextbytes, &rawEvent); err != nil {
			fmt.Println("No JSON. Error: " + err.Error())
		} else {
            var evt dt.Event = jsonrw.ReadEvent(rawEvent)
            et.TagEvent(&evt)
            newJSON, _ := json.Marshal(evt)
			//newJSON, _ := json.Marshal(rawEvent)
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
        fmt.Println("Processed event ", i)
        i=i+1
	}
}

func scanAPIPipe(file io.Reader) {
	scanner := bufio.NewScanner(file)
    var dasmap map[string]interface{}
	for scanner.Scan()  {
		dasmap = nil
		thetextbytes := []byte(scanner.Text())

		if err := json.Unmarshal(thetextbytes, &dasmap); err != nil {
			panic("No JSON. Error: " + err.Error())
		} else {
			// readAndRegister(dasmap)
		}
	}
}
