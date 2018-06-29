package impl

import(
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type SignalNamesFromPredicateVisitor struct {
    Preds map[string]common.Predicate
    SNames []striverdt.StreamName
    Momvisitor *MoMToStriverVisitor
}

func (visitor *SignalNamesFromPredicateVisitor) VisitAndPredicate(p common.AndPredicate) {
    p.Left.AcceptPred(visitor)
    p.Right.AcceptPred(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitTruePredicate(p common.TruePredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitFalsePredicate(p common.FalsePredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNotPredicate(p common.NotPredicate) {
	p.Inner.AcceptPred(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitOrPredicate(p common.OrPredicate) {
    p.Left.AcceptPred(visitor)
    p.Right.AcceptPred(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitPathPredicate(p common.PathPredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitStrCmpPredicate(p common.StrCmpPredicate) {
    p.Expected.AcceptComparableStringVisitor(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitStrMatchPredicate(p common.StrMatchPredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitTagPredicate(p common.TagPredicate) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNamedPredicate(p common.StreamNameExpr) {
    if thepred, ok := visitor.Preds[p.Stream]; ok {
        thepred.AcceptPred(visitor)
    } else {
        visitor.SNames = append(visitor.SNames, striverdt.StreamName(p.Stream))
    }
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumComparisonPredicate(p common.NumComparisonPredicate) {
    p.NumComparison.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitIfThenElsePredicate(p common.IfThenElsePredicate) {
    p.If.AcceptPred(visitor)
    p.Then.Accept(visitor)
    p.Else.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitPrevPredicate(prevExp common.PrevPredicate) {
    outStream := "prevOf::"+prevExp.Stream
    makePrevOutStream(prevExp.Stream, outStream, visitor.Momvisitor)
    visitor.SNames = append(visitor.SNames,outStream)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitIsInitPredicate(isinit common.IsInitPredicate) {
    outStream := "isInit::"+isinit.Stream
    makeIsInitStream(isinit.Stream, outStream, visitor.Momvisitor)
    visitor.SNames = append(visitor.SNames,outStream)
}

// It also visits numcomparisons!

func (visitor *SignalNamesFromPredicateVisitor) VisitNumLess(exp common.NumLess) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumLessEq(exp common.NumLessEq) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumEq(exp common.NumEq) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumGreater(exp common.NumGreater) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumGreaterEq(exp common.NumGreaterEq) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumNotEq(exp common.NumNotEq) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}

// And also visits NumExps!

func (visitor *SignalNamesFromPredicateVisitor) VisitIntLiteralExpr(exp common.IntLiteralExpr) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitFloatLiteralExpr(exp common.FloatLiteralExpr) {
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumMulExpr(exp common.NumMulExpr) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumDivExpr(exp common.NumDivExpr) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumPlusExpr(exp common.NumPlusExpr) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumMinusExpr(exp common.NumMinusExpr) {
    exp.Left.AcceptNum(visitor)
    exp.Right.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitNumPathExpr(exp common.NumPathExpr) {
}

// Come on, it visits everything..

func (visitor *SignalNamesFromPredicateVisitor) VisitAggregatorExpr(aggexp common.AggregatorExpr) {
    visitor.SNames = append(visitor.SNames,aggexp.Stream, aggexp.Session)
}

func (visitor *SignalNamesFromPredicateVisitor) VisitIfThenExpr(ifthen common.IfThenExpr) {
    ifthen.If.AcceptPred(visitor)
    ifthen.Then.Accept(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitIfThenElseExpr(common.IfThenElseExpr) {
    panic("not implemented")
}
func (visitor *SignalNamesFromPredicateVisitor) VisitPredExpr(predExp common.PredExpr) {
    predExp.Pred.AcceptPred(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitStringPathExpr(common.StringPathExpr) {
    panic("not implemented")
}
func (visitor *SignalNamesFromPredicateVisitor) VisitStreamNumExpr(numExp common.StreamNumExpr) {
    numExp.Expr.AcceptNum(visitor)
}
func (visitor *SignalNamesFromPredicateVisitor) VisitStreamNameExpr(exp common.StreamNameExpr) {
    visitor.SNames = append(visitor.SNames, striverdt.StreamName(exp.Stream))
}

// Even comparable strings

func (visitor *SignalNamesFromPredicateVisitor) VisitQuotedString(qs common.QuotedString) {
}

func (visitor *SignalNamesFromPredicateVisitor) VisitIdentifier(id common.Identifier) {
    visitor.SNames = append(visitor.SNames, striverdt.StreamName(id.Val))
}

// constructors of streams

func makePrevOutStream (inSignalName, outSignalName striverdt.StreamName, visitor *MoMToStriverVisitor) {
    hasEverFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        myprev := args[1]
        if myprev.IsSet && myprev.Val.(striverdt.EvPayload).Val.(bool) {
            return striverdt.Some(true)
        }
        return args[0].Val.(striverdt.EvPayload)
    }
    hasEverVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevEqValNode{striverdt.TNode{}, inSignalName, []striverdt.Event{}},
        &striverdt.PrevValNode{striverdt.TNode{}, outSignalName, []striverdt.Event{}},
    }, hasEverFun}
    hasEverStream := striverdt.OutStream{outSignalName, striverdt.SrcTickerNode{inSignalName}, hasEverVal}
    visitor.OutStreams = append(visitor.OutStreams, hasEverStream)
}

func makeIsInitStream (inSignalName, outSignalName striverdt.StreamName, visitor *MoMToStriverVisitor) {
    isinitFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
        strprev := args[0]
        return striverdt.Some(strprev.IsSet)
    }
    isinitVal := striverdt.FuncNode{[]striverdt.ValNode{
        &striverdt.PrevValNode{striverdt.TNode{}, inSignalName, []striverdt.Event{}},
    }, isinitFun}
    isinitStream := striverdt.OutStream{outSignalName, striverdt.SrcTickerNode{visitor.InSignalName}, isinitVal}
    visitor.OutStreams = append(visitor.OutStreams, isinitStream)
}
