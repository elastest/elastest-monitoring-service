package eventproc

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    pb "github.com/elastest/elastest-monitoring-service/protobuf"
    "math/rand"
    "strconv"
    "encoding/json"
    "fmt"
)

var tagConditions map[int]dt.TagCondition = make(map[int]dt.TagCondition)

func DeployTaggerv01(taggerDef string) *pb.MomPostReply {
    fmt.Println("Deploying def: ", taggerDef)
    var td dt.TaggerDefinition
    if err := json.Unmarshal([]byte(taggerDef), &td); err != nil {
        fmt.Println("It was an invalid definition because the JSON was malformed: ", err.Error())
        return &pb.MomPostReply{Deploymenterror:err.Error(), Momid:""}
    }
    if !validTD(td) {
        fmt.Println("Missing or empty filter or out channel")
        return &pb.MomPostReply{Deploymenterror:"Missing or empty filter or out channel", Momid:""}
    }
    var thejson map[string]interface{}
    filterbytes := []byte(td.Filter)
    if err := json.Unmarshal(filterbytes, &thejson); err != nil {
        errtext := "Filter is not a JSON. Error: "+ err.Error()
        return &pb.MomPostReply{Deploymenterror:errtext, Momid:""}
    }
    tagNode, err := getNodeFromFilter(thejson)
    if err != nil {
        fmt.Println("Error in filter definition: ", err.Error())
        return &pb.MomPostReply{Deploymenterror:err.Error(), Momid:""}
    }
    momid := rand.Int()
    fmt.Println("with momid: ", momid)
    tagConditions[momid] = dt.TagCondition{
        sets.SetFromList(td.InChannels),
        tagNode.Eval,
        dt.Channel(td.OutChannel),
    }
    return &pb.MomPostReply{Deploymenterror:"", Momid:strconv.Itoa(momid)}
}

func TagEvent(ev *dt.Event) {
    var checkConditions []dt.TagCondition
    // filter out unsatisfiable conditions
    for _,tc := range tagConditions {
        if tc.EventCondition(ev.Payload) {
            checkConditions = append(checkConditions, tc)
        }
    }
    checkChans := (*ev).Channels

    dirty:=true
    for (dirty) {
        dirty=false
        newconds := checkConditions[:0]
        var nextCheckChans dt.ChannelSet = make(dt.ChannelSet)
        for _,cond := range checkConditions {
            if !(sets.SetIn(cond.OutChannel, checkChans)) { // if it's not tagged yet
                cond.InChannels = sets.SetMinus(cond.InChannels, checkChans)
                if (sets.SetIsEmpty(cond.InChannels)) { // triggered
                    dirty=true
                    nextCheckChans = sets.SetAdd(nextCheckChans, cond.OutChannel)
                    checkChans = sets.SetAdd(checkChans, cond.OutChannel)
                } else {
                    newconds = append(newconds, cond) // check on next iteration
                }
            }
        }
        (*ev).Channels = sets.SetUnion((*ev).Channels, nextCheckChans) // add new triggered channels
        checkChans = nextCheckChans
        checkConditions = newconds
    }

}

func validTD(td dt.TaggerDefinition) bool {
    return (td.Filter != "" && td.OutChannel != "")
}
