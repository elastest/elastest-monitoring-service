package eventproc

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    pb "github.com/elastest/elastest-monitoring-service/protobuf"
    "math/rand"
    "strconv"
    "encoding/json"
    "fmt"
	"github.com/elastest/elastest-monitoring-service/go_EMS/signals"
	"github.com/elastest/elastest-monitoring-service/go_EMS/moms"
)

var momengines map[int]dt.MoMEngine01 = make(map[int]dt.MoMEngine01)

func DeploySignals01(signalsDef string) *pb.MomPostReply {
    fmt.Println("Deploying signal defs: ", signalsDef)
    var sds []dt.SignalsDefinitions
    if err := json.Unmarshal([]byte(signalsDef), &sds); err != nil {
        fmt.Println("It was an invalid array of signal definitions because the JSON was malformed: ", err.Error())
        return &pb.MomPostReply{Deploymenterror:err.Error(), Momid:""}
    }
    realdefs := make([]signals.SignalDefinition,len(sds))
    for i,sdef := range sds {
        switch(sdef.Type) {
        case "sampled":
            var realsampled signals.SampledSignalDefinition
            if err := json.Unmarshal([]byte(sdef.Def), &realsampled); err != nil {
                fmt.Println("It was an invalid real sampled signal definition because the JSON was malformed: ", err.Error())
                return &pb.MomPostReply{Deploymenterror:err.Error(), Momid:""}
            }
            realdefs[i] = realsampled
        default:
            fmt.Println("It was an invalid signal definition because the type was not recognized: ", sdef.Type)
            return &pb.MomPostReply{Deploymenterror:"Unrecognized signal type", Momid:""}
        }
    }
    momid := rand.Int()
    momengines[momid] = moms.StartEngine(realdefs)
    fmt.Println("with momid: ", momid)
    return &pb.MomPostReply{Deploymenterror:"", Momid:strconv.Itoa(momid)}
}

func UndeploySignals01(momid int) *pb.MomPostReply {
    if ek,ok := momengines[momid]; ok {
        delete(momengines, momid)
        ek.Kill()
        return &pb.MomPostReply{Deploymenterror:"", Momid:strconv.Itoa(momid)} // TODO change
    }
    return &pb.MomPostReply{Deploymenterror:"No such id", Momid:"momid"} // change
}

func ProcessEvent(evt dt.Event) {
    for _,engine := range momengines {
        samplers := engine.Samplers
        for _,sampler := range samplers {
            signals.SamplerProcessEvent(sampler, evt)
        }
    }
}

