package common
// import "fmt"

type NameToExprStreamVisitor struct {
  Namepar string
  Expr FloatLiteralExpr
  ReturnExpr StreamExpr
  ReturnPred Predicate
  ReturnNumExpr NumExpr
}

func (visitor *NameToExprStreamVisitor) VisitAggregatorExpr(aggexp AggregatorExpr) {
  visitor.ReturnExpr = aggexp
}

func (visitor *NameToExprStreamVisitor) VisitIfThenExpr(ifthen IfThenExpr) {
  //ifthen.If.AcceptPred(visitor)
  //ifthen.If = visitor.ReturnPred
  ifthen.Then.Accept(visitor)
  ifthen.Then = visitor.ReturnExpr
  visitor.ReturnExpr = ifthen
}
func (visitor *NameToExprStreamVisitor) VisitIfThenElseExpr(IfThenElseExpr) {
    panic("not implemented")
}
func (visitor *NameToExprStreamVisitor) VisitPredExpr(predExp PredExpr) {
  //predExp.Pred.AcceptPred(visitor)
  //predExp.Pred = visitor.ReturnPred
  visitor.ReturnExpr = predExp
}
func (visitor *NameToExprStreamVisitor) VisitStringPathExpr(exp StringPathExpr) {
  visitor.ReturnExpr = exp
}
func (visitor *NameToExprStreamVisitor) VisitStreamNumExpr(numExp StreamNumExpr) {
  //numExp.Expr.AcceptNum(visitor)
  //numExp.Expr = visitor.ReturnNumExpr
  visitor.ReturnExpr = numExp
}
func (visitor *NameToExprStreamVisitor) VisitStreamNameExpr(nes StreamNameExpr) {
	if nes.Stream == visitor.Namepar {
    newexpr := StreamNumExpr{visitor.Expr}
    visitor.ReturnExpr = newexpr
  } else {
    visitor.ReturnExpr = nes
  }
}

