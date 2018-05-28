package common

import(
	"fmt"
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    "github.com/elastest/elastest-monitoring-service/go_EMS/jsonrw"
)


type Predicate interface {
	Sprint() string
	Eval(e dt.Event) bool
}

type TruePredicate  struct {}
type FalsePredicate struct {}

type NotPredicate struct {
	Inner Predicate
}

type AndPredicate struct {
	Left  Predicate
	Right Predicate
}

type OrPredicate struct {
	Left  Predicate
	Right Predicate
}
type PathPredicate struct {
	Path string
}
type StrPredicate struct {
	Path string
	Expected string
}

type TagPredicate struct {
	Tag dt.Channel
}

type NamedPredicate struct {
	// This can be either a predicate "foo" defined with "pred foo :="
	// or a stream named "foo" defined "stream boolean foo :=.."
	Name string
}

var (
	True  TruePredicate
	False FalsePredicate
)

//
// sprint() functions of the different Predicates
//
func (p AndPredicate) Sprint() string {
	return fmt.Sprintf("(%s /\\ %s)",p.Left.Sprint(),p.Right.Sprint())
}
func (p OrPredicate) Sprint() string {
	return fmt.Sprintf("(%s \\/  %s)",p.Left.Sprint(),p.Right.Sprint())
}
func (p NotPredicate) Sprint() string {
	return fmt.Sprintf("~ %s",p.Inner.Sprint())
}
func (p PathPredicate) Sprint() string {
	return fmt.Sprintf("e.path(%s)",p.Path)
}
func (p StrPredicate) Sprint() string {
	return fmt.Sprintf("e.strcmp(%s,\"%s\")",p.Path,p.Expected)
}
func (p TagPredicate) Sprint() string {
	return fmt.Sprintf("e.tag(%s)",p.Tag)
}
func (p TruePredicate) Sprint() string {
	return fmt.Sprintf("true");
}
func (p FalsePredicate) Sprint() string {
	return fmt.Sprintf("false")
}
func (p NamedPredicate) Sprint() string {
	return p.Name
}

//
// eval(dt.Event e) bool
//
func (p AndPredicate) Eval(e dt.Event) bool {
	return p.Left.Eval(e) && p.Right.Eval(e)
}
func (p OrPredicate) Eval(e dt.Event) bool {
	return p.Left.Eval(e) || p.Right.Eval(e)
}
func (p NotPredicate) Eval(e dt.Event) bool {
	return !p.Inner.Eval(e)
}
func (p PathPredicate) Eval(e dt.Event) bool {
    _,err := jsonrw.ExtractFromMap(e.Payload, dt.JSONPath(p.Path))
	return err == nil
}
func (p StrPredicate) Eval(e dt.Event) bool {
    strif,err := jsonrw.ExtractFromMap(e.Payload, dt.JSONPath(p.Path))
    if err != nil {
        return false
    }
    fmt.Println("Comparing",strif, "with", p.Expected)
    return strif.(string) == p.Expected
}
func (p TagPredicate) Eval(e dt.Event) bool {
    return sets.SetIn(p.Tag, e.Channels)
}
func (p TruePredicate) Eval(e dt.Event) bool {
	return true
}
func (p FalsePredicate) Eval(e dt.Event) bool {
	return false
}
func (p NamedPredicate) Eval(e dt.Event) bool {
	//
	// Need to access the mathine to get the body of the predicate
	// or stream and evaluate
	//
	return false
}


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
