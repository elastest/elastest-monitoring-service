package impl

import(
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    parsercommon "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
)

var TRIGGER_PREFIX striverdt.StreamName = "trigger::"

func makeTriggerStreams(visitor *MoMToStriverVisitor, predicate parsercommon.Predicate, triggerAction session.EmitAction) {
	//StreamName striverdt.StreamName
	//TagName    common.Tag

    inputSignalName := triggerAction.StreamName
    outputSignalName := TRIGGER_PREFIX + inputSignalName
    condFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        rawevent := args[0]
        then := args[1]
        if !rawevent.IsSet {
            panic("No raw event!")
        }
        theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
        theEvalVisitor := EvalVisitor{false, theEvent}
        predicate.Accept(&theEvalVisitor)
        if theEvalVisitor.Result && then.IsSet {
            theVal := then.Val.(striverdt.EvPayload).Val
            return striverdt.Some(dt.Event{
                        dt.ChannelSet{triggerAction.TagName.Tag:nil},
                        map[string]interface{}{"value":theVal},
                        theEvent.Timestamp})
        } else {
            return striverdt.NothingPayload
        }
    }
    condVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, inputSignalName, []striverdt.Event{}},
    }, condFun}
    condStream := striverdt.OutStream{outputSignalName, striverdt.SrcTickerNode{visitor.InSignalName}, condVal} // TODO Check the source tick
    visitor.OutStreams = append(visitor.OutStreams, condStream)
}
