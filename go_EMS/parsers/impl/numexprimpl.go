package impl

import(
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	"github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
)

type NumExprToStriverVisitor struct {
    streamvisitor StreamExprToStriverVisitor
}

func (visitor NumExprToStriverVisitor) VisitNumPathExpr(ipathexpr common.NumPathExpr) {
    makeNumPathStream(ipathexpr.Path, visitor.streamvisitor.streamname, visitor.streamvisitor.momvisitor)
}
func (visitor NumExprToStriverVisitor) VisitFloatLiteralExpr(common.FloatLiteralExpr) {
    panic("not implemented")
}
func (visitor NumExprToStriverVisitor) VisitIntLiteralExpr(common.IntLiteralExpr) {
    panic("not implemented")
}
func (visitor NumExprToStriverVisitor) VisitStreamNameExpr(common.StreamNameExpr) {
    panic("not implemented")
}
func (visitor NumExprToStriverVisitor) VisitNumMulExpr(common.NumMulExpr) {
    panic("not implemented")
}
func (visitor NumExprToStriverVisitor) VisitNumDivExpr(common.NumDivExpr) {
    panic("not implemented")
}
func (visitor NumExprToStriverVisitor) VisitNumPlusExpr(common.NumPlusExpr) {
    panic("not implemented")
}
func (visitor NumExprToStriverVisitor) VisitNumMinusExpr(common.NumMinusExpr) {
    panic("not implemented")
}

func makeNumPathStream(path dt.JSONPath, mysignalname striverdt.StreamName, visitor *MoMToStriverVisitor) {
    extractFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        rawevent := args[0]
        if !rawevent.IsSet {
            panic("No raw event!")
        }
        theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
        valif, err := jsonrw.ExtractFromMap(theEvent.Payload, path)
        if err != nil {
            /* This can happen: the stream might be guarded by an if statement upper in the AST */
            return striverdt.NothingPayload
        }
        return striverdt.Some(float32(valif.(float64)))
    }
    extractVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
    }, extractFun}
    extractStream := striverdt.OutStream{mysignalname, striverdt.SrcTickerNode{visitor.InSignalName}, extractVal}
    visitor.OutStreams = append(visitor.OutStreams, extractStream)
}
