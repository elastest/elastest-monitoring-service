package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type event struct {
	Type          string `json:"type,omitempty"`
	ContainerName string `json:"containerName,omitempty"`
	Net           *Net   `json:"net,omitempty"`
	Message       string `json:"message,omitempty"`
	EmsType       string `json:"ems_type,omitempty"`
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
        log.Printf("Received event:\n        %s\n", input)
		if e.Type == "net" {
			if strings.HasSuffix(e.ContainerName, "full-teaching-openvidu-server-kms_1") {
				//log.Printf("Checking %f\n", e.Net.TxBytesPs)
				if e.Net.TxBytesPs > 100000.0 {
					os.Exit(1)
				}
			}
		}
		if e.EmsType == "webrtc" {
			log.Printf("8182 INPUT:\n        %s\n", input)
		}
		if strings.Contains(e.Message, "Finished at") {
			os.Exit(0)
		}
	}
}
