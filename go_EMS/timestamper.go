// Note that this version is completely sequential and doesn't use goroutines at all
package main

import "encoding/json"
import "fmt"
import "os"
import "bufio"
import "time"

func processMap(dasmap map[string]interface{}, offset int) {

	if _, ok := dasmap["@timestamp"]; ok {
		now := time.Now()
		dasmap["@timestamp"] = now.UTC().Format("2006-01-02T15:04:05.000") + "Z"
	}

}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
    var dasmap map[string]interface{}
	i := 0
	for ; scanner.Scan(); i++ {
		thetextbytes := []byte(scanner.Text())

		if err := json.Unmarshal(thetextbytes, &dasmap); err != nil {
			fmt.Println("No JSON. Error: " + err.Error())
		} else {
			processMap(dasmap, i);
			newJSON, _ := json.Marshal(dasmap)
			fmt.Println(string(newJSON))
		}
	}
}

