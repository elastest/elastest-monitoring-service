// Example taken from https://golangcode.com/download-a-file-from-a-url/
package main

import (
  "encoding/json"
  "net/http"
  "time"
  "os"
  "os/exec"
  "fmt"
  "github.com/gorilla/websocket"
  "log"
  "net/url"
  "io/ioutil"
  "io"
  "bytes"
)

type event struct {
  Channels      []string `json:"channels,omitempty"`
  Value         bool `json:"value,omitempty"`
}

//Net net structure
type Net struct {
  TxBytesPs float64 `json:"txBytes_ps,omitempty"`
}

func postMoMs(ems string) {
  // curl -H "Content-Type:text/plain"  --data-binary @stampers.txt http://${ET_EMS_LSBEATS_HOST}:8888/stamper/tag0.1
  dat, err := ioutil.ReadFile("stampers.txt")
  if err != nil {
    panic(err)
  }
  req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:8888/stamper/tag0.1", ems), bytes.NewBuffer(dat))
  if err != nil {
    panic(err)
  }
  req.Header.Set("Content-Type", "text/plain")
  client := &http.Client{}
  resp, err := client.Do(req)
  f, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }
  log.Println(string(f))
  if err != nil {
    panic(err)
  }

  // curl -H "Content-Type:text/plain" --data-binary @sessiondef.txt http://${ET_EMS_LSBEATS_HOST}:8888/MonitoringMachine/signals0.1
  dat, err = ioutil.ReadFile("sessiondef.txt")
  if err != nil {
    panic(err)
  }
  req, err = http.NewRequest("POST", fmt.Sprintf("http://%s:8888/MonitoringMachine/signals0.1", ems), bytes.NewBuffer(dat))
  if err != nil {
    panic(err)
  }
  req.Header.Set("Content-Type", "text/plain")
  client = &http.Client{}
  resp, err = client.Do(req)
  f, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }
  log.Println(string(f))
  if err != nil {
    panic(err)
  }
}

func main() {
  ems := os.Getenv("ET_EMS_LSBEATS_HOST")
  postMoMs(ems)
  go stress()
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
    log.Printf("Received event:[%s]\n", string(input))
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

func stress() {
  host:=os.Getenv("ET_SUT_HOST")
  fileUrl := "http://"+host+"/sparse"

  fmt.Println("STARTING FIRST DOWNLOAD")
  go DownloadFile(fileUrl)
  // curl -sS http://${ET_SUT_HOST}/sparse >/dev/null & 
  time.Sleep(60 * time.Second)
  fmt.Println("STARTING SECOND DOWNLOAD")
  go DownloadFile(fileUrl)
  time.Sleep(60 * time.Second)
  fmt.Println("FINISHING SECOND DOWNLOAD")
  time.Sleep(5 * time.Second)
  fmt.Println("FINISHING TEST")
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(url string) {

  //fmt.Println(`Executing: "sh", "-c", "curl -sS "`+url+`">/dev/null"`)
  cmd := exec.Command("sh", "-c", "curl -sS --limit-rate 10M "+url+">/dev/null")
  err := cmd.Run()
  return

  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  // Write the body to file
  fmt.Println("copying file")
  _, err = io.Copy(os.Stderr, resp.Body)
  fmt.Println("finished copying file")
  if err != nil {
    panic(err)
  }
}
