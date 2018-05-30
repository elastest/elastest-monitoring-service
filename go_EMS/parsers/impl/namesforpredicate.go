package impl

import(
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type SignalNamesFromPredicateVisitor struct {
    Preds map[string]common.Predicate
    SNames []striverdt.StreamName
}

func (visitor *SignalNamesFromPredicateVisitor) VisitAndPredicate(p common.AndPredicate) {
    p.Left.Accept(visitor)
    p.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitTruePredicate(p common.TruePredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitFalsePredicate(p common.FalsePredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNotPredicate(p common.NotPredicate) {
	p.Inner.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitOrPredicate(p common.OrPredicate) {
    p.Left.Accept(visitor)
    p.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitPathPredicate(p common.PathPredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitStrPredicate(p common.StrPredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitTagPredicate(p common.TagPredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNamedPredicate(p common.NamedPredicate) {
    if thepred, ok := visitor.Preds[p.Name]; ok {
        thepred.Accept(visitor)
    } else {
        visitor.SNames = append(visitor.SNames, striverdt.StreamName(p.Name))
    }
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumComparisonPredicate(p common.NumComparisonPredicate) {
    p.NumComparison.Accept(visitor)
}

// It also visits numcomparisons!

func (visitor *SignalNamesFromPredicateVisitor) VisitNumLess(exp common.NumLess) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumLessEq(exp common.NumLessEq) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumEq(exp common.NumEq) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumGreater(exp common.NumGreater) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumGreaterEq(exp common.NumGreaterEq) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumNotEq(exp common.NumNotEq) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}

// And also visits NumExps!

func (visitor *SignalNamesFromPredicateVisitor) VisitIntLiteralExpr(exp common.IntLiteralExpr) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitFloatLiteralExpr(exp common.FloatLiteralExpr) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitStreamNameExpr(exp common.StreamNameExpr) {
    visitor.SNames = append(visitor.SNames, exp.StreamName)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumMulExpr(exp common.NumMulExpr) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumDivExpr(exp common.NumDivExpr) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumPlusExpr(exp common.NumPlusExpr) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumMinusExpr(exp common.NumMinusExpr) {
    exp.Left.Accept(visitor)
    exp.Right.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumPathExpr(exp common.NumPathExpr) {
}
