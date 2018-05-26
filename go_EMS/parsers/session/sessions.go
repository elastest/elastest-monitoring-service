package session

import(
//	"errors"
	"fmt"
//	"log"
	"strconv"
	"strings"
)

var (
	True  TruePredicate
	False FalsePredicate
)


type Predicate interface {
	Sprint() string
	Eval(e Event) bool
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

type Event struct {
	Payload string // changeme
	Stamp []Tag
}

type NamedPredicate struct {
	// This can be either a predicate "foo" defined with "pred foo :="
	// or a stream named "foo" defined "stream boolean foo :=.."
	Name string
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
func (p NamedPredicate) Sprint() string {
	return p.Name
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

func (p NamedPredicate) Eval(e Event) bool {
	//
	// Need to access the mathine to get the body of the predicate
	// or stream and evaluate
	//
	return false
}

func newNamedPredicate(n interface{}) NamedPredicate {
	name := n.(Identifier).Val
	return NamedPredicate{name}
}

func newAndPredicate(a, b interface{}) Predicate {
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

func newOrPredicate(a, b interface{}) Predicate {

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

func newNotPredicate(p interface{}) NotPredicate {
	return NotPredicate{p.(Predicate)}
}

func newPathPredicate(p interface{}) PathPredicate {
	path := p.(PathName).Val
	return PathPredicate{path}
}

func newStrPredicate(p,v interface{}) StrPredicate {
	path     :=p.(PathName).Val
	expected :=v.(QuotedString).Val
	return StrPredicate{path,expected}
}
func newTagPredicate(t interface{}) TagPredicate {
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

//
// Action
//
type EmitAction struct {
	StreamName string
	TagName    Tag
}
func (a EmitAction) Sprint() string {
	return fmt.Sprintf("emit %s on %s\n",a.StreamName,a.TagName.Tag)
}

type Trigger struct {
	Pred    Predicate
	Action  EmitAction
}

type IntPathExpr struct {
	Path string
}
type StringPathExpr struct {
	Path string
}
func (i IntPathExpr) Sprint() string {
	return fmt.Sprintf("e.getint(%s)",i.Path)
}
func (i StringPathExpr) Sprint() string {
	return fmt.Sprintf("e.getstr(%s)",i.Path)
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

func newStringPathExpr(p interface{}) (StringPathExpr) {
	path := p.(PathName).Val

	return StringPathExpr{path}
}


type StreamType int
const (
	IntT  StreamType   = iota
	BoolT
	StringT
	LastType = StringT
)

func (t StreamType) Sprint() string {

	type_names := []string{"int","bool","string"}

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

	if t>=LastType { return "" }
	return fmt.Sprintf("%s",type_names[t])
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

//
// Expressions
//
type StreamExpr interface {
	// add functions here
	Sprint() string
}

type AggregatorExpr struct {
	Operation string
	Stream    string //StreamName
	Session   string //StreamName
}
type IfThenExpr struct {
	If   Predicate
	Then StreamExpr
}
type IfThenElseExpr struct {
	If   Predicate
	Then StreamExpr
	Else StreamExpr
}
type StringExpr struct {
	Path string// so far only e.get(path) claiming to return a string
}
type NumExpr interface { // See cases below
	Sprint() string
}
type PredExpr struct {
	Pred Predicate 
}
func (p AggregatorExpr) Sprint() string {
	return fmt.Sprintf("%s(%s within %s)",p.Operation,p.Stream,p.Session)
}

func (p IfThenExpr) Sprint() string {
	return fmt.Sprintf("if %s then %s",p.If.Sprint(),p.Then.Sprint())
}
func (p IfThenElseExpr) Sprint() string {
	return fmt.Sprintf("if %s then %s else %s",p.If.Sprint(),p.Then.Sprint(),p.Else.Sprint())
}

func (p PredExpr) Sprint() string {
	return p.Pred.Sprint()
}


//
// Expression Node constructors
//
func newAggregatorExpr(op, str, ses interface{}) AggregatorExpr {
	operation := op.(string)
	stream    := str.(Identifier).Val
	session   := ses.(Identifier).Val
	
	return AggregatorExpr{operation,stream,session}
}

func newIfThenExpr(p,e interface{}) IfThenExpr {
	if_part   := p.(Predicate)
	then_part := e.(StreamExpr)
	return IfThenExpr{if_part,then_part}
}
func newIfThenElseExpr(p,a,b interface{}) IfThenElseExpr {
	if_part   := p.(Predicate)
	then_part := a.(StreamExpr)
	else_part := b.(StreamExpr)
	return IfThenElseExpr{if_part, then_part, else_part}
}
func newPredExpr(p interface{}) PredExpr {
	return PredExpr{p.(Predicate)}
}

//
// Numeric:
//  NumExpressions and NumComparison
//
type NumComparison interface {
	Sprint() string
//	Eval() bool // miss context to perform the evaluation
}

type NumLess struct {
	Left  NumExpr
	Right NumExpr
}

type NumLessEq struct {
	Left  NumExpr
	Right NumExpr
}

type NumEq struct {
	Left  NumExpr
	Right NumExpr
}

type NumGreater struct {
	Left  NumExpr
	Right NumExpr
}

type NumGreaterEq struct {
	Left  NumExpr
	Right NumExpr
}

type NumNotEq struct {
	Left  NumExpr
	Right NumExpr
}

func newNumLess(a,b interface{}) NumLess {
	return NumLess{a.(NumExpr),b.(NumExpr)}
}
func newNumLessEq(a,b interface{}) NumLessEq {
	return NumLessEq{a.(NumExpr),b.(NumExpr)}
}
func newNumGreater(a,b interface{}) NumGreater {
	return NumGreater{a.(NumExpr),b.(NumExpr)}
}
func newNumGreaterEq(a,b interface{}) NumGreaterEq {
	return NumGreaterEq{a.(NumExpr),b.(NumExpr)}
}
func newNumEq(a,b interface{}) NumEq {
	return NumEq{a.(NumExpr),b.(NumExpr)}
}
func newNumNotEq(a,b interface{}) NumNotEq {
	return NumNotEq{a.(NumExpr),b.(NumExpr)}
}

//
// Numeric Expressions
// 
type IntLiteralExpr struct {
	Num int
}
type FloatLiteralExpr struct {
	Num float32
}
type NumStreamExpr struct {
	StreamName string
}
type NumMulExpr struct {
	Left  NumExpr
	Right NumExpr
}
type NumDivExpr struct {
	Left NumExpr
	Right NumExpr
}
type NumPlusExpr struct {
	Left NumExpr
	Right NumExpr
}
type NumMinusExpr struct {
	Left NumExpr
	Right NumExpr
}

func newMulExpr(a,b interface{}) NumMulExpr {
	return NumMulExpr{a.(NumExpr),b.(NumExpr)}
}
func newDivExpr(a,b interface{}) NumDivExpr {
	return NumDivExpr{a.(NumExpr),b.(NumExpr)}
}
func newPlusExpr(a,b interface{}) NumPlusExpr {
	return NumPlusExpr{a.(NumExpr),b.(NumExpr)}
}
func newMinusExpr(a,b interface{}) NumMinusExpr {
	return NumMinusExpr{a.(NumExpr),b.(NumExpr)}
}
func newNumStreamExpr(a interface{}) NumStreamExpr {
	return NumStreamExpr{a.(Identifier).Val}
}
func newIntLiteralExpr(a interface{}) IntLiteralExpr {
	return IntLiteralExpr{a.(int)}
}
func newFloatLiteralExpr(a interface{}) FloatLiteralExpr {
	return FloatLiteralExpr{a.(float32)}
}

func (e NumMulExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),'*',e.Right.Sprint())
}
func (e NumDivExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),'/',e.Right.Sprint())
}
func (e NumPlusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),'+',e.Right.Sprint())
}
func (e NumMinusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),'-',e.Right.Sprint())
}
func (e NumStreamExpr) Sprint() string {
	return e.StreamName
}
func (e IntLiteralExpr) Sprint() string {
	return strconv.Itoa(e.Num)
}
func (e FloatLiteralExpr) Sprint() string {
	return strconv.FormatFloat(float64(e.Num),'f',4,32)
}

//
// Declaration Node constructors
//

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

func newPredicateDeclaration(n,p interface{}) PredicateDecl {
	name := n.(Identifier).Val
	pred := p.(Predicate)

	return PredicateDecl{name,pred}
}

type PredicateDecl struct {
	Name string
	Pred Predicate 
}

type MonitorMachine struct {
	Stampers []Filter
	Sessions []Session
	Streams  []Stream
	Triggers []Trigger
	Preds    []PredicateDecl
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
		case PredicateDecl:
			machine.Preds    = append(machine.Preds,val)
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
	fmt.Printf("There are %d predicates\n", len(mon.Preds))

	for _,v := range mon.Stampers {
		
		fmt.Printf("when %s do %s\n", v.Pred.Sprint(), v.Tag)
	}
	for _,v := range mon.Sessions {
		fmt.Printf("session %s := (begin=>%s,end=>%s)\n",v.Name,v.Begin.Sprint(),v.End.Sprint())
	}
	for _,v := range mon.Streams {
		fmt.Printf("stream %s %s := %s\n",v.Type.Sprint(),v.Name,v.Expr.Sprint())
	}
	for _,v := range mon.Triggers {
		fmt.Printf("trigger %s do %s\n",v.Pred.Sprint(),v.Action.Sprint())
	}

}


