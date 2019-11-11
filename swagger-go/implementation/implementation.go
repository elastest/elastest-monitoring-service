package implementation

import (
  "net/http"
  runtime "github.com/go-openapi/runtime"
  middleware "github.com/go-openapi/runtime/middleware"
  "io/ioutil"

  "../restapi/operations/subscribers"
	"../restapi/operations/offline"
  "../models"
  "fmt"
  "math/rand"
  "strconv"
  "os"
  "regexp"
  "strings"
  "os/exec"
)

type LocalESEndpoint models.ESEndpoint
type LocalRMQEndpoint models.RMQEndpoint

type IEndpoint interface {
  getInjectableString(subId string) string
}

type SubscribeOk struct {

  // In: body
  subscriptionId string
}


type UnSubscribeOk struct {
}

type SubscribeManyOk struct {

  // In: body
  subscriptionIds []string
}

type OfflineResponder struct { events string }

// WriteResponse to the client
func (or OfflineResponder) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
  rw.WriteHeader(200)
	if err := producer.Produce(rw, or.events); err != nil {
		panic(err) // let the recovery middleware deal with this
  }
}

// WriteResponse to the client
func (o SubscribeOk) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

  rw.WriteHeader(200)
  if err := producer.Produce(rw, o.subscriptionId); err != nil {
    panic(err) // let the recovery middleware deal with this
  }
}

func (o UnSubscribeOk) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

  rw.WriteHeader(200)
  if err := producer.Produce(rw, "ok"); err != nil {
    panic(err) // let the recovery middleware deal with this
  }
}

func (o SubscribeManyOk) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

  rw.WriteHeader(200)
  if err := producer.Produce(rw, o.subscriptionIds); err != nil {
    panic(err) // let the recovery middleware deal with this
  }
}

func (endpoint LocalESEndpoint) getInjectableString(subId string) string {
  template := `
  # SUBID %s
  if "%s" in [channels] {
    elasticsearch {
      hosts => ["%s:%v"]
      user => %s
      password => %s
    }
  }
  # ENDOF %s
}`
return fmt.Sprintf(template, subId, endpoint.Channel, endpoint.IP, endpoint.Port, endpoint.User, endpoint.Password, subId)
}

func (endpoint LocalRMQEndpoint) getInjectableString(subId string) string {
  template := `
  # SUBID %s
  if "%s" in [channels] {
    rabbitmq {
      key => "%s"
      exchange => "%s"
      exchange_type => "%s"
      user => "%s"
      password => "%s"
      host => "%s"
      port => %v
      durable => true
      persistent => true
    }
  }
  # ENDOF %s
}`
return fmt.Sprintf(template, subId, endpoint.Channel, endpoint.Key, endpoint.Exchange, endpoint.ExchangeType, endpoint.User, endpoint.Password, endpoint.IP, endpoint.Port, subId)
}

func injectNewOutput(injStr string) {
  // open input file
  conffile := "/usr/share/logstash/pipeline/outlogstash.conf"
  file, err := os.OpenFile(conffile, os.O_RDWR, 0755)
  if err != nil {
    panic(err)
  }
  // close file on exit and check for its returned error
  defer func() {
    if err := file.Close(); err != nil {
      panic(err)
    }
  }()

  finfo, err := file.Stat()
  if err != nil {
    panic(err)
  }

  fsize := finfo.Size()
  //fmt.Printf("The file is %d bytes long", fsize)

  if _, err := file.WriteAt([]byte(injStr), fsize-2); err != nil {
    panic(err)
  }
}

func subscribeEndpoint(endpoint IEndpoint) SubscribeOk {
  subId := strconv.Itoa(rand.Int())
  injStr := endpoint.getInjectableString(subId)
  injectNewOutput(injStr)
  //fmt.Println(injStr)
  return SubscribeOk{subId}
}


func SubscribeES(params subscribers.SubscribeElasticSearchParams) middleware.Responder {
  endpoint := LocalESEndpoint(*params.Endpoint)
  return subscribeEndpoint(endpoint)
}

func SubscribeRMQ(params subscribers.SubscribeRabbitMQParams) middleware.Responder {
  endpoint := LocalRMQEndpoint(*params.Endpoint)
  return subscribeEndpoint(endpoint)
}

func subscribeDefaultEndpoint(ep string, channel string) string {
  edm_es := os.Getenv("ET_EDM_ELASTICSEARCH_API")

  // Sanity check
  r, _ := regexp.Compile("http://[^:]+:[0-9]+/")
  if !(r.MatchString(edm_es)) {
    return "invalid_EDM_ES"
  }


  edm_es = edm_es[7:len(edm_es)-1]
  i := strings.Index(edm_es, ":")
  ip := edm_es[:i]
  port, _ := strconv.ParseInt(edm_es[i+1:], 10, 64)

  switch ep {
  case "persistence","dashboard":
    return subscribeEndpoint(LocalESEndpoint(models.ESEndpoint{channel, ip, "changeme",  port, "elastic"})).subscriptionId
  default:
    return "unknown default endpoint"
  }
}

func SubscribeElastestEndpoint(params subscribers.SubscribeElastestEndpointsParams) middleware.Responder {
  channel := params.Endpoints.Channel
  endpoints := params.Endpoints.Endpoints

  subids := make([]string, len(endpoints))
  for i, v := range endpoints {
    subids[i] = subscribeDefaultEndpoint(v, channel)
  }
  return SubscribeManyOk{subids}
}

func UnsubscribeHandler(params subscribers.UnsubscribeParams) middleware.Responder {
  // sed -i '/SUBID 5577006791947779410/,/ENDOF 5577006791947779410/d' testfile.txt
  cmd := "sed"
  args := []string{"-i", "/SUBID "+params.SubID+"/,/ENDOF "+params.SubID+"/d", "/usr/share/logstash/pipeline/outlogstash.conf"}
  if err := exec.Command(cmd, args...).Run(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  return UnSubscribeOk{}
}

func offline_pipe_reader(c chan string) {
  dat, err := ioutil.ReadFile("/usr/share/logstash/pipes/offlinepipe")
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  c <- string(dat)
  close(c)
}

func Offline(params offline.OfflineParams) middleware.Responder {
  c := make(chan string)
  go offline_pipe_reader(c)
  template := `
input {
  jdbc {
    jdbc_driver_class => "com.mysql.jdbc.Driver"
    jdbc_connection_string => "jdbc:mysql://172.25.0.4:3306/ETM"
    jdbc_user => "elastest"
    jdbc_password => "elastest"
    parameters => {"execid" => "%s"}
    statement => "select * from Trace where exec = :execid"
  }
}

output {
  file {
    path => "/usr/share/logstash/pipes/offlinepipe"
    codec => json_lines
  }
}`
  execid := params.ExecID
  inlscontent := fmt.Sprintf(template, execid)
  if err := ioutil.WriteFile("/usr/share/logstash/pipeline/inlogstash.conf", []byte(inlscontent), 0644); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  evs := <-c
  return OfflineResponder{"holis " + evs}
}
