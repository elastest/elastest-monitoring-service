package impl

import(
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    parsercommon "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
)

type StreamExprToStriverVisitor struct {
    momvisitor *MoMToStriverVisitor
    streamname striverdt.StreamName
}

func (visitor StreamExprToStriverVisitor) VisitAggregatorExpr(aggexp session.AggregatorExpr) {
    switch aggexp.Operation {
    case "avg":
        makeAvgOutStream(aggexp.Stream, aggexp.Session, visitor.streamname, visitor.momvisitor)
    default:
        panic("Operation "+aggexp.Operation+" not implemented")
    }
}

func (visitor StreamExprToStriverVisitor) VisitIfThenExpr(ifthen session.IfThenExpr) {
    mysignalname := visitor.streamname
    thensignalname := "ifthen_thenstream::"+mysignalname
    visitor.streamname = thensignalname
    ifthen.Then.Accept(visitor)
    makeIfThenStream(ifthen.If, thensignalname, mysignalname, visitor.momvisitor)
}
func (visitor StreamExprToStriverVisitor) VisitIfThenElseExpr(session.IfThenElseExpr) {
    panic("not implemented")
}
func (visitor StreamExprToStriverVisitor) VisitPredExpr(predExp session.PredExpr) {
    makePredicateStream(predExp.Pred, visitor.streamname, visitor.momvisitor)
}
func (visitor StreamExprToStriverVisitor) VisitStringPathExpr(session.StringPathExpr) {
    panic("not implemented")
}
func (visitor StreamExprToStriverVisitor) VisitNumExprStream(nes session.NumExprStream) {
    nes.NumExpr.Accept(NumExprToStriverVisitor{visitor})
}

func makeAvgOutStream(inSignalName, sessionSignalName, outSignalName striverdt.StreamName, visitor *MoMToStriverVisitor) {
    condCounterName := "condcounter::"+outSignalName
    condCounterFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        cond := args[0]
        if !cond.IsSet || !cond.Val.(striverdt.EvPayload).Val.(bool) {
            return striverdt.NothingPayload
        }
        myprev := args[1]
        prev := 0
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(int)
        }
        return striverdt.Some(prev+1)
    }
    condCounterVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, sessionSignalName, []striverdt.Event{}},
        &striverdt.PrevValNode{striverdt.TNode{}, condCounterName, []striverdt.Event{}},
    }, condCounterFun}
    condCounterStream := striverdt.OutStream{condCounterName, striverdt.SrcTickerNode{inSignalName}, condCounterVal}
    visitor.OutStreams = append(visitor.OutStreams, condCounterStream)

    condAvgFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        cond := args[0]
        if !cond.IsSet || !cond.Val.(striverdt.EvPayload).Val.(bool) {
            return striverdt.NothingPayload
        }
        myprev := args[1]
        currentval := args[2].Val.(striverdt.EvPayload).Val.(float32)
        kplusone := float32(args[3].Val.(striverdt.EvPayload).Val.(int))
        prev := float32(0.0)
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(float32)
        }
        res := (prev*(kplusone-1)+currentval)/kplusone
        return striverdt.Some(res)
    }
    condAvgVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, sessionSignalName, []striverdt.Event{}},
        &striverdt.PrevValNode{striverdt.TNode{}, outSignalName, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, inSignalName, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, condCounterName, []striverdt.Event{}},
    }, condAvgFun}
    condAvgStream := striverdt.OutStream{outSignalName, striverdt.SrcTickerNode{inSignalName}, condAvgVal}
    visitor.OutStreams = append(visitor.OutStreams, condAvgStream)
}

func makeIfThenStream(ifpred parsercommon.Predicate, thensignalname, mysignalname striverdt.StreamName, visitor *MoMToStriverVisitor) {
    signalNamesVisitor := SignalNamesFromPredicateVisitor{[]striverdt.StreamName{}}
    ifpred.Accept(&signalNamesVisitor)
    signalNames := signalNamesVisitor.SNames
    condFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        rawevent := args[0]
        if !rawevent.IsSet {
            panic("No raw event!")
        }
        theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
        then := args[1]
        argsMap := map[striverdt.StreamName]interface{}{}
        for i,sname := range signalNames {
            argsMap[sname] = args[i+2].Val.(striverdt.EvPayload).Val
        }
        theEvalVisitor := EvalVisitor{false, theEvent, argsMap}
        ifpred.Accept(&theEvalVisitor)
        if theEvalVisitor.Result && then.IsSet {
            return then.Val.(striverdt.EvPayload)
        } else {
            return striverdt.NothingPayload
        }
    }
    argSignals := []striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, thensignalname, []striverdt.Event{}},
    }
    for _,sname := range signalNames {
        argSignals = append(argSignals,
        &striverdt.PrevEqValNode{striverdt.TNode{}, sname, []striverdt.Event{}})
    }
    condVal := striverdt.FuncNode{argSignals, condFun}
    condStream := striverdt.OutStream{mysignalname, striverdt.SrcTickerNode{visitor.InSignalName}, condVal} // TODO Check the source tick
    visitor.OutStreams = append(visitor.OutStreams, condStream)
}

func makePredicateStream(pred parsercommon.Predicate, mysignalname striverdt.StreamName, visitor *MoMToStriverVisitor) {
    signalNamesVisitor := SignalNamesFromPredicateVisitor{[]striverdt.StreamName{}}
    pred.Accept(&signalNamesVisitor)
    signalNames := signalNamesVisitor.SNames
    predFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        rawevent := args[0]
        if !rawevent.IsSet {
            panic("No raw event!")
        }
        theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
        argsMap := map[striverdt.StreamName]interface{}{}
        for i,sname := range signalNames {
            argsMap[sname] = args[i+1].Val.(striverdt.EvPayload).Val
        }
        theEvalVisitor := EvalVisitor{false, theEvent, argsMap}
        pred.Accept(&theEvalVisitor)
        return striverdt.Some(theEvalVisitor.Result)
    }

    argSignals := []striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
    }
    for _,sname := range signalNames {
        argSignals = append(argSignals,
        &striverdt.PrevEqValNode{striverdt.TNode{}, sname, []striverdt.Event{}})
    }


    predVal := striverdt.FuncNode{argSignals, predFun}
    predStream := striverdt.OutStream{mysignalname, striverdt.SrcTickerNode{visitor.InSignalName}, predVal}
    visitor.OutStreams = append(visitor.OutStreams, predStream)
}
