package implementation

import (
	"net/http"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"../restapi/operations/subscribers"
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
  elasticsearch {
    hosts => ["%s:%v"]
    user => %s
    password => %s
  }
# ENDOF %s
}`
  return fmt.Sprintf(template, subId, endpoint.IP, endpoint.Port, endpoint.User, endpoint.Password, subId)
}

func (endpoint LocalRMQEndpoint) getInjectableString(subId string) string {
	template := `
# SUBID %s
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
# ENDOF %s
}`
  return fmt.Sprintf(template, subId, endpoint.Key, endpoint.Exchange, endpoint.ExchangeType, endpoint.User, endpoint.Password, endpoint.IP, endpoint.Port, subId)
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


