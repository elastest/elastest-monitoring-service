package implementation

import (
	"net/http"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"../restapi/operations/monitoring_machine"
	"fmt"
	"math/rand"
	"strconv"
	"encoding/json"
	"os"
)

type DeployOk struct {

    // In: body
	deploymentId string
}

// WriteResponse to the client
func (o DeployOk) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

    rw.WriteHeader(200)
	if err := producer.Produce(rw, o.deploymentId); err != nil {
		panic(err) // let the recovery middleware deal with this
    }
}

type DeployNotAllowed struct { }

// WriteResponse to the client
func (o DeployNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

    rw.WriteHeader(405)
	if err := producer.Produce(rw, "Deployment is not allowed in this version of the ElasTest Monitoring Service"); err != nil {
		panic(err) // let the recovery middleware deal with this
    }
}

func PostMOM(params monitoring_machine.PostMoMParams) middleware.Responder {
	if _, err := os.Stat("/usr/share/logstash/pipes/swagpipe"); os.IsNotExist(err) {
		return DeployNotAllowed{}
	}
	definition := params.Mom.Definition
	momType := params.Mom.MomType
    var dasmap map[string]interface{} = nil
	deploymentId := strconv.Itoa(rand.Int())
	if err := json.Unmarshal([]byte(definition), &dasmap); err != nil {
		fmt.Println("No JSON. Error: " + err.Error())
	} else {
		dasmap["momType"] = momType
		dasmap["deploymentId"] = deploymentId
	}
	newJSON, _ := json.Marshal(dasmap)
	fmt.Println("JSON: " + string(newJSON))
	file, err := os.Create("/usr/share/logstash/pipes/swagpipe")
    if err != nil {
        panic(err)
    }
	fmt.Fprintln(file, string(newJSON))
	return DeployOk{deploymentId}
}
