// Note that this version is completely sequential and doesn't use goroutines at all
package main

import "encoding/json"
import "fmt"
import "os"
import "bufio"
import "io"

func main() {
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
    var dasmap map[string]interface{}
    i:=0
	for scanner.Scan() {
		dasmap = nil
		thetextbytes := []byte(scanner.Text())
        fmt.Println("Read event ", i)

		if err := json.Unmarshal(thetextbytes, &dasmap); err != nil {
			fmt.Println("No JSON. Error: " + err.Error())
		} else {
			evt := getEvent(dasmap)
			checkSamples(evt)
			checkWriteDefs(evt.Timestamp)
			if (evt.Channel == "undefined") {
				newJSON, _ := json.Marshal(evt)
				newJSON = newJSON
				//fmt.Println(string(newJSON))
			}
			//newJSON, _ := json.Marshal(dasmap)
            evstring := string(thetextbytes)+"\n"
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
			readAndRegister(dasmap)
		}
	}
}
