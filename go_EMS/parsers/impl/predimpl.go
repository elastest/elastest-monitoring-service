package impl

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
	"github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "regexp"
    "time"
)

type ComparableStringEvaluator struct {
    ArgsMap map[striverdt.StreamName]interface{}
    Result string
}

func (visitor *ComparableStringEvaluator) VisitQuotedString(qs common.QuotedString) {
    visitor.Result = qs.Val
}

func (visitor *ComparableStringEvaluator) VisitIdentifier(id common.Identifier) {
    visitor.Result = visitor.ArgsMap[striverdt.StreamName(id.Val)].(string)
}

type EvalVisitor struct {
    Result bool
    Event dt.Event
    Preds map[string]common.Predicate
    ArgsMap map[striverdt.StreamName]interface{}
}

type EvalNumVisitor struct {
    Result float32
    ResultInt int64
    IsInt bool
    Event dt.Event
    ArgsMap map[striverdt.StreamName]interface{}
}

func (visitor *EvalVisitor) VisitAndPredicate(p common.AndPredicate) {
    p.Left.AcceptPred(visitor)
    rLeft := visitor.Result
    p.Right.AcceptPred(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft && rRight
}
func (visitor *EvalVisitor) VisitTruePredicate(p common.TruePredicate) {
    visitor.Result = true
}
func (visitor *EvalVisitor) VisitFalsePredicate(p common.FalsePredicate) {
    visitor.Result = false
}
func (visitor *EvalVisitor) VisitNotPredicate(p common.NotPredicate) {
	p.Inner.AcceptPred(visitor)
    visitor.Result = !visitor.Result
}
func (visitor *EvalVisitor) VisitOrPredicate(p common.OrPredicate) {
    p.Left.AcceptPred(visitor)
    rLeft := visitor.Result
    p.Right.AcceptPred(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft || rRight
}
func (visitor *EvalVisitor) VisitPathPredicate(p common.PathPredicate) {
    _,err := jsonrw.ExtractFromMap2(visitor.Event.Payload, p.Path)
	visitor.Result = err == nil
}
func (visitor *EvalVisitor) VisitStrCmpPredicate(p common.StrCmpPredicate) {
    strif,err := jsonrw.ExtractFromMap2(visitor.Event.Payload, p.Path)
    if err != nil {
        visitor.Result = false
        //fmt.Println("No string found in event ", visitor.Event)
        return
    }
    compeval := ComparableStringEvaluator{visitor.ArgsMap, ""}
    p.Expected.AcceptComparableStringVisitor(&compeval)
    visitor.Result = compeval.Result == strif.(string)
    //fmt.Println("Comparing",strif, "with", p.Expected, "and the result is", visitor.Result)
}
func (visitor *EvalVisitor) VisitStrMatchPredicate(p common.StrMatchPredicate) {
    strif,err := jsonrw.ExtractFromMap2(visitor.Event.Payload, []dt.JSONPath(p.Path))
    if err != nil {
        visitor.Result = false
        //fmt.Println("No string found in event ", visitor.Event)
        return
    }
    matched, err := regexp.MatchString(p.Expected, strif.(string))
    visitor.Result = err == nil && matched
    // FIXME compile regex
    //fmt.Println("Comparing",strif, "with", p.Expected, "and the result is", visitor.Result)
}
func (visitor *EvalVisitor) VisitTagPredicate(p common.TagPredicate) {
    visitor.Result = sets.SetIn(p.Tag, visitor.Event.Channels)
}
func (visitor *EvalVisitor) VisitNamedPredicate(p common.StreamNameExpr) {
    if thepred, ok := visitor.Preds[p.Stream]; ok {
        thepred.AcceptPred(visitor)
    } else {
        visitor.Result = visitor.ArgsMap[striverdt.StreamName(p.Stream)].(bool)
    }
}
func (visitor *EvalVisitor) VisitNumComparisonPredicate(p common.NumComparisonPredicate) {
    p.NumComparison.Accept(visitor)
}

func (visitor *EvalVisitor) VisitPrevPredicate(p common.PrevPredicate) {
    visitor.Result = visitor.ArgsMap["prevOf::"+p.Stream].(bool)
}

func (visitor *EvalVisitor) VisitIsInitPredicate(p common.IsInitPredicate) {
    visitor.Result = visitor.ArgsMap["isInit::"+p.Stream].(bool)
}

// It also visits numcomparisons!

func (visitor *EvalVisitor) VisitNumLess(exp common.NumLess) {
    numvisitor := EvalNumVisitor{0, 0, false, visitor.Event, visitor.ArgsMap}
    exp.Left.AcceptNum(&numvisitor)
    a := numvisitor.Result
    exp.Right.AcceptNum(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a<b
}
func (visitor *EvalVisitor) VisitNumLessEq(exp common.NumLessEq) {
    numvisitor := EvalNumVisitor{0, 0, false, visitor.Event, visitor.ArgsMap}
    exp.Left.AcceptNum(&numvisitor)
    a := numvisitor.Result
    exp.Right.AcceptNum(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a<=b
}
func (visitor *EvalVisitor) VisitNumEq(exp common.NumEq) {
    numvisitor := EvalNumVisitor{0, 0, false, visitor.Event, visitor.ArgsMap}
    exp.Left.AcceptNum(&numvisitor)
    a := numvisitor.Result
    exp.Right.AcceptNum(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a==b
}
func (visitor *EvalVisitor) VisitNumGreater(exp common.NumGreater) {
    numvisitor := EvalNumVisitor{0, 0, false, visitor.Event, visitor.ArgsMap}
    exp.Left.AcceptNum(&numvisitor)
    a := numvisitor.Result
    exp.Right.AcceptNum(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a>b
}
func (visitor *EvalVisitor) VisitNumGreaterEq(exp common.NumGreaterEq) {
    numvisitor := EvalNumVisitor{0, 0, false, visitor.Event, visitor.ArgsMap}
    exp.Left.AcceptNum(&numvisitor)
    a := numvisitor.Result
    exp.Right.AcceptNum(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a>=b
}
func (visitor *EvalVisitor) VisitNumNotEq(exp common.NumNotEq) {
    numvisitor := EvalNumVisitor{0, 0, false, visitor.Event, visitor.ArgsMap}
    exp.Left.AcceptNum(&numvisitor)
    a := numvisitor.Result
    exp.Right.AcceptNum(&numvisitor)
    b := numvisitor.Result
    visitor.Result = a!=b
}

// And also visits NumExps!

func (visitor *EvalNumVisitor) VisitIntLiteralExpr(exp common.IntLiteralExpr) {
    visitor.Result = float32(exp.Num)
}
func (visitor *EvalNumVisitor) VisitFloatLiteralExpr(exp common.FloatLiteralExpr) {
    visitor.Result = exp.Num
}
func (visitor *EvalNumVisitor) VisitStreamNameExpr(exp common.StreamNameExpr) {
    switch castedval:= visitor.ArgsMap[striverdt.StreamName(exp.Stream)].(type) {
    case float32:
      visitor.Result = castedval
      visitor.IsInt = false
    case int64:
      visitor.Result = float32(castedval)
      visitor.ResultInt = castedval
      visitor.IsInt = true
    }
}
func (visitor *EvalNumVisitor) VisitNumMulExpr(exp common.NumMulExpr) {
    exp.Left.AcceptNum(visitor)
    rLeft := visitor.Result
    visitor.IsInt = false
    exp.Right.AcceptNum(visitor)
    rRight := visitor.Result
      visitor.IsInt = false
	visitor.Result = rLeft * rRight
}
func (visitor *EvalNumVisitor) VisitNumDivExpr(exp common.NumDivExpr) {
    exp.Left.AcceptNum(visitor)
    rLeft := visitor.Result
    visitor.IsInt = false
    exp.Right.AcceptNum(visitor)
    rRight := visitor.Result
      visitor.IsInt = false
	visitor.Result = rLeft / rRight
}
func (visitor *EvalNumVisitor) VisitNumPlusExpr(exp common.NumPlusExpr) {
    exp.Left.AcceptNum(visitor)
    rLeft := visitor.Result
    visitor.IsInt = false
    exp.Right.AcceptNum(visitor)
    rRight := visitor.Result
      visitor.IsInt = false
	visitor.Result = rLeft + rRight
}
func (visitor *EvalNumVisitor) VisitNumMinusExpr(exp common.NumMinusExpr) {
    exp.Left.AcceptNum(visitor)
    rLeft := visitor.Result
    leftIsInt := visitor.IsInt
    leftInt := visitor.ResultInt
    visitor.IsInt = false
    exp.Right.AcceptNum(visitor)
    rRight := visitor.Result
    if visitor.IsInt && leftIsInt {
      visitor.Result = float32(leftInt - visitor.ResultInt)
    } else {
      visitor.Result = rLeft - rRight
    }
    visitor.IsInt = false
}
func (visitor *EvalNumVisitor) VisitNumPathExpr(exp common.NumPathExpr) {
    theEvent := visitor.Event
    valif, err := jsonrw.ExtractFromMap2(theEvent.Payload, exp.ExtractPaths)
    visitor.Result = -9999999
    visitor.ResultInt = -9999999
    visitor.IsInt = false
    if err == nil {
        /* This may not happen: the stream might be guarded by an if statement upper in the AST.
        Perhaps we should panic and fix if statements to not evaluate
        the inner function if the result is false? */
        switch castedval := valif.(type) {
        case float64:
          visitor.Result = float32(castedval)
        case string:
          t,e := time.Parse(time.RFC3339, castedval)
          if e==nil {
            visitor.ResultInt = t.UnixNano()/1000000
            visitor.IsInt = true
          }
        }
    }
}
