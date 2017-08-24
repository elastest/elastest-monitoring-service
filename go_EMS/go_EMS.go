// Note that this version is completely sequential and doesn't use goroutines at all
package main

import "encoding/json"
import "fmt"
import "os"
import "bufio"

func main() {
	go scanAPIPipe()
	scanStdIn()
}

func scanStdIn() {
	scanner := bufio.NewScanner(os.Stdin)
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
				fmt.Println(string(newJSON))
			}
			//newJSON, _ := json.Marshal(dasmap)
			//fmt.Println(string(newJSON))
		}
	}
}

func scanAPIPipe() {
	file, err := os.Open("/usr/share/logstash/pipes/swagpipe")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var dasmap map[string]interface{}
	for {
    scanner := bufio.NewScanner(file)
		for scanner.Scan()  {
			dasmap = nil
			thetextbytes := []byte(scanner.Text())

			if err := json.Unmarshal(thetextbytes, &dasmap); err != nil {
				fmt.Println("No JSON. Error: " + err.Error())
			} else {
				fmt.Printf("JSON read: %v\n", dasmap)
			}
	}
	}
	fmt.Println("leaving")
}
