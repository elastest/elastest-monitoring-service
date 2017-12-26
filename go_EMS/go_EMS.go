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
	}
	panic("leaving!")
}

func scanStdIn(file io.Reader) {
	scanner := bufio.NewScanner(file)
    var dasmap map[string]interface{}
	for scanner.Scan() {
		dasmap = nil
		thetextbytes := []byte(scanner.Text())

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
			fmt.Println(string(thetextbytes))
		}
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
