package common

import(
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
	"fmt"
	"errors"
)

type PredicateVisitor interface {
    VisitTruePredicate(TruePredicate)
    VisitFalsePredicate(FalsePredicate)
    VisitNotPredicate(NotPredicate)
    VisitAndPredicate(AndPredicate)
    VisitOrPredicate(OrPredicate)
    VisitPathPredicate(PathPredicate)
    VisitStrPredicate(StrPredicate)
    VisitStrMatchPredicate(StrMatchPredicate)
    VisitTagPredicate(TagPredicate)
    VisitNamedPredicate(StreamNameExpr)
	VisitNumComparisonPredicate(NumComparisonPredicate)
	VisitIfThenElsePredicate(IfThenElsePredicate)
	VisitPrevPredicate(PrevPredicate)

}

type Predicate interface {
	AcceptPred(PredicateVisitor)
	Sprint() string
}

//
// Predicates can become StreamExpr
//
func getStreamExpr(p Predicate) StreamExpr {
	return NewPredExpr(p)
}

type TruePredicate  struct {}

func (this TruePredicate) AcceptPred (visitor PredicateVisitor) {
    visitor.VisitTruePredicate(this)
}
func (this TruePredicate) Sprint () string {
	return "true"
}
type FalsePredicate struct {}

func (this FalsePredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitFalsePredicate(this)
}
func (this FalsePredicate) Sprint () string {
	return "false"
}

type NotPredicate struct {
	Inner Predicate
}
func (this NotPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitNotPredicate(this)
}
func (this NotPredicate) Sprint () string {
	return fmt.Sprintf("~ %s", this.Inner.Sprint())
}
type AndPredicate struct {
	Left  Predicate
	Right Predicate
}
func (this AndPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitAndPredicate(this)
}
func (this AndPredicate) Sprint () string {
	return fmt.Sprintf("%s /\\ %s",this.Left.Sprint(),this.Right.Sprint())
}
type OrPredicate struct {
	Left  Predicate
	Right Predicate
}
func (this OrPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitOrPredicate(this)
}
func (this OrPredicate) Sprint () string {
	return fmt.Sprintf("%s \\/ %s",this.Left.Sprint(),this.Right.Sprint())
}
type PathPredicate struct {
	Path string
}
func (this PathPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitPathPredicate(this)
}
func (this PathPredicate) Sprint () string {
	return fmt.Sprintf("e.Path(%s)",this.Path)
}
type StrPredicate struct {
	Path string
	Expected string
}
func (this StrPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitStrPredicate(this)
}
func (this StrPredicate) Sprint () string {
	return fmt.Sprintf("e.strcmp(%s,%s)",this.Path,this.Expected)
}
type StrMatchPredicate struct {
	Path string
	Expected string
}
func (this StrMatchPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitStrMatchPredicate(this)
}
func (this StrMatchPredicate) Sprint () string {
	return fmt.Sprintf("e.strmatch(%s,%s)",this.Path,this.Expected)
}
type TagPredicate struct {
	Tag dt.Channel
}
func (this TagPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitTagPredicate(this)
}
func (this TagPredicate) Sprint () string {
	return fmt.Sprintf("e.tag(%s)",this.Tag)
}

// type NamedPredicate struct {
//	// This can be either a predicate "foo" defined with "pred foo :="
//	// or a stream named "foo" defined "stream boolean foo :=.."
//	Name string
//}

func (this StreamNameExpr) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitNamedPredicate(this)
}

type NumComparisonPredicate struct {
    NumComparison NumComparison
}

func (this NumComparisonPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitNumComparisonPredicate(this)
}
func (this NumComparisonPredicate) Sprint() string {
	return this.NumComparison.Sprint()
}

type IfThenElsePredicate struct {
	If Predicate
	Then StreamExpr
	Else StreamExpr
}
func (this IfThenElsePredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitIfThenElsePredicate(this)
}
func (this IfThenElsePredicate) Sprint() string {
	expr := IfThenElseExpr{this.If,this.Then,this.Else}
	return expr.Sprint()
}
func IfThenElseExpr2Pred(e IfThenElseExpr) IfThenElsePredicate {
	return IfThenElsePredicate{e.If,e.Then,e.Else}
}

type PrevPredicate struct {
	Stream striverdt.StreamName
}

func (this PrevPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitPrevPredicate(this)
}
func (this PrevPredicate) Sprint() string {
	return fmt.Sprintf("Prev %s",this.Stream)
}
func NewPrevPred(p interface{}) (PrevPredicate) {
	return PrevPredicate{striverdt.StreamName(p.(Identifier).Val)}
}


var (
	True  TruePredicate
	False FalsePredicate
	TrueExpr = PredExpr{True}
	FalseExpr = PredExpr{False}
)


func getPredExpr(a interface{}) (Predicate,error) {
	if v,ok:=a.(PredExpr) ; ok {
		return v.Pred,nil
	} else if v,ok:=a.(StreamNameExpr); ok {
		return v,nil         // StreamNAmeExpr implements Predicate
	} else {
		str := fmt.Sprintf("cannot convert to pred \"%s\"\n",a.(StreamExpr).Sprint())
		fmt.Printf(str)
		return nil,errors.New(str)
	}
}

func NewAndPredicate(a, b interface{}) (Predicate) {
	preds := ToSlice(b)
	first,_ := getPredExpr(a)
	if len(preds)==0 {
		return first
	}
	right,_ := getPredExpr(preds[len(preds)-1])
	for i := len(preds)-2; i >= 0; i-- {
		left,_ := getPredExpr(preds[i])
		right = AndPredicate{left,right}
	}
	ret := AndPredicate{first,right}
	return ret
}

func NewOrPredicate(a, b interface{}) (Predicate) {
	preds  := ToSlice(b)
	first,_:= getPredExpr(a)
	if len(preds)==0 {
		return first
	}
	right,_ :=getPredExpr(preds[len(preds)-1])
	for i := len(preds)-2; i >= 0; i-- {
		left,_ := getPredExpr(preds[i])
		right = OrPredicate{left,right}
	}
	return OrPredicate{first,right}
}

func NewNotPredicate(p interface{}) (NotPredicate) {
	return NotPredicate{p.(Predicate)}
}

func NewPathPredicate(p interface{}) (PathPredicate) {
	path := p.(PathName).Val
	return PathPredicate{path}
}

func NewStrPredicate(p,v interface{}) (StrPredicate) {
	path     :=p.(PathName).Val
	expected :=v.(QuotedString).Val
	return StrPredicate{path,expected}
}

func NewStrMatchPredicate(p,v interface{}) (StrMatchPredicate) {
	path     :=p.(PathName).Val
	expected :=v.(QuotedString).Val
	return StrMatchPredicate{path,expected}
}

func NewTagPredicate(t interface{}) (TagPredicate) {
	tag := t.(Tag).Tag
	return TagPredicate{tag}
}

// func NewNamedPredicate(n interface{}) NamedPredicate {
// 	name := n.(Identifier).Val
//	return NamedPredicate{name}
//}

func NewNumComparisonPredicate(n interface{}) NumComparisonPredicate {
    return NumComparisonPredicate{n.(NumComparison)}
}

// Helper functions

func ToSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})

}
