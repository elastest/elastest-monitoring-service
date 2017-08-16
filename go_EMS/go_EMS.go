// Note that this version is completely sequential and doesn't use goroutines at all
package main

import "encoding/json"
import "fmt"
import "os"
import "bufio"

func main() {

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
			//newJSON, _ := json.Marshal(dasmap)
			//fmt.Println(string(newJSON))
		}
	}
}
