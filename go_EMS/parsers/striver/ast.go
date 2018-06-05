package striver
import(
//	"errors"
	"fmt"
//	"log"
//	"strconv"
//	"strings"
)


type StreamType int
const (
	NumT  StreamType   = iota
	BoolT
	StringT
	Unknown	// we use this in the parser for unknow type values (offset expressions) that will be resolved later
	LastType = StringT
)

type StreamName string

func (s StreamName) Sprint() string {
	return string(s)
}

func (t StreamType) Sprint() string {

	type_names := []string{"num","bool","string"}

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


type ConstDecl struct {   	// const int one_sec := 1s  
	Name StreamName
	Type StreamType
	Val Expr
}
type InputDecl struct {        // input int bar 
	Name StreamName
	Type StreamType
}
type OutputDecl struct {  // output int foo /* this is just a decl, later a tick and a def will be given */
	Name StreamName
	Type StreamType
}

type TicksDecl struct {
	Name StreamName
	Ticks TickingExpr
}

type OutputDefinition struct {
	Name StreamName
	Type StreamType
	Expr Expr                // chango to ValueExpr?
}

func NewConstDecl(n,t,e interface{}) ConstDecl {
	name := getStreamName(n)
	return ConstDecl{name,t.(StreamType),e.(Expr)}
}
func NewInputDecl(n,t interface{}) InputDecl {
	name := getStreamName(n)
	return InputDecl{name,t.(StreamType)}
}
func NewOutputDecl(n,t interface{}) OutputDecl {
	name := getStreamName(n)
	return OutputDecl{name,t.(StreamType)}
}
func NewTicksDecl(n,t interface{}) TicksDecl {
	name := getStreamName(n)
	expr := t.(TickingExpr)
	return TicksDecl{name,expr}
}
func NewOutputDefinition(n,t,e interface{}) OutputDefinition {
	name := getStreamName(n)
	expr := e.(Expr)
	return OutputDefinition{name,t.(StreamType),expr}
}

func getStreamName(a interface{}) StreamName {
	return StreamName(a.(Identifier).Val)
}

//
// DEPRECATED (MOVED ELSEWHERE)
//
// type Event struct {
// 	Payload string // changeme
// //	Stamp []Tag
// }
// //
// // eval(Event e) bool
// //
// func (p AndPredicate) Eval(e Event) bool {
// 	return p.Left.Eval(e) && p.Right.Eval(e)
// }
// func (p OrPredicate) Eval(e Event) bool {
// 	return p.Left.Eval(e) || p.Right.Eval(e)
// }
// func (p NotPredicate) Eval(e Event) bool {
// 	return !p.Inner.Eval(e)
// }
// func (p TruePredicate) Eval(e Event) bool {
// 	return true
// }
// func (p FalsePredicate) Eval(e Event) bool {
// 	return false
// }

//type Monitor Filters

type Tag struct {
	//	Tag dt.Channel
	Tag string
}


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

func NewIdentifier(s string) (Identifier) {
	return Identifier{s}
}
func NewPathName(s string) (PathName) {
	return PathName{s}
}
func NewQuotedString(s string) (QuotedString) {
	return QuotedString{s}
}

func ToSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}





