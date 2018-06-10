package striver

import (
	"fmt"
	"strconv"
	"errors"
)

//
// NumComparison
//
type NumComparison interface {
	Sprint() string
	AcceptNumComp(NumComparisonVisitor)
}

type NumComparisonVisitor interface {
    VisitNumLess(NumLess)
    VisitNumLessEq(NumLessEq)
    VisitNumEq(NumEq)
    VisitNumGreater(NumGreater)
    VisitNumGreaterEq(NumGreaterEq)
    VisitNumNotEq(NumNotEq)
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

func (this NumLess) AcceptNumComp(visitor NumComparisonVisitor) {
    visitor.VisitNumLess(this)
}
func (this NumLessEq) AcceptNumComp(visitor NumComparisonVisitor) {
    visitor.VisitNumLessEq(this)
}
func (this NumGreater) AcceptNumComp(visitor NumComparisonVisitor) {
    visitor.VisitNumGreater(this)
}
func (this NumGreaterEq) AcceptNumComp(visitor NumComparisonVisitor) {
    visitor.VisitNumGreaterEq(this)
}
func (this NumEq) AcceptNumComp(visitor NumComparisonVisitor) {
    visitor.VisitNumEq(this)
}
func (this NumNotEq) AcceptNumComp(visitor NumComparisonVisitor) {
    visitor.VisitNumNotEq(this)
}


func (this NumLess) Sprint() string {
	return fmt.Sprintf("%s < %s",this.Left.Sprint(),this.Right.Sprint())
}
func (this NumLessEq) Sprint() string {
	return fmt.Sprintf("%s <= %s",this.Left.Sprint(),this.Right.Sprint())
}
func (this NumGreater) Sprint() string {
	return fmt.Sprintf("%s > %s",this.Left.Sprint(),this.Right.Sprint())
}
func (this NumGreaterEq) Sprint() string {
	return fmt.Sprintf("%s >= %s",this.Left.Sprint(),this.Right.Sprint())
}
func (this NumNotEq) Sprint() string {
	return fmt.Sprintf("%s != %s",this.Left.Sprint(),this.Right.Sprint())
}
func (this NumEq) Sprint() string {
	return fmt.Sprintf("%s = %s",this.Left.Sprint(),this.Right.Sprint())
}


//
// Numeric Expressions
//
type NumExpr interface{ 
	AcceptNum(NumExprVisitor)
	Sprint() string
}

type NumExprVisitor interface {
	VisitIntLiteralExpr(IntLiteralExpr) //
	VisitFloatLiteralExpr(FloatLiteralExpr) //
	VisitStreamOffsetExpr(StreamOffsetExpr)
	VisitConst(ConstExpr)
	VisitNumMulExpr(NumMulExpr) //
	VisitNumDivExpr(NumDivExpr) //
	VisitNumPlusExpr(NumPlusExpr) //
	VisitNumMinusExpr(NumMinusExpr) //
//	VisitNumPathExpr(NumPathExpr)
}


type IntLiteralExpr struct {
	Num int
}
type FloatLiteralExpr struct {
		Num float32
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
func NewIntLiteralExpr(a interface{}) IntLiteralExpr {
	return IntLiteralExpr{a.(int)}
}
func NewFloatLiteralExpr(a interface{}) FloatLiteralExpr {
	//	return FloatLiteralExpr{a.(float32)}
	return FloatLiteralExpr{float32(a.(float64))}
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
func (e IntLiteralExpr) Sprint() string {
	return strconv.Itoa(e.Num)
}
func (e FloatLiteralExpr) Sprint() string {
	return strconv.FormatFloat(float64(e.Num),'f',4,32)
}


func (e NumMulExpr) AcceptNum(v NumExprVisitor) {
	v.VisitNumMulExpr(e)
}
func (e NumDivExpr) AcceptNum(v NumExprVisitor) {
	v.VisitNumDivExpr(e)
}
func (e NumPlusExpr) AcceptNum(v NumExprVisitor) {
	v.VisitNumPlusExpr(e)
}
func (e NumMinusExpr) AcceptNum(v NumExprVisitor) {
	v.VisitNumMinusExpr(e)
}
func (e IntLiteralExpr) AcceptNum(v NumExprVisitor) {
	v.VisitIntLiteralExpr(e)
}
func (e FloatLiteralExpr) AcceptNum(v NumExprVisitor) {
	v.VisitFloatLiteralExpr(e)
}
func (e StreamOffsetExpr) AcceptNum(v NumExprVisitor) {
	v.VisitStreamOffsetExpr(e)
}	
func (e ConstExpr) AcceptNum(v NumExprVisitor) {
	v.VisitConst(e)
}


//
//
//
func getNumExpr(e interface{}) (NumExpr,error) {
	if v,ok:=e.(NumericExpr);ok {
		return v.NExpr,nil
	} else if v,ok:=e.(StreamOffsetExpr);ok {
		return v,nil
	} else if k,ok := e.(ConstExpr); ok {
		return k,nil
	} else {
		str := fmt.Sprintf("cannot convert to num \"%s\"\n",e.(Expr).Sprint())
		return nil,errors.New(str)
	}
}
func NumExprToExpr(expr NumExpr) Expr {
	if s,ok:=expr.(StreamOffsetExpr); ok {
		return s
	} else if k,ok:=expr.(ConstExpr); ok {
		return k
	} else {
		return NewNumericExpr(expr)
	}
}

// FLATTEN
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
func Flatten(a,b interface{}) Expr {
	exprs := ToSlice(b)
	first,_ := getNumExpr(a)
	if len(exprs)==0 {
		return NumExprToExpr(first)
	}
	right := exprs[len(exprs)-1].(RightSubexpr)
	curr  := right.getInner()
	for i := len(exprs)-2; i>=0; i-- {
		left:=  exprs[i].(RightSubexpr)
		curr = right.buildBinaryExpr(left.getInner(),curr)
		right = left
	}
	ret := NumExprToExpr(right.buildBinaryExpr(first,curr))
	return ret
}
