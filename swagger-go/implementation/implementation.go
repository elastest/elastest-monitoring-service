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

// WriteResponse to the client
func (o SubscribeOk) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

    rw.WriteHeader(200)
	if err := producer.Produce(rw, o.subscriptionId); err != nil {
		panic(err) // let the recovery middleware deal with this
    }
}

func (endpoint LocalESEndpoint) getInjectableString(subId string) string {
	template := `# SUBID %s
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
	template := `# SUBID %s
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

func subscribeEndpoint(endpoint IEndpoint) middleware.Responder {
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
