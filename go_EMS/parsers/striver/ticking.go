package striver

import (
	"fmt"
)

type TickingExpr interface {
	// Accept(TickingExprVisitor)
	Sprint() string
}
type TimeConstantTickingExpr struct {
	Val Time
}
type StreamTickingExpr struct {
	Name StreamName
}
type UnionTickingExpr struct {
	First  TickingExpr
	Second TickingExpr
}
type DelayTickingExpr struct {
	Delay DelayExpr
}
type DelayExpr interface {
	Sprint() string
}

type NamedDelayExpr struct {
	Name StreamName
}
type SconstDelayExpr struct {
	Const   Time // either literal or const name
	Carrier StreamName
}

func NewTimeConstantTickingExpr(c interface{}) TimeConstantTickingExpr {
	return TimeConstantTickingExpr{c.(Time)}
}
func NewUnionTickingExpr(a,b interface{}) UnionTickingExpr {
	fst := a.(TickingExpr)
	snd := b.(TickingExpr)
	return UnionTickingExpr{fst,snd}
}
func NewStreamTickingExpr(n interface{}) StreamTickingExpr {
	return StreamTickingExpr{getStreamName(n)}
}
func NewDelayTickingExpr(d interface{}) DelayTickingExpr {
	return DelayTickingExpr{d.(DelayExpr)}
}
func NewSconstDelayExpr(c,s interface{}) SconstDelayExpr {
	t := c.(Time)
	return SconstDelayExpr{t,getStreamName(s)}
}
func NewNamedDelayExpr(n interface{}) NamedDelayExpr {
	
	return NamedDelayExpr{getStreamName(n)}
}

func (e TimeConstantTickingExpr) Sprint() string {
	return fmt.Sprintf("{ %s }",e.Val.Sprint())
}
func (e StreamTickingExpr) Sprint() string {
	return string(e.Name)
}
func (e UnionTickingExpr) Sprint() string {
	str := fmt.Sprintf("%s U %s",e.First.Sprint(),e.Second.Sprint())
	return str
}
func (e DelayTickingExpr) Sprint() string {
	str := fmt.Sprintf("delay %s",e.Delay.Sprint())
	return str
}
func (e NamedDelayExpr) Sprint() string {
	str := fmt.Sprintf("delay %s",string(e.Name))
	return str
}
func (e SconstDelayExpr) Sprint() string {
	str := fmt.Sprintf("sconst %s %s",e.Const.Sprint(),string(e.Carrier))
	return str
}

