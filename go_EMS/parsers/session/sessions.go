package session

import(
//	"errors"
	"fmt"
//	"log"
  "strconv"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type Filter struct {
	Pred common.Predicate
	Tag dt.Channel
}

type Filters struct {
	Defs []Filter
}

type Monitor Filters

func newFiltersNode(defs interface{}) (Filters) {
	parsed_defs := common.ToSlice(defs)
	ds := make([]Filter, len(parsed_defs))
	for i,v := range parsed_defs {
		ds[i] = v.(Filter)
	}
	return Filters{ds}
}

// Monitoring Machines

type MoMVisitor interface {
	VisitFilter(Filter)
	VisitSession(Session)
	VisitStream(Stream)
	VisitTrigger(Trigger)
	VisitPredicateDecl(PredicateDecl)
}

type MoM interface {
    Accept(MoMVisitor)
}


//
// Action
//
type EmitAction struct {
  DaTemplate string
	TagName    common.Tag
}
// func (a EmitAction) Sprint() string {
// 	return fmt.Sprintf("emit %s on %s\n",a.StreamName,a.TagName.Tag)
// }

type Trigger struct {
	Pred    common.Predicate
	Action  EmitAction
}

func (this Trigger) Accept(visitor MoMVisitor) {
    visitor.VisitTrigger(this)
}

func newEmitAction(n, t interface{}) (EmitAction) {
	tag  := t.(common.Tag)
	name,ok := n.(common.Identifier)
  var dajson string
  if ok {
    dajson = `{"value": "%`+name.Val+`"}`
  } else {
    dajson =n.(common.BackQuotedString).Val
  }
  return EmitAction{dajson,tag}

}

func newTrigger(p, a interface{}) (Trigger) {
	pred:= p.(common.Predicate)
	act := a.(EmitAction)

	return Trigger{pred,act}
}


type StreamType int
const (
	NumT  StreamType   = iota
	BoolT
	StringT
	LastType = StringT
)

func (e StreamType) Sprint() string {
	switch e {
	case NumT:
		return "num"
	case BoolT:
		return "bool"
	case StringT:
		return "string"
	default:
		return "*error*" // error
	}
}

// func (t StreamType) Sprint() string {
// 
// 	type_names := []string{"int","bool","string"}
// 
// 	// str string 
// 	// switch t {
// 	// case Int:
// 	// 	str = "int"
// 	// case Bool:
// 	// 	str = "bool"
// 	// case String:
// 	// 	str = "string"
// 	// }
// 	// return str
// 
// 	if t>=LastType { return "" }
// 	return fmt.Sprintf("%s",type_names[t])
// }


//
// We need a dictionary of streams (so all streams used are defined)
//
type Streams struct { // a Stream is a Name:=Expr
  DaStreams []Stream
}

type Stream struct { // a Stream is a Name:=Expr
	SType StreamType
	Name striverdt.StreamName
	Expr common.StreamExpr
}

func (this Stream) Accept(visitor MoMVisitor) {
    visitor.VisitStream(this)
}

type Session struct {
	Name  striverdt.StreamName
	Begin common.Predicate
	End   common.Predicate
}

func (this Session) Accept(visitor MoMVisitor) {
    visitor.VisitSession(this)
}


//
// Declaration Node constructors
//

func newStreamDeclaration(ipars,t,n,e interface{}) Streams {
  the_type := t.(StreamType)
  name     := n.(common.Identifier).Val
  expr     := e.(common.StreamExpr)
  var streams []Stream
  if ipars == nil {
    stream := Stream{the_type,striverdt.StreamName(name),expr}
    return Streams{append(streams, stream)}
  }
  pars := ipars.(common.ParamDef)
  fmt.Println(pars)
  for i := int(pars.Fst.Num); i<=int(pars.Lst.Num); i++ {
    newname, newexpr := processParameterizedStream(name, expr, pars.Name, i)
    stream := Stream{the_type,striverdt.StreamName(newname),newexpr}
    streams = append(streams,stream)
  }
  return Streams{streams}
}

func processParameterizedStream(namestr string, expr common.StreamExpr, namepar string, index int) (string, common.StreamExpr) {
  parstr := strconv.Itoa(index)
  return (namestr+"["+parstr+"]"),expr
}

func newSessionDeclaration(n,b,e interface{}) Session {
	name  := n.(common.Identifier).Val
	begin := b.(common.Predicate)
	end   := e.(common.Predicate)
	return Session{striverdt.StreamName(name),begin,end}
}

func newPredicateDeclaration(n,p interface{}) PredicateDecl {
	name := n.(common.Identifier).Val
	pred := p.(common.Predicate)
	return PredicateDecl{name,pred}
}

type PredicateDecl struct {
	Name string
	Pred common.Predicate
}

func (this PredicateDecl) Accept(visitor MoMVisitor) {
    visitor.VisitPredicateDecl(this)
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
	// Additionally, we should check that all
	// used streams are defined and that there is no circularity

	return machine
}

var (
	Verbose bool
)

func Print(mon MonitorMachine) {
	if Verbose {
		fmt.Printf("There are %d stampers\n",len(mon.Stampers))
		fmt.Printf("There are %d sessions\n",len(mon.Sessions))
		fmt.Printf("There are %d streams\n", len(mon.Streams))
		fmt.Printf("There are %d triggers\n", len(mon.Triggers))
		fmt.Printf("There are %d predicates\n", len(mon.Preds))
	}		
		
	for _,v := range mon.Stampers {
		fmt.Printf("when %s do %s\n", v.Pred.Sprint(), v.Tag)
	}
	for _,v := range mon.Sessions {
		fmt.Printf("session %s := [%s, %s]\n",v.Name,v.Begin,v.End)
	}
	for _,v := range mon.Streams {
		fmt.Printf("stream %s %s := %s\n",v.SType.Sprint(),v.Name,v.Expr.Sprint())
	}
	for _,v := range mon.Triggers {
		fmt.Printf("trigger %s do %s\n",v.Pred,v.Action)
	}

}


