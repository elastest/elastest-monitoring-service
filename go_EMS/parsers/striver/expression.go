package striver

import (
	"fmt"
)


//
// Every expression that can be parsed is an Expr
//    including Predicates, Numeric Expressions, Constants and Streams
//

//
// 1. There are two possibilities for NumExpr and BoolExpr to be Expr
// One is to force every type that implements NumExpr to implement Expr too
// The second is to Embed a NumExpr into a struct that implements Expr
// This second option facilitates the extension of NumExpr to new numeric expressions
// 2. Additionally, ConstExpr and StreamExpr must implement Expr, NumExpr
// and BoolExpr because these are untyped. Also If_Then_Else must also
// implement Expr, NumExpr and BoolExpr  as these can appear everywhere.
//
// We use duck typing for this
//


type Expr interface {
	Sprint() string
	Accept(ExprVisitor)
}
type ExprVisitor interface {
	VisitLetExpr(LetExpr)	
	VisitIfThenElseExpr(IfThenElseExpr)
	// VisitStringExpr(StringExpr)
	VisitNumericExpr(NumericExpr)
	VisitTimeExpr(TimeExpr)
	VisitBoolExpr(BoolExpr)
	VisitStreamOffsetExpr(StreamOffsetExpr)
	VisitOutsideExpr(OutsideExpr)
	VisitNoTickExpr(NoTickExpr)
	VisitConstant(ConstExpr)
}


type ConstExpr struct { // implements Expr,NumExpr,BoolExpr
	Name StreamName
}
type LetExpr struct {
	Name StreamName
	Bind Expr
	Body Expr
}
type IfThenElseExpr struct { // implements Expr,NumExpr,BoolExpr
	If   Expr
	Then Expr
	Else Expr
}
type StreamOffsetExpr struct {  // StreamOffsetExpr implements Expr,NumExpr,BoolExpr
	SExpr StreamExpr
}
type BoolExpr struct {      
	BExpr BooleanExpr
}
type NumericExpr struct {       
	 NExpr NumExpr
}
type TimeExpr struct {
	TExpr Time
}
type OutsideExpr struct {}
type NoTickExpr struct {}

// Accept
func (this ConstExpr) Accept(visitor ExprVisitor) {
	visitor.VisitConstant(this)
}
func (this OutsideExpr) Accept(visitor ExprVisitor) {
	visitor.VisitOutsideExpr(this)
}
func (this LetExpr) Accept(visitor ExprVisitor) {
	visitor.VisitLetExpr(this)
}
func (this IfThenElseExpr) Accept (visitor ExprVisitor) {
	visitor.VisitIfThenElseExpr(this)
}
func (this NumericExpr) Accept (visitor ExprVisitor) {
	visitor.VisitNumericExpr(this)
}
func (this TimeExpr) Accept (visitor ExprVisitor) {
	visitor.VisitTimeExpr(this)
}
func (this StreamOffsetExpr) Accept (visitor ExprVisitor) {
	visitor.VisitStreamOffsetExpr(this)
}
func (this BoolExpr) Accept (visitor ExprVisitor) {
	visitor.VisitBoolExpr(this)
}
func (this NoTickExpr) Accept (visitor ExprVisitor) {
	visitor.VisitNoTickExpr(this)
}

// Sprint()
func (this ConstExpr) Sprint() string {
	return string(this.Name)
}
func (this LetExpr) Sprint() string {
	bind := this.Bind.Sprint()
	body := this.Bind.Sprint()
	return fmt.Sprintf("let %s = %s in %s",this.Name,bind,body)
}
func (this IfThenElseExpr) Sprint() string {
	if_part   := this.If.Sprint()
	then_part := this.Then.Sprint()
	else_part := this.Else.Sprint()
	return fmt.Sprintf("if %s then %s else %s",if_part,then_part,else_part)
}
func (this NumericExpr) Sprint() string {
	return this.NExpr.Sprint()
}
func (this TimeExpr) Sprint() string {
	return this.TExpr.Sprint()
}
func (this StreamOffsetExpr) Sprint() string {
	return this.SExpr.Sprint()
}
func (this BoolExpr) Sprint() string {
	return this.BExpr.Sprint()
}
func (this NoTickExpr) Sprint() string {
	return "notick"
}
func (this OutsideExpr) Sprint() string {
	return "outside"
}

var (
	TheOutsideExpr OutsideExpr
	TheNoTickExpr  NoTickExpr
)

func NewConstExpr(a interface{}) ConstExpr {
	return ConstExpr{getStreamName(a)}
}
func NewNumericExpr(a interface{}) NumericExpr {
	return NumericExpr{a.(NumExpr)}
}
func NewStreamOffsetExpr(a interface{}) StreamOffsetExpr {
	return StreamOffsetExpr{a.(StreamExpr)}
}
func NewTimeExpr(a interface{}) TimeExpr {
	return TimeExpr{a.(Time)}
}
func NewBoolExpr(b interface{}) BoolExpr {
	return BoolExpr{b.(BooleanExpr)}
}
func NewIfThenElseExpr(p,a,b interface{}) IfThenElseExpr {
	return IfThenElseExpr{p.(Expr),a.(Expr),b.(Expr)}
}
func NewLetExpr(n,e,b interface{}) LetExpr {
	name := getStreamName(n)
	return LetExpr{name,e.(Expr),b.(Expr)}
}


//
// StreamExpr : s(~t) s(t~) s(<t) s(t>) s(s<t) ...
//
type StreamExpr interface { 
	AcceptStream(StreamExprVisitor)
	Sprint() string
}

type StreamExprVisitor interface {
	VisitPrevEqValExpr(PrevEqValExpr)
	VisitPrevValExpr(PrevValExpr)
	VisitSuccEqValExpr(SuccEqValExpr)
	VisitSuccValExpr(SuccValExpr)
	VisitStreamFetchExpr(StreamFetchExpr)
}

type PrevEqValExpr struct { //implements StreamExpr
	Name StreamName
	Expr Time
}
type PrevValExpr struct { //implements StreamExpr
	Name StreamName
	Expr Time
}
type SuccEqValExpr struct { //implements StreamExpr
	Name StreamName
	Expr Time
}
type SuccValExpr struct { //implements StreamExpr
	Name StreamName
	Expr Time
}
type StreamFetchExpr struct { //implements StreamExpr
	Name StreamName
	Offset OffsetExpr
}
func NewPrevEqValExpr(s,t interface{}) PrevEqValExpr {
	name := getStreamName(s)
	return PrevEqValExpr{name,t.(Time)}
}
func NewPrevValExpr(s,t interface{}) PrevValExpr {
	name := getStreamName(s)
	return PrevValExpr{name,t.(Time)}
}
func NewSuccEqValExpr(s,t interface{}) SuccEqValExpr {
	name := getStreamName(s)
	return SuccEqValExpr{name,t.(Time)}
}
func NewSuccValExpr(s,t interface{}) SuccValExpr {
	name := getStreamName(s)
	return SuccValExpr{name,t.(Time)}
}
func NewStreamFetchExpr(s,t interface{}) StreamFetchExpr {
	offset := t.(OffsetExpr)
	name := getStreamName(s)
	return StreamFetchExpr{name,offset}
}

func (e PrevEqValExpr) Sprint() string {
	return fmt.Sprintf("%s(~ %s )",e.Name.Sprint(),e.Expr.Sprint())
}
func (e PrevValExpr) Sprint() string {
	return fmt.Sprintf("%s(< %s )",e.Name.Sprint(),e.Expr.Sprint())
}
func (e SuccEqValExpr) Sprint() string {
	return fmt.Sprintf("%s( %s ~)",e.Name.Sprint(),e.Expr.Sprint())
}
func (e SuccValExpr) Sprint() string {
	return fmt.Sprintf("%s( %s >)",e.Name.Sprint(),e.Expr.Sprint())
}
func (e StreamFetchExpr) Sprint() string {
	return fmt.Sprintf("%s(%s)",e.Name.Sprint(),e.Offset.Sprint())
}


func (this PrevEqValExpr) AcceptStream(v StreamExprVisitor) {
	v.VisitPrevEqValExpr(this)
}
func (this PrevValExpr) AcceptStream(v StreamExprVisitor) {
	v.VisitPrevValExpr(this)
}
func (this SuccEqValExpr) AcceptStream(v StreamExprVisitor) {
	v.VisitSuccEqValExpr(this)
}
func (this SuccValExpr) AcceptStream(v StreamExprVisitor) {
	v.VisitSuccValExpr(this)
}
func (this StreamFetchExpr) AcceptStream(v StreamExprVisitor) {
	v.VisitStreamFetchExpr(this)
}

//
// Time
//
type TimeBasic interface { 
	//Accept(...)
	Sprint() string
}
type TimeLiteral struct  { // Implement TimeBasic
	Val float32
}
type TimeConstant struct { // Implement TimeBasic
	Name StreamName     // this is for an symbolic time contant
}
func NewTimeLiteral(n interface{}) TimeLiteral {
	num := n.(FloatLiteralExpr).Num
	return TimeLiteral{num}
}
func NewTimeConstant(s interface{}) TimeConstant {
	k := getStreamName(s)
	return TimeConstant{k}
}

func (l TimeLiteral) Sprint() string {
	return fmt.Sprintf("%f",l.Val)
}
func (c TimeConstant) Sprint() string {
	return string(c.Name)
}


type Time interface {
	Sprint() string
}

type Time_t struct {  } // implements Time
var (
	T Time_t   // "t" as in "define foo t := ..."
)

func (t Time_t) Sprint() string {
	return "t"
}

//
// Offset
//
type OffsetExpr interface { // OffsetExpr "implements" Time
	Sprint() string
}

type PrevEqExpr struct { // implements OffsetExpr
	Name StreamName
	Time Time
}
type PrevExpr struct { // implements OffsetExpr
	Name StreamName
	Time Time
}
type SuccEqExpr struct { // implements OffsetExpr
	Name StreamName
	Time Time
}
type SuccExpr struct { // implements OffsetExpr
	Name StreamName
	Time Time
}

// offsets
func NewPrevEqExpr(n,t interface{}) PrevEqExpr {
	name := getStreamName(n)
	return PrevEqExpr{name,t.(Time)}
}
func NewPrevExpr(n,t interface{}) PrevExpr {
	name := getStreamName(n)
	return PrevExpr{name,t.(Time)}
}
func NewSuccEqExpr(n,t interface{}) SuccEqExpr {
	name := getStreamName(n)
	return SuccEqExpr{name,t.(Time)}
}
func NewSuccExpr(n,t interface{}) SuccExpr {
	name := getStreamName(n)
	return SuccExpr{name,t.(Time)}
}

func (e PrevEqExpr) Sprint() string {
	return fmt.Sprintf("%s<~%s",e.Name.Sprint(),e.Time.Sprint())
}
func (e PrevExpr) Sprint() string {
	return fmt.Sprintf("%s<<%s",e.Name.Sprint(),e.Time.Sprint())
}
func (e SuccEqExpr) Sprint() string {
	return fmt.Sprintf("%s~>%s",e.Name.Sprint(),e.Time.Sprint())
}
func (e SuccExpr) Sprint() string {
	return fmt.Sprintf("%s>>%s",e.Name.Sprint(),e.Time.Sprint())
}

