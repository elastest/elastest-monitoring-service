package session

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type MoMToStriverVisitor struct {
    Samplers []dt.Sampler
    OutStreams []striverdt.OutStream
    InStreams []striverdt.InStream
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
    panic("No PredicateDeclallowed!")
}

type StreamExprToStriverVisitor struct {
    momvisitor *MoMToStriverVisitor
    streamname striverdt.StreamName
}

func (visitor StreamExprToStriverVisitor) visitAggregatorExpr(aggexp AggregatorExpr) {
    switch aggexp.Operation {
    case "avg":
        getAvgOutStream(aggexp.Stream, aggexp.Session, visitor.streamname, visitor.momvisitor)
    default:
        panic("Operation "+aggexp.Operation+" not implemented")
    }
}

func (visitor StreamExprToStriverVisitor) visitIfThenExpr(IfThenExpr) {
    panic("not implemented")
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

func getAvgOutStream(inSignalName, sessionSignalName, outSignalName striverdt.StreamName, visitor *MoMToStriverVisitor) {
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
        cpuval := args[2].Val.(striverdt.EvPayload).Val.(float64)
        kplusone := float64(args[3].Val.(striverdt.EvPayload).Val.(int))
        prev := 0.0
        if myprev.IsSet {
            prev = myprev.Val.(striverdt.EvPayload).Val.(float64)
        }
        res := (prev*(kplusone-1)+cpuval)/kplusone
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
