package impl

import(
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    parsercommon "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
)

var TRIGGER_PREFIX striverdt.StreamName = "trigger::"

func makeTriggerStreams(visitor *MoMToStriverVisitor, predicate parsercommon.Predicate, triggerAction session.EmitAction) {
    inputSignalName := triggerAction.StreamName
    outputSignalName := TRIGGER_PREFIX + inputSignalName
    signalNamesVisitor := SignalNamesFromPredicateVisitor{visitor.Preds, []striverdt.StreamName{}}
    predicate.Accept(&signalNamesVisitor)
    signalNames := signalNamesVisitor.SNames
    predFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        rawevent := args[0]
        if !rawevent.IsSet {
            panic("No raw event!")
        }
        theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
        then := args[1]
        argsMap := map[striverdt.StreamName]interface{}{}
        for i,sname := range signalNames {
            if !args[i+2].IsSet {
                return striverdt.NothingPayload
            }
            argsMap[sname] = args[i+2].Val.(striverdt.EvPayload).Val
        }
        theEvalVisitor := EvalVisitor{false, theEvent, visitor.Preds, argsMap}
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

    argSignals := []striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, inputSignalName, []striverdt.Event{}},
    }
    for _,sname := range signalNames {
        argSignals = append(argSignals,
        &striverdt.PrevEqValNode{striverdt.TNode{}, sname, []striverdt.Event{}})
    }


    predVal := striverdt.FuncNode{argSignals, predFun}
    predStream := striverdt.OutStream{outputSignalName, striverdt.SrcTickerNode{visitor.InSignalName}, predVal}
    visitor.OutStreams = append(visitor.OutStreams, predStream)
}
