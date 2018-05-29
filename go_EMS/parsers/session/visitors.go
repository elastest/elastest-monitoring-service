package session

import(
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type MoMToStriverVisitor struct {
    OutStreams []striverdt.OutStream
    InSignalName striverdt.StreamName
}

func (visitor *MoMToStriverVisitor) VisitStream(stream Stream) {
    stream.Expr.Accept(StreamExprToStriverVisitor{visitor, stream.Name})
}
func (visitor *MoMToStriverVisitor) VisitTrigger(trigger Trigger) {
    makeTriggerStreams(visitor, trigger.Pred, trigger.Action)
}
func (visitor *MoMToStriverVisitor) VisitFilter(Filter) {
    panic("No filters allowed!")
}
func (visitor *MoMToStriverVisitor) VisitSession(Session) {
    panic("No session allowed!")
}
func (visitor *MoMToStriverVisitor) VisitPredicateDecl(PredicateDecl) {
    panic("No PredicateDecl allowed!")
}
