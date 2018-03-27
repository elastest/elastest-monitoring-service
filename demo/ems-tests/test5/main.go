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

type message struct {
	Message string `yaml:"status,omitempty"`
	System  System `yaml:"system,omitempty"`
}

type System struct {
	Cpu Cpu `yaml:"cpu,omitempty"`
}

type Cpu struct {
	User User `yaml:"user,omitempty"`
}

type User struct {
	Pct float64 `yaml:"pct,omitempty"`
}

func main() {
	ems := os.Getenv("ET_EMS_LSBEATS_HOST")
	ems = fmt.Sprintf("%s:3232", ems)
	u := url.URL{Scheme: "ws", Host: ems, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan message)
	iterations := 0
	state := ""
	cpu := 0.0
	items := 0.0

	go func() {
		defer close(done)
		for {
			_, input, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var m message
			json.Unmarshal(input, &m)
			if strings.Contains(m.Message, "STATUS_ON") {
				if state == "" {
					state = "on"
				} else {
					average := cpu / items
					if average > 0.5 {
						fmt.Println("High CPU in 'off' state")
						os.Exit(1)
					}
					log.Printf("OFF state ok with %f CPU\n", average)
					iterations += 1
					items = 0.0
					cpu = 0.0
				}
				log.Println("Starting ON state")
			}
			if strings.Contains(m.Message, "STATUS_OFF") {
				if state == "" {
					state = "off"
				} else {
					average := cpu / items
					if average < 0.5 {
						fmt.Println("Low CPU in 'on' state")
						os.Exit(1)
					}
					log.Printf("ON state ok with %f CPU\n", average)
					iterations += 1
					items = 0.0
					cpu = 0.0
				}
				log.Println("Starting OFF state")
			}
			if m.System.Cpu.User.Pct != 0.0 && state != "" {
				cpu += m.System.Cpu.User.Pct
				items += 1

			}
		}
	}()

	for iterations < 10 {
		time.Sleep(1 * time.Second)
	}
	log.Println("Test finished successfully!")

}
