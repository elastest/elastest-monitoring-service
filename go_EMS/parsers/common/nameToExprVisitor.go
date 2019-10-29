package common
import (
  // "fmt"
  "regexp"
  "strconv"
)

type NameToExprStreamVisitor struct {
  Namepar string
  Expr FloatLiteralExpr
  ReturnExpr StreamExpr
  ReturnPred Predicate
  ReturnNumExpr NumExpr
  ReturnNumComparison NumComparison
}

func (visitor *NameToExprStreamVisitor) VisitAggregatorExpr(aggexp AggregatorExpr) {
  visitor.ReturnExpr = aggexp
}

func (visitor *NameToExprStreamVisitor) VisitIfThenExpr(ifthen IfThenExpr) {
  ifthen.If.AcceptPred(visitor)
  ifthen.If = visitor.ReturnPred
  ifthen.Then.Accept(visitor)
  ifthen.Then = visitor.ReturnExpr
  visitor.ReturnExpr = ifthen
}
func (visitor *NameToExprStreamVisitor) VisitIfThenElseExpr(IfThenElseExpr) {
    panic("not implemented")
}
func (visitor *NameToExprStreamVisitor) VisitPredExpr(predExp PredExpr) {
  predExp.Pred.AcceptPred(visitor)
  predExp.Pred = visitor.ReturnPred
  visitor.ReturnExpr = predExp
}
func (visitor *NameToExprStreamVisitor) VisitStringPathExpr(exp StringPathExpr) {
  visitor.ReturnExpr = exp
}
func (visitor *NameToExprStreamVisitor) VisitStreamNumExpr(numExp StreamNumExpr) {
  numExp.Expr.AcceptNum(visitor)
  numExp.Expr = visitor.ReturnNumExpr
  visitor.ReturnExpr = numExp
}
func (visitor *NameToExprStreamVisitor) VisitStreamNameExpr(nes StreamNameExpr) {
  daregex := regexp.MustCompile("_"+visitor.Namepar+"$")
	if nes.Stream == visitor.Namepar {
    newexpr := StreamNumExpr{visitor.Expr}
    visitor.ReturnExpr = newexpr
    visitor.ReturnNumExpr = visitor.Expr
    visitor.ReturnPred = nil
  } else if daregex.MatchString(nes.Stream)  {
    indexstr := strconv.Itoa(int(visitor.Expr.Num))
    newNes := StreamNameExpr{daregex.ReplaceAllString(nes.Stream, "_"+indexstr)}
    visitor.ReturnExpr = newNes
    visitor.ReturnNumExpr = newNes
    visitor.ReturnPred = newNes
  } else {
    visitor.ReturnExpr = nes
    visitor.ReturnNumExpr = nes
    visitor.ReturnPred = nes
  }
}

func (visitor *NameToExprStreamVisitor) VisitLastOfStreamNameExpr(nes LastOfStreamNameExpr) {
  daregex := regexp.MustCompile("_"+visitor.Namepar+"$")
	if nes.Stream == visitor.Namepar {
    newexpr := StreamNumExpr{visitor.Expr}
    visitor.ReturnExpr = newexpr
    visitor.ReturnNumExpr = visitor.Expr
    visitor.ReturnPred = nil
  } else if daregex.MatchString(nes.Stream) {
    indexstr := strconv.Itoa(int(visitor.Expr.Num))
    newNes := LastOfStreamNameExpr{daregex.ReplaceAllString(nes.Stream, "_"+indexstr)}
    visitor.ReturnExpr = newNes
    visitor.ReturnNumExpr = nil
    visitor.ReturnPred = nil
  } else {
    visitor.ReturnExpr = nes
    visitor.ReturnNumExpr = nil
    visitor.ReturnPred = nil
  }
}

func (visitor *NameToExprStreamVisitor) VisitAndPredicate(p AndPredicate) {
    p.Left.AcceptPred(visitor)
    pLeft := visitor.ReturnPred
    p.Right.AcceptPred(visitor)
    pRight := visitor.ReturnPred
    visitor.ReturnPred = AndPredicate{pLeft, pRight}
}
func (visitor *NameToExprStreamVisitor) VisitTruePredicate(p TruePredicate) {
    visitor.ReturnPred = p
}
func (visitor *NameToExprStreamVisitor) VisitFalsePredicate(p FalsePredicate) {
    visitor.ReturnPred = p
}
func (visitor *NameToExprStreamVisitor) VisitNotPredicate(p NotPredicate) {
	p.Inner.AcceptPred(visitor)
  visitor.ReturnPred = NotPredicate{visitor.ReturnPred}
}
func (visitor *NameToExprStreamVisitor) VisitOrPredicate(p OrPredicate) {
    p.Left.AcceptPred(visitor)
    pLeft := visitor.ReturnPred
    p.Right.AcceptPred(visitor)
    pRight := visitor.ReturnPred
    visitor.ReturnPred = OrPredicate{pLeft, pRight}
}
func (visitor *NameToExprStreamVisitor) VisitPathPredicate(p PathPredicate) {
    visitor.ReturnPred = p
}
func (visitor *NameToExprStreamVisitor) VisitStrCmpPredicate(p StrCmpPredicate) {
    visitor.ReturnPred = p
}
func (visitor *NameToExprStreamVisitor) VisitStrMatchPredicate(p StrMatchPredicate) {
    visitor.ReturnPred = p
}
func (visitor *NameToExprStreamVisitor) VisitTagPredicate(p TagPredicate) {
    visitor.ReturnPred = p
}
func (visitor *NameToExprStreamVisitor) VisitNamedPredicate(p StreamNameExpr) {
  p.Accept(visitor)
}
func (visitor *NameToExprStreamVisitor) VisitLastOfStreamNamedPredicate(p LastOfStreamNameExpr) {
  p.Accept(visitor)
}
func (visitor *NameToExprStreamVisitor) VisitNumComparisonPredicate(p NumComparisonPredicate) {
    p.NumComparison.Accept(visitor)
    visitor.ReturnPred = NumComparisonPredicate{visitor.ReturnNumComparison}
}

func (visitor *NameToExprStreamVisitor) VisitIfThenElsePredicate(p IfThenElsePredicate) {
    panic("not implemented")
}

func (visitor *NameToExprStreamVisitor) VisitPrevPredicate(p PrevPredicate) {
    visitor.ReturnPred = p
}

func (visitor *NameToExprStreamVisitor) VisitIsInitPredicate(p IsInitPredicate) {
    visitor.ReturnPred = p
}

// It also visits numcomparisons!

func (visitor *NameToExprStreamVisitor) VisitNumLess(exp NumLess) {
    exp.Left.AcceptNum(visitor)
    pLeft := visitor.ReturnNumExpr
    exp.Right.AcceptNum(visitor)
    pRight := visitor.ReturnNumExpr
    visitor.ReturnNumComparison = NumLess{pLeft, pRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumLessEq(exp NumLessEq) {
    exp.Left.AcceptNum(visitor)
    pLeft := visitor.ReturnNumExpr
    exp.Right.AcceptNum(visitor)
    pRight := visitor.ReturnNumExpr
    visitor.ReturnNumComparison = NumLessEq{pLeft, pRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumEq(exp NumEq) {
    exp.Left.AcceptNum(visitor)
    pLeft := visitor.ReturnNumExpr
    exp.Right.AcceptNum(visitor)
    pRight := visitor.ReturnNumExpr
    visitor.ReturnNumComparison = NumEq{pLeft, pRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumGreater(exp NumGreater) {
    exp.Left.AcceptNum(visitor)
    pLeft := visitor.ReturnNumExpr
    exp.Right.AcceptNum(visitor)
    pRight := visitor.ReturnNumExpr
    visitor.ReturnNumComparison = NumGreater{pLeft, pRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumGreaterEq(exp NumGreaterEq) {
    exp.Left.AcceptNum(visitor)
    pLeft := visitor.ReturnNumExpr
    exp.Right.AcceptNum(visitor)
    pRight := visitor.ReturnNumExpr
    visitor.ReturnNumComparison = NumGreaterEq{pLeft, pRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumNotEq(exp NumNotEq) {
    exp.Left.AcceptNum(visitor)
    pLeft := visitor.ReturnNumExpr
    exp.Right.AcceptNum(visitor)
    pRight := visitor.ReturnNumExpr
    visitor.ReturnNumComparison = NumNotEq{pLeft, pRight}
}

// And also visits NumExps!

func (visitor *NameToExprStreamVisitor) VisitIntLiteralExpr(exp IntLiteralExpr) {
  visitor.ReturnNumExpr = exp
}
func (visitor *NameToExprStreamVisitor) VisitFloatLiteralExpr(exp FloatLiteralExpr) {
  visitor.ReturnNumExpr = exp
}
func (visitor *NameToExprStreamVisitor) VisitNumMulExpr(exp NumMulExpr) {
  exp.Left.AcceptNum(visitor)
  rLeft := visitor.ReturnNumExpr
  exp.Right.AcceptNum(visitor)
  rRight := visitor.ReturnNumExpr
	visitor.ReturnNumExpr = NumMulExpr{rLeft, rRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumDivExpr(exp NumDivExpr) {
  exp.Left.AcceptNum(visitor)
  rLeft := visitor.ReturnNumExpr
  exp.Right.AcceptNum(visitor)
  rRight := visitor.ReturnNumExpr
	visitor.ReturnNumExpr = NumDivExpr{rLeft, rRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumPlusExpr(exp NumPlusExpr) {
  exp.Left.AcceptNum(visitor)
  rLeft := visitor.ReturnNumExpr
  exp.Right.AcceptNum(visitor)
  rRight := visitor.ReturnNumExpr
	visitor.ReturnNumExpr = NumPlusExpr{rLeft, rRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumMinusExpr(exp NumMinusExpr) {
  exp.Left.AcceptNum(visitor)
  rLeft := visitor.ReturnNumExpr
  exp.Right.AcceptNum(visitor)
  rRight := visitor.ReturnNumExpr
	visitor.ReturnNumExpr = NumMinusExpr{rLeft, rRight}
}
func (visitor *NameToExprStreamVisitor) VisitNumPathExpr(exp NumPathExpr) {
  visitor.ReturnNumExpr = exp
}
