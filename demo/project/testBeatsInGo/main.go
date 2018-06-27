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
	inSimplex := false
	inDuplex := false
	simplexReceived := 0.0
	duplexReceived := 0.0
	metrics := 0.0
	for {
		_, input, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var e event
		json.Unmarshal(input, &e)
		if e.Type == "net" {
			if strings.HasSuffix(e.ContainerName, "full-teaching-openvidu-server-kms_1") {
				metrics = metrics + 1
				if inSimplex {
					simplexReceived = simplexReceived + e.Net.TxBytesPs
				}
				if inDuplex {
					duplexReceived = duplexReceived + e.Net.TxBytesPs
				}
			}
		}
		if strings.Contains(e.Message, "STARTING SIMPLEX SESSION") {
			metrics = 0
			inSimplex = true
		}
		if strings.Contains(e.Message, "ENDING SIMPLEX SESSION") {
			inSimplex = false
			simplexReceived = simplexReceived / metrics
		}
		if strings.Contains(e.Message, "STARTING DUPLEX SESSION") {
			metrics = 0
			inDuplex = true
		}
		if strings.Contains(e.Message, "ENDING DUPLEX SESSION") {
			inDuplex = false
			duplexReceived = duplexReceived / metrics
		}
		if strings.Contains(e.Message, "Finished at") {
			log.Printf("SIMPLEX: %f, DUPLEX: %f, RANGE(%f,%f)", simplexReceived, duplexReceived, simplexReceived*1.7, simplexReceived*2.3)
			if duplexReceived < simplexReceived*1.7 && duplexReceived > simplexReceived*2.3 {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}
