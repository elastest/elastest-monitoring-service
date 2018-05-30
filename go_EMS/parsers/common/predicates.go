package common

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
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
    VisitNamedPredicate(NamedPredicate)
    VisitNumComparisonPredicate(NumComparisonPredicate)
}

type Predicate interface {
    Accept (PredicateVisitor)
}

type TruePredicate  struct {}

func (this TruePredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitTruePredicate(this)
}

type FalsePredicate struct {}

func (this FalsePredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitFalsePredicate(this)
}

type NotPredicate struct {
	Inner Predicate
}

func (this NotPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitNotPredicate(this)
}

type AndPredicate struct {
	Left  Predicate
	Right Predicate
}

func (this AndPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitAndPredicate(this)
}

type OrPredicate struct {
	Left  Predicate
	Right Predicate
}

func (this OrPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitOrPredicate(this)
}

type PathPredicate struct {
	Path string
}

func (this PathPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitPathPredicate(this)
}

type StrPredicate struct {
	Path string
	Expected string
}

func (this StrPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitStrPredicate(this)
}

type TagPredicate struct {
	Tag dt.Channel
}

func (this TagPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitTagPredicate(this)
}

type NamedPredicate struct {
	// This can be either a predicate "foo" defined with "pred foo :="
	// or a stream named "foo" defined "stream boolean foo :=.."
	Name string
}

func (this NamedPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitNamedPredicate(this)
}

type NumComparisonPredicate struct {
    NumComparison NumComparison
}

func (this NumComparisonPredicate) Accept (visitor PredicateVisitor) {
    visitor.VisitNumComparisonPredicate(this)
}

var (
	True  TruePredicate
	False FalsePredicate
)


// Constructors

func NewAndPredicate(a, b interface{}) (Predicate) {
	preds := ToSlice(b)
	first := a.(Predicate)
	if len(preds)==0 {
		return first
	}
	right := preds[len(preds)-1].(Predicate)
	for i := len(preds)-2; i >= 0; i-- {
		left := preds[i].(Predicate)
		right = AndPredicate{left,right}
	}
	ret := AndPredicate{first,right}
	return ret
}

func NewOrPredicate(a, b interface{}) (Predicate) {
	preds := ToSlice(b)
	first := a.(Predicate)
	if len(preds)==0 {
		return first
	}
	right := preds[len(preds)-1].(Predicate)
	for i := len(preds)-2; i >= 0; i-- {
		left := preds[i].(Predicate)
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

func NewNamedPredicate(n interface{}) NamedPredicate {
	name := n.(Identifier).Val
	return NamedPredicate{name}
}

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
