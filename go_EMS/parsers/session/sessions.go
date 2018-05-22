package stamp

import(
//	"errors"
	"fmt"
//	"log"
	//	"strconv"
	"strings"
)

var (
	True  TruePredicate
	False FalsePredicate
)

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

type Event struct {
	Payload string // changeme
	Stamp []Tag
}

type Predicate interface {
	Sprint() string
	Eval(e Event) bool
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
// eval(Event e) bool
//
func (p AndPredicate) Eval(e Event) bool {
	return p.Left.Eval(e) && p.Right.Eval(e)
}
func (p OrPredicate) Eval(e Event) bool {
	return p.Left.Eval(e) || p.Right.Eval(e)
}
func (p NotPredicate) Eval(e Event) bool {
	return !p.Inner.Eval(e)
}
func (p PathPredicate) Eval(e Event) bool {
	return true // FIXME
}
func (p StrPredicate) Eval(e Event) bool {
	return true // FIXME
}
func (p TagPredicate) Eval(e Event) bool {
	for _,v := range e.Stamp {
		if strings.Compare(v.Tag,p.Tag)==0 {
			return true
		}
	}
	return false
}
func (p TruePredicate) Eval(e Event) bool {
	return true
}
func (p FalsePredicate) Eval(e Event) bool {
	return false
}




func newAndPredicate(a, b interface{}) (Predicate) {
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

func newOrPredicate(a, b interface{}) (Predicate) {

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
	Tag string
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
	parsed_defs := ToSlice(defs)
	ds := make([]Filter, len(parsed_defs))
	for i,v := range parsed_defs {
		ds[i] = v.(Filter)
	}
	return Filters{ds}
}


func ToSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}


type EmitAction struct {
	StreamName string
	TagName    Tag
}

type Trigger struct {
	Pred    Predicate
	Action  EmitAction
}

type IntPathExpr struct {
	Path string
}


func newEmitAction(n, t interface{}) (EmitAction) {
	name := n.(Identifier).Val
	tag  := t.(Tag)
	
	return EmitAction{name,tag}
}

func newTrigger(p, a interface{}) (Trigger) {
	pred:= p.(Predicate)
	act := a.(EmitAction)

	return Trigger{pred,act} 
}

func newIntPathExpr(p interface{}) (IntPathExpr) {
	path := p.(PathName).Val

	return IntPathExpr{path}
}


type StreamType int
const (
	Int  StreamType   = iota
	Bool     
	String   
)

func (t StreamType) Sprint() string {
	// str string 
	// switch t {
	// case Int:
	// 	str = "int"
	// case Bool:
	// 	str = "bool"
	// case String:
	// 	str = "string"
	// }
	// return str
	return fmt.Sprintf("%s",t)
}


//
// We need a dictionary of streams (so all streams used are defined)
//
type Stream struct { // a Stream is a Name:=Expr
	Type StreamType
	Name string
	Expr StreamExpr 
}
type Session struct {
	Name  string
	Begin Predicate
	End   Predicate
}

type StreamExpr interface {
	// add functions here
	Sprint() string
}

type AggregatorExpr struct {
	Operation string
	Stream    string //StreamName
	Session   string //StreamName
}

type IfExpr struct {
	Antecedent Predicate
	Path       string
}

func (p AggregatorExpr) Sprint() string {
	return fmt.Sprintf("%s(%s within %s)",p.Operation,p.Stream,p.Session)
}

func (p IfExpr) Sprint() string {
	conseq := fmt.Sprintf("e.get(%s)",p.Path)
	return fmt.Sprintf("if %s then %s",p.Antecedent.Sprint(),conseq)
}

func newAggregatorExpr(op, str, ses interface{}) AggregatorExpr {
	operation := op.(string)
	stream    := str.(Identifier).Val
	session   := ses.(Identifier).Val
	
	return AggregatorExpr{operation,stream,session}
}

func newIfExpr(p,e interface{}) IfExpr {
	ante := p.(Predicate)
	path := e.(IntPathExpr).Path
	return IfExpr{ante,path}
}

func newStreamDeclaration(t,n,e interface{}) Stream {
	the_type := t.(StreamType)
	name     := n.(Identifier).Val
	expr     := e.(StreamExpr)
	return Stream{the_type,name,expr}
}



func newSessionDeclaration(n,b,e interface{}) Session {
	name  := n.(Identifier).Val
	begin := b.(Predicate)
	end   := e.(Predicate)
	
	return Session{name,begin,end}
}

type MonitorMachine struct {
	Stampers []Filter
	Sessions []Session
	Streams  []Stream
	Triggers []Trigger
}


//
// after ParseFile returns a []interafce{} all the elements in the slice
//    are Filer or Session or Stream or Trigger
//    this function creates a MonitorMachine from such a mixed slice
//

func ProcessDeclarations(ds []interface{}) MonitorMachine {
	machine := MonitorMachine{}
	for _,v := range ds {
		switch val := v.(type) {
		case Filter:
			machine.Stampers = append(machine.Stampers,val)
		case Session:
			machine.Sessions = append(machine.Sessions,val)
		case Stream:
			machine.Streams  = append(machine.Streams,val)
		case Trigger:
			machine.Triggers = append(machine.Triggers,val)
		}
	}

	// FIXME
	// Additionally, we should checked that all
	// used streams are defined and that there is no circularity
	
	return machine
}

func Print(mon MonitorMachine) {
	fmt.Printf("There are %d stampers\n",len(mon.Stampers))
	fmt.Printf("There are %d sessions\n",len(mon.Sessions))
	fmt.Printf("There are %d streams\n", len(mon.Streams))
	fmt.Printf("There are %d triggers\n", len(mon.Triggers))

	for _,v := range mon.Stampers {
		
		fmt.Printf("when %s do %s\n", v.Pred.Sprint(), v.Tag)
	}
	for _,v := range mon.Sessions {
		fmt.Printf("session %s := (begin=>%s,end=>%s)\n",v.Name,v.Begin.Sprint(),v.End.Sprint())
	}
	for _,v := range mon.Streams {
		fmt.Printf("stream %s %s := %s",v.Type,v.Name,v.Expr.Sprint())
	}
}

