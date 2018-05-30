package common

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
)

type EvalVisitor struct {
    Result bool
    Event dt.Event
}

func (visitor *EvalVisitor) visitAndPredicate(p AndPredicate) {
    p.Left.Accept(visitor)
    rLeft := visitor.Result
    p.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft && rRight
}
func (visitor *EvalVisitor) visitTruePredicate(p TruePredicate) {
    visitor.Result = true
}
func (visitor *EvalVisitor) visitFalsePredicate(p FalsePredicate) {
    visitor.Result = false
}
func (visitor *EvalVisitor) visitNotPredicate(p NotPredicate) {
	p.Inner.Accept(visitor)
    visitor.Result = !visitor.Result
}
func (visitor *EvalVisitor) visitOrPredicate(p OrPredicate) {
    p.Left.Accept(visitor)
    rLeft := visitor.Result
    p.Right.Accept(visitor)
    rRight := visitor.Result
	visitor.Result = rLeft || rRight
}
func (visitor *EvalVisitor) visitPathPredicate(p PathPredicate) {
    _,err := jsonrw.ExtractFromMap(visitor.Event.Payload, dt.JSONPath(p.Path))
	visitor.Result = err == nil
}
func (visitor *EvalVisitor) visitStrPredicate(p StrPredicate) {
    strif,err := jsonrw.ExtractFromMap(visitor.Event.Payload, dt.JSONPath(p.Path))
    if err != nil {
        visitor.Result = false
        //fmt.Println("No string found in event ", visitor.Event)
        return
    }
    visitor.Result = strif.(string) == p.Expected
    //fmt.Println("Comparing",strif, "with", p.Expected, "and the result is", visitor.Result)
}
func (visitor *EvalVisitor) visitTagPredicate(p TagPredicate) {
    visitor.Result = sets.SetIn(p.Tag, visitor.Event.Channels)
}
func (visitor *EvalVisitor) visitNamedPredicate(p NamedPredicate) {
    // TODO
}
func (visitor *EvalVisitor) visitNumComparisonPredicate(p NumComparisonPredicate) {
    // TODO
}


//
// sprint() functions of the different Predicates. TODO: Use visitor
//
// func (p AndPredicate) Sprint() string {
// 	return fmt.Sprintf("(%s /\\ %s)",p.Left.Sprint(),p.Right.Sprint())
// }
// func (p OrPredicate) Sprint() string {
// 	return fmt.Sprintf("(%s \\/  %s)",p.Left.Sprint(),p.Right.Sprint())
// }
// func (p NotPredicate) Sprint() string {
// 	return fmt.Sprintf("~ %s",p.Inner.Sprint())
// }
// func (p PathPredicate) Sprint() string {
// 	return fmt.Sprintf("e.path(%s)",p.Path)
// }
// func (p StrPredicate) Sprint() string {
// 	return fmt.Sprintf("e.strcmp(%s,\"%s\")",p.Path,p.Expected)
// }
// func (p TagPredicate) Sprint() string {
// 	return fmt.Sprintf("e.tag(%s)",p.Tag)
// }
// func (p TruePredicate) Sprint() string {
// 	return fmt.Sprintf("true");
// }
// func (p FalsePredicate) Sprint() string {
// 	return fmt.Sprintf("false")
// }
// func (p NamedPredicate) Sprint() string {
// 	return p.Name
// }
