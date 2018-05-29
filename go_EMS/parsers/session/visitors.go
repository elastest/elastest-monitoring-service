package session

import(
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    parsercommon "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	"github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
)

type MoMToStriverVisitor struct {
    OutStreams []striverdt.OutStream
    InSignalName striverdt.StreamName
}

func (visitor *MoMToStriverVisitor) VisitStream(stream Stream) {
    stream.Expr.Accept(StreamExprToStriverVisitor{visitor, stream.Name})
}
func (visitor *MoMToStriverVisitor) VisitFilter(Filter) {
    panic("No filters allowed!")
}
func (visitor *MoMToStriverVisitor) VisitSession(Session) {
    panic("No session allowed!")
}
func (visitor *MoMToStriverVisitor) VisitTrigger(Trigger) {
    panic("No trigger allowed!")
}
func (visitor *MoMToStriverVisitor) VisitPredicateDecl(PredicateDecl) {
    panic("No PredicateDecl allowed!")
}

type StreamExprToStriverVisitor struct {
    momvisitor *MoMToStriverVisitor
    streamname striverdt.StreamName
}

func (visitor StreamExprToStriverVisitor) visitAggregatorExpr(aggexp AggregatorExpr) {
    switch aggexp.Operation {
    case "avg":
        makeAvgOutStream(aggexp.Stream, aggexp.Session, visitor.streamname, visitor.momvisitor)
    default:
        panic("Operation "+aggexp.Operation+" not implemented")
    }
}

func (visitor StreamExprToStriverVisitor) visitIfThenExpr(ifthen IfThenExpr) {
    mysignalname := visitor.streamname
    thensignalname := "ifthen_thenstream::"+mysignalname
    visitor.streamname = thensignalname
    ifthen.Then.Accept(visitor)
    makeIfThenStream(ifthen.If, thensignalname, mysignalname, visitor.momvisitor)
}
func (visitor StreamExprToStriverVisitor) visitIntPathExpr(ipathexpr IntPathExpr) {
    makeIntPathStream(ipathexpr.Path, visitor.streamname, visitor.momvisitor)
}
func (visitor StreamExprToStriverVisitor) visitIfThenElseExpr(IfThenElseExpr) {
    panic("not implemented")
}
func (visitor StreamExprToStriverVisitor) visitNumExpr(NumExpr) {
    panic("not implemented")
}
func (visitor StreamExprToStriverVisitor) visitPredExpr(PredExpr) {
    panic("not implemented")
}
func (visitor StreamExprToStriverVisitor) visitStringPathExpr(StringPathExpr) {
    panic("not implemented")
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
        currentval := args[2].Val.(striverdt.EvPayload).Val.(float64)
        kplusone := float64(args[3].Val.(striverdt.EvPayload).Val.(int))
        prev := 0.0
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(float64)
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
    condFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        rawevent := args[0]
        then := args[1]
        if !rawevent.IsSet {
            panic("No raw event!")
        }
        theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
        theEvalVisitor := parsercommon.EvalVisitor{false, theEvent}
        ifpred.Accept(&theEvalVisitor)
        if theEvalVisitor.Result {
            return then
        } else {
            return striverdt.NothingPayload
        }
    }
    condVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
        &striverdt.PrevEqValNode{striverdt.TNode{}, thensignalname, []striverdt.Event{}},
    }, condFun}
    condStream := striverdt.OutStream{mysignalname, striverdt.SrcTickerNode{visitor.InSignalName}, condVal} // TODO Check the source tick
    visitor.OutStreams = append(visitor.OutStreams, condStream)
}

func makeIntPathStream(path dt.JSONPath, mysignalname striverdt.StreamName, visitor *MoMToStriverVisitor) {
    extractFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        rawevent := args[0]
        if !rawevent.IsSet {
            panic("No raw event!")
        }
        theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
        valif, err := jsonrw.ExtractFromMap(theEvent.Payload, path)
        if err != nil {
            panic("No path found!")
        }
        return striverdt.Some(valif.(float64))
    }
    extractVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
    }, extractFun}
    extractStream := striverdt.OutStream{mysignalname, striverdt.SrcTickerNode{visitor.InSignalName}, extractVal}
    visitor.OutStreams = append(visitor.OutStreams, extractStream)
}
