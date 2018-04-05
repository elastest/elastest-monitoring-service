package eventproc

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
)

var tagConditions []dt.TagCondition = []dt.TagCondition {
    dt.TagCondition{
        dt.ChannelSet(map[dt.Channel]interface{}{dt.Channel("a"):nil,}),
        func(ev dt.Event) bool {return true},
        dt.Channel("C"),
    },
    dt.TagCondition{
        dt.ChannelSet(map[dt.Channel]interface{}{dt.Channel("C"):nil,}),
        func(ev dt.Event) bool {return true},
        dt.Channel("D"),
    },
    dt.TagCondition{
        dt.ChannelSet(map[dt.Channel]interface{}{dt.Channel("D"):nil,dt.Channel("b"):nil}),
        func(ev dt.Event) bool {return true},
        dt.Channel("E"),
    },
}

func TagEvent(ev *dt.Event) {
    var checkConditions []dt.TagCondition
    // filter out unsatisfiable conditions
    for _,tc := range tagConditions {
        if tc.EventCondition(*ev) {
            checkConditions = append(checkConditions, tc)
        }
    }

    checkChans := (*ev).Channels
    for (len(checkConditions) > 0 && len(checkChans) > 0) {

        newconds := checkConditions[:0]
        var nextCheckChans dt.ChannelSet = make(dt.ChannelSet)
        for _,cond := range checkConditions {
            if !(sets.SetIn(cond.OutChannel, checkChans)) { // if it's not tagged yet
                cond.InChannels = sets.SetMinus(cond.InChannels, checkChans)
                if (sets.SetIsEmpty(cond.InChannels)) { // triggered
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
