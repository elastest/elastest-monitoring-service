package common

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
//    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "fmt"
	"strconv"
	"errors"
)

//
//  NumComparisons
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
	Accept(NumComparisonVisitor)
	Sprint() string
}

type NumLess struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumLess) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumLess(this)
}
func (this NumLess) Sprint() string {
	return fmt.Sprintf("%s < %s",this.Left.Sprint(),this.Right.Sprint())
}

type NumLessEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumLessEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumLessEq(this)
}
func (this NumLessEq) Sprint() string {
	return fmt.Sprintf("%s <= %s",this.Left.Sprint(),this.Right.Sprint())
}

type NumEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumEq(this)
}
func (this NumEq) Sprint() string {
	return fmt.Sprintf("%s = %s",this.Left.Sprint(),this.Right.Sprint())
}


type NumGreater struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumGreater) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumGreater(this)
}
func (this NumGreater) Sprint() string {
	return fmt.Sprintf("%s > %s",this.Left.Sprint(),this.Right.Sprint())
}


type NumGreaterEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumGreaterEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumGreaterEq(this)
}
func (this NumGreaterEq) Sprint() string {
	return fmt.Sprintf("%s >= %s",this.Left.Sprint(),this.Right.Sprint())
}

type NumNotEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumNotEq) Accept(visitor NumComparisonVisitor) {
    visitor.VisitNumNotEq(this)
}
func (this NumNotEq) Sprint() string {
	return fmt.Sprintf("%s != %s",this.Left.Sprint(),this.Right.Sprint())
}

func NewNumLess(a,b interface{}) NumLess {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumLess{left,right}
}
func NewNumLessEq(a,b interface{}) NumLessEq {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumLessEq{left,right}
}
func NewNumGreater(a,b interface{}) NumGreater {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumGreater{left,right}
}
func NewNumGreaterEq(a,b interface{}) NumGreaterEq {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumGreaterEq{left,right}
}
func NewNumEq(a,b interface{}) NumEq {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumEq{left,right}
}
func NewNumNotEq(a,b interface{}) NumNotEq {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumNotEq{left,right}
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
    AcceptNum(NumExprVisitor)
}

// transforms a StreamNumExpr or StreamNameExpr into a NumExpr
func getNumExpr(e interface{}) (NumExpr,error) {
	if v,ok:=e.(StreamNumExpr);ok {
		return v.Expr,nil
	} else if v,ok:=e.(StreamNameExpr);ok {
		return v,nil
	} else {
		str := fmt.Sprintf("cannot convert to num \"%s\"\n",e.(StreamExpr).Sprint())
		return nil,errors.New(str)
	}
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
	return NumMulExpr{left,right}
}
func (r RightMultExpr) getInner() NumExpr { return r.E }
func (r RightDivExpr) buildBinaryExpr(left NumExpr, right NumExpr) NumExpr {
	return NumDivExpr{left,right}
}
func (r RightDivExpr) getInner() NumExpr { return r.E }

func (r RightPlusExpr) buildBinaryExpr(left NumExpr, right NumExpr) NumExpr {
	return NumPlusExpr{left,right}
}
func (r RightPlusExpr) getInner() NumExpr { return r.E }
func (r RightMinusExpr) buildBinaryExpr(left NumExpr, right NumExpr) NumExpr {
	return NumMinusExpr{left,right}
}
func (r RightMinusExpr) getInner() NumExpr { return r.E }

func NewRightMultExpr(a interface{}) RightMultExpr {
	n,_ := getNumExpr(a)
	return RightMultExpr{n}
}
func NewRightDivExpr(a interface{}) RightDivExpr {
	n,_ := getNumExpr(a)
	return RightDivExpr{n}
}
func NewRightPlusExpr(a interface{}) RightPlusExpr {
	n,_ := getNumExpr(a)
	return RightPlusExpr{n}
}
func NewRightMinusExpr(a interface{}) RightMinusExpr {
	n,_ := getNumExpr(a)
	return RightMinusExpr{n}
}

func Flatten(a,b interface{}) StreamNumExpr {
	exprs := ToSlice(b)
	first,_ := getNumExpr(a)
	if len(exprs)==0 {
		return NewStreamNumExpr(first)
	}
	right := exprs[len(exprs)-1].(RightSubexpr)
	curr  := right.getInner()
	for i := len(exprs)-2; i>=0; i-- {
		left:=  exprs[i].(RightSubexpr)
		curr = right.buildBinaryExpr(left.getInner(),curr)
		right = left
	}
	ret := NewStreamNumExpr(right.buildBinaryExpr(first,curr))
	return ret
}

type NumPathExpr struct {
	Path dt.JSONPath
}
func (this NumPathExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitNumPathExpr(this)
}

type IntLiteralExpr struct {
	Num int
}
func (this IntLiteralExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitIntLiteralExpr(this)
}

type FloatLiteralExpr struct {
	Num float32
}
func (this FloatLiteralExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitFloatLiteralExpr(this)
}

// type StreamNameExpr struct {
//	StreamName striverdt.StreamName
//}
func (this StreamNameExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitStreamNameExpr(this)
}

type NumMulExpr struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumMulExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitNumMulExpr(this)
}

type NumDivExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumDivExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitNumDivExpr(this)
}

type NumPlusExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumPlusExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitNumPlusExpr(this)
}

type NumMinusExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumMinusExpr) AcceptNum(visitor NumExprVisitor) {
    visitor.VisitNumMinusExpr(this)
}

func NewNumPathExpr(p interface{}) (NumPathExpr) {
	path := p.(PathName).Val
	return NumPathExpr{dt.JSONPath(path)}
}
func NewMulExpr(a,b interface{}) NumMulExpr {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumMulExpr{left,right}
}
func NewDivExpr(a,b interface{}) NumDivExpr {
	left ,_ := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumDivExpr{left,right}
}
func NewPlusExpr(a,b interface{}) NumPlusExpr {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumPlusExpr{left,right}
}
func NewMinusExpr(a,b interface{}) NumMinusExpr {
	left,_  := getNumExpr(a)
	right,_ := getNumExpr(b)
	return NumMinusExpr{left,right}
}
//func NewStreamNameExpr(a interface{}) StreamNameExpr {
//	return StreamNameExpr{striverdt.StreamName(a.(Identifier).Val)}
//}
func NewIntLiteralExpr(a interface{}) IntLiteralExpr {
	return IntLiteralExpr{a.(int)}
}
//func NewFloatLiteralExpr(val float64) FloatLiteralExpr {
func NewFloatLiteralExpr(val float64) FloatLiteralExpr {
	return FloatLiteralExpr{float32(val)}
}

func (i NumPathExpr) Sprint() string {
	return fmt.Sprintf("e.getint(%s)",i.Path)
}

func (e NumMulExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),"*",e.Right.Sprint())
}
func (e NumDivExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),"/",e.Right.Sprint())
}
func (e NumPlusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),"+",e.Right.Sprint())
}
func (e NumMinusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)",e.Left.Sprint(),"-",e.Right.Sprint())
}
//func (e StreamNameExpr) Sprint() string {
//	return string(e.Stream)
//}
func (e IntLiteralExpr) Sprint() string {
	return strconv.Itoa(e.Num)
}
func (e FloatLiteralExpr) Sprint() string {
	return "a number"
}
