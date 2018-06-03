package common

import(
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
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
    VisitTagPredicate(TagPredicate)
    VisitNamedPredicate(StreamNameExpr)
    VisitNumComparisonPredicate(NumComparisonPredicate)
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
	return fmt.Sprint("~ %s", this.Inner.Sprint())
}
type AndPredicate struct {
	Left  Predicate
	Right Predicate
}
func (this AndPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitAndPredicate(this)
}
func (this AndPredicate) Sprint () string {
	return fmt.Sprint("%s /\\ %s",this.Left.Sprint(),this.Right.Sprint())
}
type OrPredicate struct {
	Left  Predicate
	Right Predicate
}
func (this OrPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitOrPredicate(this)
}
func (this OrPredicate) Sprint () string {
	return fmt.Sprint("%s \\/ %s",this.Left.Sprint(),this.Right.Sprint())
}
type PathPredicate struct {
	Path string
}
func (this PathPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitPathPredicate(this)
}
func (this PathPredicate) Sprint () string {
	return fmt.Sprint("e.Path(%s)",this.Path)
}
type StrPredicate struct {
	Path string
	Expected string
}
func (this StrPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitStrPredicate(this)
}
func (this StrPredicate) Sprint () string {
	return fmt.Sprint("e.strcmp(%s,%s)",this.Path,this.Expected)
}
type TagPredicate struct {
	Tag dt.Channel
}
func (this TagPredicate) AcceptPred(visitor PredicateVisitor) {
    visitor.VisitTagPredicate(this)
}
func (this TagPredicate) Sprint () string {
	return fmt.Sprint("e.tag(%s)",this.Tag)
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
