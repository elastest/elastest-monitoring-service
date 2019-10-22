package common

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "fmt"
    "errors"
//    "strconv"
)

type StreamExprVisitor interface {
	VisitAggregatorExpr(AggregatorExpr)
	VisitIfThenExpr(IfThenExpr)
	VisitIfThenElseExpr(IfThenElseExpr)
	// VisittringExpr(StringExpr)
	VisitStreamNumExpr(StreamNumExpr)
	VisitPredExpr(PredExpr)
	VisitStringPathExpr(StringPathExpr)
	VisitStreamNameExpr(StreamNameExpr)
}

type StreamExpr interface {
	// add functions here
	Sprint() string
	Accept (StreamExprVisitor)
}

type StreamNumExpr struct {
    Expr NumExpr
}
func (this StreamNumExpr) Accept(visitor StreamExprVisitor) {
	visitor.VisitStreamNumExpr(this)
}
func (this StreamNumExpr) Sprint() string {
	str := this.Expr.Sprint()
	return str
}

type AggregatorExpr struct {
	Operation string
	Stream    striverdt.StreamName //StreamName
	Session   striverdt.StreamName //StreamName
}
func (a AggregatorExpr) Sprint() string {
	return "FIXME"
}

func (this AggregatorExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitAggregatorExpr(this)
}

type IfThenExpr struct {
	If   Predicate
	Then StreamExpr
}
func (this IfThenExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitIfThenExpr(this)
}
func (this IfThenExpr) Sprint() string {
	return fmt.Sprintf("if %s then %s",this.If.Sprint(),this.Then.Sprint())
}


type IfThenElseExpr struct {
	If   Predicate
	Then StreamExpr
	Else StreamExpr
}

func (this IfThenElseExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitIfThenElseExpr(this)
}
func (a IfThenElseExpr) Sprint() string {
	return fmt.Sprintf("if %s then %s else %s",a.If.Sprint(),a.Then.Sprint(),a.Else.Sprint())
}
// Is this ever used?
//type StringExpr struct {
//	Path string// so far only e.get(path) claiming to return a string
//}

type PredExpr struct {
	Pred Predicate
}

func (this PredExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitPredExpr(this)
}
func (this PredExpr) Sprint() string {
	return this.Pred.Sprint()
}

type StreamNameExpr struct {
	Stream string
}
func (this StreamNameExpr) Accept(visitor StreamExprVisitor) {
	visitor.VisitStreamNameExpr(this)
}
func (this StreamNameExpr) Sprint() string {
	return string(this.Stream)
}

type StringPathExpr struct {
	Path dt.JSONPath
}

func (this StringPathExpr) Accept(visitor StreamExprVisitor) {
    visitor.VisitStringPathExpr(this)
}

func (this StringPathExpr) Sprint() string {
	return fmt.Sprintf("e.getstr(%s)",this.Path)
}

//
// Expression Node constructors
//
func NewStreamNumExpr(n interface{}) StreamNumExpr {
	return StreamNumExpr{n.(NumExpr)}
}
func NewAggregatorExpr(op, str, ses interface{}) AggregatorExpr {
	operation := op.(string)
	stream    := str.(Identifier).Val
	session   := ses.(Identifier).Val

	return AggregatorExpr{operation,striverdt.StreamName(stream),striverdt.StreamName(session)}
}

func NewIfThenExpr(p,e interface{}) IfThenExpr {
	if_part   := p.(Predicate)
	then_part := e.(StreamExpr)
	return IfThenExpr{if_part,then_part}
}
func NewIfThenElseExpr(p,a,b interface{}) IfThenElseExpr {
	if_part   := p.(Predicate)
	then_part := a.(StreamExpr)
	else_part := b.(StreamExpr)
	return IfThenElseExpr{if_part, then_part, else_part}
}
func NewPredExpr(p interface{}) PredExpr {
	return PredExpr{p.(Predicate)}
}

func NewStringPathExpr(p interface{}) StringPathExpr {
	path := p.(PathName).Val
	return StringPathExpr{dt.JSONPath(path)}
}

func NewStreamNameExpr(p interface{}) StreamNameExpr {
	return StreamNameExpr{p.(Identifier).Val}
}

func NewJSONExpr(suffixes interface{}) (JSONExpr, error) {
  isuffixes := suffixes.([]interface{})
  var paths []dt.JSONPath
  if len(isuffixes)==1 {
    paths = append(paths, dt.JSONPath(isuffixes[0].(PathName).Val))
  }
  if len(isuffixes)>1 {
	    return JSONExpr{paths},errors.New("Only one level of JSON embedding is supported")
  }
  return JSONExpr{paths}, nil
}
