package stamp

import(
	"fmt"
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
)

var (
	True  TruePredicate
	False FalsePredicate
)

func Print(mon Monitor) {
	fmt.Printf("There are %d stampers\n",len(mon.Defs))
	for _,v := range mon.Defs {
		//fmt.Printf("when %s do %s\n", v.pred.pred, v.tag.tag)
		if (v == v) {}
	}
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
	Tag string
}

type Predicate interface {
	Sprint() string
	Eval(e dt.Event) bool
}

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
	return true // FIXME
}
func (p StrPredicate) Eval(e dt.Event) bool {
	return true // FIXME
}
func (p TagPredicate) Eval(e dt.Event) bool {
    return sets.SetIn(dt.Channel(p.Tag), e.Channels)
}
func (p TruePredicate) Eval(e dt.Event) bool {
	return true
}
func (p FalsePredicate) Eval(e dt.Event) bool {
	return false
}




func newAndPredicate(a, b interface{}) (Predicate) {
	preds := toSlice(b)
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

func newOrPredicate(a, b interface{}) (Predicate) {

	preds := toSlice(b)
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

func newNotPredicate(p interface{}) (NotPredicate) {
	return NotPredicate{p.(Predicate)}
}

func newPathPredicate(p interface{}) (PathPredicate) {
	path := p.(PathName).Val
	return PathPredicate{path}
}

func newStrPredicate(p,v interface{}) (StrPredicate) {
	path     :=p.(PathName).Val
	expected :=v.(QuotedString).Val
	return StrPredicate{path,expected}
}
func newTagPredicate(t interface{}) (TagPredicate) {
	tag := t.(Tag).Tag
	return TagPredicate{tag}
}

type Tag struct {
	Tag string
}


type Filter struct {
	Pred Predicate
	Tag Tag
}

type Filters struct {
	Defs []Filter
}

type Version struct {
	Num string
}

type Monitor Filters

type Identifier struct{
	Val string
}

type PathName struct {
	Val string
}
type QuotedString struct {
	Val string
}

type Alphanum struct {
	Val string
}

type Keyword struct {
	Val string
}

func newIdentifier(s string) (Identifier) {
	return Identifier{s}
}
func newPathName(s string) (PathName) {
	return PathName{s}
}
func newQuotedString(s string) (QuotedString) {
	return QuotedString{s}
}

func newFiltersNode(defs interface{}) (Filters) {
	parsed_defs := toSlice(defs)
	ds := make([]Filter, len(parsed_defs))
	for i,v := range parsed_defs {
		ds[i] = v.(Filter)
	}
	return Filters{ds}
}


func toSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
