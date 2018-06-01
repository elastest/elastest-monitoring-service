package common

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "fmt"
    "strconv"
)

//
// Numeric:
//  NumExpressions and NumComparison
//
type NumComparisonVisitor interface {
    VisitNumLess(NumLess)
    VisitNumLessEq(NumLessEq)
    VisitNumEq(NumEq)
    VisitNumGreater(NumGreater)
    VisitNumGreaterEq(NumGreaterEq)
    VisitNumNotEq(NumNotEq)
}

type NumComparison interface {
    Accept (NumComparisonVisitor)
}

type NumLess struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumLess) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumLess(this)
}

type NumLessEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumLessEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumLessEq(this)
}

type NumEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumEq(this)
}

type NumGreater struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumGreater) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumGreater(this)
}

type NumGreaterEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumGreaterEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumGreaterEq(this)
}

type NumNotEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumNotEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumNotEq(this)
}

func NewNumLess(a,b interface{}) NumLess {
	return NumLess{a.(NumExpr),b.(NumExpr)}
}
func NewNumLessEq(a,b interface{}) NumLessEq {
	return NumLessEq{a.(NumExpr),b.(NumExpr)}
}
func NewNumGreater(a,b interface{}) NumGreater {
	return NumGreater{a.(NumExpr),b.(NumExpr)}
}
func NewNumGreaterEq(a,b interface{}) NumGreaterEq {
	return NumGreaterEq{a.(NumExpr),b.(NumExpr)}
}
func NewNumEq(a,b interface{}) NumEq {
	return NumEq{a.(NumExpr),b.(NumExpr)}
}
func NewNumNotEq(a,b interface{}) NumNotEq {
	return NumNotEq{a.(NumExpr),b.(NumExpr)}
}

//
// Numeric Expressions
// 

type NumExprVisitor interface {
    VisitIntLiteralExpr(IntLiteralExpr)
    VisitFloatLiteralExpr(FloatLiteralExpr)
    VisitStreamNameExpr(StreamNameExpr)
    VisitNumMulExpr(NumMulExpr)
    VisitNumDivExpr(NumDivExpr)
    VisitNumPlusExpr(NumPlusExpr)
    VisitNumMinusExpr(NumMinusExpr)
    VisitNumPathExpr(NumPathExpr)
}
type NumExpr interface {
    Sprint() string
    Accept (NumExprVisitor)
}

//
// Flatten ExpressionLists
//
type RightSubexpr interface{
	buildBinaryExpr(left NumExpr, right NumExpr) NumExpr
	getInner() NumExpr
}
type RightMultExpr   struct { E NumExpr }
type RightDivExpr    struct { E NumExpr }
type RightPlusExpr   struct { E NumExpr }
type RightMinusExpr  struct { E NumExpr }

func (r RightMultExpr) buildBinaryExpr(left NumExpr, right NumExpr) NumExpr {
	return NewMulExpr(left,right)
}
func (r RightMultExpr) getInner() NumExpr { return r.E }
func (r RightDivExpr) buildBinaryExpr(left NumExpr, right NumExpr) NumExpr {
	return NewDivExpr(left,right)
}
func (r RightDivExpr) getInner() NumExpr { return r.E }

func (r RightPlusExpr) buildBinaryExpr(left NumExpr, right NumExpr) NumExpr {
	return NewPlusExpr(left,right)
}
func (r RightPlusExpr) getInner() NumExpr { return r.E }
func (r RightMinusExpr) buildBinaryExpr(left NumExpr, right NumExpr) NumExpr {
	return NewMinusExpr(left,right)
}
func (r RightMinusExpr) getInner() NumExpr { return r.E }


func NewRightMultExpr(a interface{}) RightMultExpr {
	return RightMultExpr{a.(NumExpr)}
}
func NewRightDivExpr(a interface{}) RightDivExpr {
	return RightDivExpr{a.(NumExpr)}
}
func NewRightPlusExpr(a interface{}) RightPlusExpr {
	return RightPlusExpr{a.(NumExpr)}
}
func NewRightMinusExpr(a interface{}) RightMinusExpr {
	return RightMinusExpr{a.(NumExpr)}
}

func Flatten(a,b interface{}) NumExpr {
	exprs := ToSlice(b)
	first := a.(NumExpr)
	if len(exprs)==0 {
		return first
	}
	right := exprs[len(exprs)-1].(RightSubexpr)
	curr  := right.getInner()
	for i := len(exprs)-2; i>=0; i-- {
		left:=  exprs[i].(RightSubexpr)
		curr = right.buildBinaryExpr(left.getInner(),curr)
		right = left
	}
	ret := right.buildBinaryExpr(first,curr)
	return ret
}

type NumPathExpr struct {
	Path dt.JSONPath
}
func (this NumPathExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitNumPathExpr(this)
}

type IntLiteralExpr struct {
	Num int
}
func (this IntLiteralExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitIntLiteralExpr(this)
}

type FloatLiteralExpr struct {
	Num float32
}
func (this FloatLiteralExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitFloatLiteralExpr(this)
}

type StreamNameExpr struct {
	StreamName striverdt.StreamName
}
func (this StreamNameExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitStreamNameExpr(this)
}

type NumMulExpr struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumMulExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitNumMulExpr(this)
}

type NumDivExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumDivExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitNumDivExpr(this)
}

type NumPlusExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumPlusExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitNumPlusExpr(this)
}

type NumMinusExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumMinusExpr) Accept(visitor NumExprVisitor) {
    visitor.VisitNumMinusExpr(this)
}

func NewNumPathExpr(p interface{}) (NumPathExpr) {
	path := p.(PathName).Val
	return NumPathExpr{dt.JSONPath(path)}
}
func NewMulExpr(a,b interface{}) NumMulExpr {
	return NumMulExpr{a.(NumExpr),b.(NumExpr)}
}
func NewDivExpr(a,b interface{}) NumDivExpr {
	return NumDivExpr{a.(NumExpr),b.(NumExpr)}
}
func NewPlusExpr(a,b interface{}) NumPlusExpr {
	return NumPlusExpr{a.(NumExpr),b.(NumExpr)}
}
func NewMinusExpr(a,b interface{}) NumMinusExpr {
	return NumMinusExpr{a.(NumExpr),b.(NumExpr)}
}
func NewStreamNameExpr(a interface{}) StreamNameExpr {
	return StreamNameExpr{striverdt.StreamName(a.(Identifier).Val)}
}
func NewIntLiteralExpr(a interface{}) IntLiteralExpr {
	return IntLiteralExpr{a.(int)}
}
func NewFloatLiteralExpr(val float64) FloatLiteralExpr {
	return FloatLiteralExpr{float32(val)}
}

func (i NumPathExpr) Sprint() string {
	return fmt.Sprintf("e.getint(%s)",i.Path)
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
func (e StreamNameExpr) Sprint() string {
	return string(e.StreamName)
}
func (e IntLiteralExpr) Sprint() string {
	return strconv.Itoa(e.Num)
}
func (e FloatLiteralExpr) Sprint() string {
	return "a number"
}
