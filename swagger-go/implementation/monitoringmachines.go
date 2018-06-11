package implementation

import (
	"net/http"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

    "log"
    "time"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
	"../restapi/operations/monitoring_machine"
    pb "github.com/elastest/elastest-monitoring-service/protobuf"
)

type MomPostReply pb.MomPostReply

// WriteResponse to the client
func (o MomPostReply) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
    var ret string
    if len(o.Deploymenterror) == 0 {
        rw.WriteHeader(200)
        ret = o.Momid
    } else {
        rw.WriteHeader(406)
        ret = o.Deploymenterror
    }
	if err := producer.Produce(rw, ret); err != nil {
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
	/*if _, err := os.Stat("/usr/share/logstash/pipes/swagpipe"); os.IsNotExist(err) {
		return DeployNotAllowed{}
	}
	definition := params.Mom.Definition
	momType := params.Mom.MomType
    var dasmap map[string]interface{} = nil
	if err := json.Unmarshal([]byte(definition), &dasmap); err != nil {
		fmt.Println("No JSON. Error: " + err.Error())
	} else {
		dasmap["momType"] = momType
		dasmap["deploymentId"] = deploymentId
	}
	newJSON, _ := json.Marshal(dasmap)
	file, err := os.Create("/usr/share/logstash/pipes/swagpipe")
    if err != nil {
        panic(err)
    }
	fmt.Fprintln(file, string(newJSON))*/
    req := pb.MomPostRequest{Momtype:params.Version, Momdefinition:*params.Mom}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
        // return error instead
    }
    defer conn.Close()
    c := pb.NewEngineClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.PostMoM(ctx, &req)
    if err != nil {
        log.Fatalf("could not greet: %v", err)
        // return error instead
    }
	return MomPostReply(*r)
}
