package implementation

import (
	"net/http"
	openapiruntime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"runtime"

	"encoding/json"
)

type HealthIsOk struct {
	Status string
}

// WriteResponse to the client
func (healthOk HealthIsOk) WriteResponse(rw http.ResponseWriter, producer openapiruntime.Producer) {
    rw.WriteHeader(200)
	healthokjson, _ := json.Marshal(healthOk)
	if err := producer.Produce(rw, string(healthokjson)); err != nil {
		panic(err) // let the recovery middleware deal with this
    }
}

type EnvironmentInfo struct { generalinfo map[string]interface{} }

// WriteResponse to the client
func (info EnvironmentInfo) WriteResponse(rw http.ResponseWriter, producer openapiruntime.Producer) {
    rw.WriteHeader(200)
	envjson, _ := json.Marshal(info.generalinfo)
	if err := producer.Produce(rw, string(envjson)); err != nil {
		panic(err) // let the recovery middleware deal with this
    }
}

var healthok HealthIsOk = HealthIsOk{"up"}

func GetHealth() middleware.Responder {
	return healthok
}

func GetEnvironment() middleware.Responder {
	// gather info here!
	var info map[string]interface{} = map[string]interface{} {
		"OS" : runtime.GOOS,
		"Architecture" : runtime.GOARCH,
	}
	return EnvironmentInfo {info}
}
