package impl

import(
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
    parsercommon "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
)

type MoMToStriverVisitor struct {
    OutStreams []striverdt.OutStream
    InSignalName striverdt.StreamName
    Preds map[string]parsercommon.Predicate
}

func (visitor *MoMToStriverVisitor) VisitStream(stream session.Stream) {
    stream.Expr.Accept(StreamExprToStriverVisitor{visitor, stream.Name})
}
func (visitor *MoMToStriverVisitor) VisitTrigger(trigger session.Trigger) {
    makeTriggerStreams(visitor, trigger.Pred, trigger.Action)
}
func (visitor *MoMToStriverVisitor) VisitFilter(session.Filter) {
    panic("No filters allowed!")
}
func (visitor *MoMToStriverVisitor) VisitSession(session.Session) {
    panic("No session allowed!")
}
func (visitor *MoMToStriverVisitor) VisitPredicateDecl(pd session.PredicateDecl) {
    visitor.Preds[pd.Name] = pd.Pred
}
