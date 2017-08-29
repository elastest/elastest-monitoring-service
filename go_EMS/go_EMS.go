// Note that this version is completely sequential and doesn't use goroutines at all
package main

import "encoding/json"
import "fmt"
import "os"
import "io"
import "bufio"

func main() {
	go scanAPIPipe("/usr/share/logstash/pipes/swagpipe")
	scanStdIn(os.Stdin)
}

func scanStdIn(fdes io.Reader) {
	scanner := bufio.NewScanner(fdes)
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
			//fmt.Println(string(newJSON))
		}
	}
}

func scanAPIPipe(pipename string) {
	file, err := os.Open(pipename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var dasmap map[string]interface{}
	//for {
		scanner := bufio.NewScanner(file)
			for scanner.Scan()  {
				dasmap = nil
				thetextbytes := []byte(scanner.Text())

				if err := json.Unmarshal(thetextbytes, &dasmap); err != nil {
					panic("No JSON. Error: " + err.Error())
				} else {
					//readAndRegister(dasmap)
				}
		}
	//}
	//leaving
}
