package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type event struct {
    Channels      []string `json:"channels,omitempty"`
    Value         bool `json:"value,omitempty"`
}

//Net net structure
type Net struct {
	TxBytesPs float64 `json:"txBytes_ps,omitempty"`
}

func main() {
	ems := os.Getenv("ET_EMS_LSBEATS_HOST")
	ems = fmt.Sprintf("%s:3232", ems)
	u := url.URL{Scheme: "ws", Host: ems, Path: "/"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		time.Sleep(300 * time.Second)
	}
	defer c.Close()
	log.Printf("Done!")

	done := make(chan struct{})
	defer close(done)
	log.Printf("Looping for metrics...")
	for {
		_, input, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var e event
		json.Unmarshal(input, &e)
        log.Printf("Received event:[%s]", string(input))
        if inList("#testresult", e.Channels) {
            if e.Value {
                os.Exit(0)
            } else {
                os.Exit(1)
            }
        }
	}
}

func inList(in string, chans []string) bool {
   for _, s := range chans {
        if s == in {
            return true
        }
    }
    return false
}
