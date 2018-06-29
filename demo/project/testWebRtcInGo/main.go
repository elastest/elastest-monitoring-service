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
	"github.com/icza/dyno"
)

type browser struct {
	Name     string
	finished bool
	AudioIn  float64
	AudioOut float64
	VideoIn  float64
	VideoOut float64
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
	inDuplex := false
	teacherBrowser := &browser{}
	studentBrowser := &browser{}
	var thisBrowser *browser
	numTeacher := 0
	numStudent := 0
	for {
		_, input, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var e interface{}
		if err := json.Unmarshal([]byte(input), &e); err != nil {
			panic(err)
		}
		message, _ := dyno.GetString(e, "message")
		component, _ := dyno.GetString(e, "component")
		if strings.Contains(component, "tss_eus_browser") {
			if teacherBrowser.Name == "" {
				log.Printf("INITIALIZING TEACHER TO %s\n", component)
				teacherBrowser.Name = component
				teacherBrowser.finished = false
			} else {
				if studentBrowser.Name == "" && component != teacherBrowser.Name {
					log.Printf("INITIALIZING STUDENT TO %s\n", component)
					studentBrowser.Name = component
					studentBrowser.finished = false
				}
			}
		}
		if strings.Contains(message, "STARTING DUPLEX SESSION") {
			if teacherBrowser.Name == "" || studentBrowser.Name == "" {
				log.Printf("DUAL SESSION AND TEACHER/STUDENT NOT INITIALIZED\n")
				os.Exit(1)
			}
			inDuplex = true
		}
		if strings.Contains(message, "ENDING DUPLEX SESSION") {
			inDuplex = false
			log.Printf("TEACHER: %s-%d-%f:%f:%f:%f\n", teacherBrowser.Name, numTeacher, teacherBrowser.AudioIn, teacherBrowser.AudioOut, teacherBrowser.VideoIn, teacherBrowser.VideoOut)
			log.Printf("STUDENT: %s-%d-%f:%f:%f:%f\n", studentBrowser.Name, numStudent, studentBrowser.AudioIn, studentBrowser.AudioOut, studentBrowser.VideoIn, studentBrowser.VideoOut)
		}

		if inDuplex && (component == teacherBrowser.Name || component == studentBrowser.Name) {
			if component == teacherBrowser.Name {
				thisBrowser = teacherBrowser
				numTeacher++
			} else {
				thisBrowser = studentBrowser
				numStudent++
			}
			typeField, _ := dyno.GetString(e, "type")
			info, _ := dyno.Get(e, typeField)
			var bytes int64
			if strings.Contains(typeField, "inbound") {
				bytes, _ = dyno.GetInteger(info, "bytesReceived")
			} else {
				bytes, _ = dyno.GetInteger(info, "bytesSent")
			}
			if strings.Contains(typeField, "inbound_audio") {
				thisBrowser.AudioIn += float64(bytes)
			}
			if strings.Contains(typeField, "outbound_audio") {
				thisBrowser.AudioOut += float64(bytes)
			}
			if strings.Contains(typeField, "inbound_video") {
				thisBrowser.VideoIn += float64(bytes)
			}
			if strings.Contains(typeField, "outbound_video") {
				thisBrowser.VideoOut += float64(bytes)
			}
		}

		if strings.Contains(message, "Finished at") {
			if studentBrowser.AudioIn < teacherBrowser.AudioOut*0.9 || studentBrowser.AudioIn > teacherBrowser.AudioOut*1.1 {
				log.Printf("ERROR in student AudioIn %f vs %f\n", studentBrowser.AudioIn, teacherBrowser.AudioOut)
				os.Exit(1)
			}
			if studentBrowser.AudioOut < teacherBrowser.AudioIn*0.9 || studentBrowser.AudioOut > teacherBrowser.AudioIn*1.1 {
				log.Printf("ERROR in student AudioOut %f vs %f\n", studentBrowser.AudioOut, teacherBrowser.AudioIn)
				os.Exit(1)
			}
			if studentBrowser.VideoIn < teacherBrowser.VideoOut*0.9 || studentBrowser.VideoIn > teacherBrowser.VideoOut*1.1 {
				log.Printf("ERROR in student VideoIn %f vs %f\n", studentBrowser.VideoIn, teacherBrowser.VideoOut)
				os.Exit(1)
			}
			if studentBrowser.VideoOut < teacherBrowser.VideoIn*0.9 || studentBrowser.VideoOut > teacherBrowser.VideoIn*1.1 {
				log.Printf("ERROR in student VideoOut %f vs %f\n", studentBrowser.VideoOut, teacherBrowser.VideoIn)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}
