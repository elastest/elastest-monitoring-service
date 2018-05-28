package common

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
)

type PredicateVisitor interface {
    visitTruePredicate(TruePredicate)
    visitFalsePredicate(FalsePredicate)
    visitNotPredicate(NotPredicate)
    visitAndPredicate(AndPredicate)
    visitOrPredicate(OrPredicate)
    visitPathPredicate(PathPredicate)
    visitStrPredicate(StrPredicate)
    visitTagPredicate(TagPredicate)
    visitNamedPredicate(NamedPredicate)
}

type Predicate interface {
    Accept (PredicateVisitor)
	Sprint() string // TODO remove later on
}

type TruePredicate  struct {}

func (this TruePredicate) Accept (visitor PredicateVisitor) {
    visitor.visitTruePredicate(this)
}

type FalsePredicate struct {}

func (this FalsePredicate) Accept (visitor PredicateVisitor) {
    visitor.visitFalsePredicate(this)
}

type NotPredicate struct {
	Inner Predicate
}

func (this NotPredicate) Accept (visitor PredicateVisitor) {
    visitor.visitNotPredicate(this)
}

type AndPredicate struct {
	Left  Predicate
	Right Predicate
}

func (this AndPredicate) Accept (visitor PredicateVisitor) {
    visitor.visitAndPredicate(this)
}

type OrPredicate struct {
	Left  Predicate
	Right Predicate
}

func (this OrPredicate) Accept (visitor PredicateVisitor) {
    visitor.visitOrPredicate(this)
}

type PathPredicate struct {
	Path string
}

func (this PathPredicate) Accept (visitor PredicateVisitor) {
    visitor.visitPathPredicate(this)
}

type StrPredicate struct {
	Path string
	Expected string
}

func (this StrPredicate) Accept (visitor PredicateVisitor) {
    visitor.visitStrPredicate(this)
}

type TagPredicate struct {
	Tag dt.Channel
}

func (this TagPredicate) Accept (visitor PredicateVisitor) {
    visitor.visitTagPredicate(this)
}

type NamedPredicate struct {
	// This can be either a predicate "foo" defined with "pred foo :="
	// or a stream named "foo" defined "stream boolean foo :=.."
	Name string
}

func (this NamedPredicate) Accept (visitor PredicateVisitor) {
    visitor.visitNamedPredicate(this)
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

// Helper functions

func ToSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})

}
