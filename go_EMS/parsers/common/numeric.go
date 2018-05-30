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
    visitNumLess(NumLess)
    visitNumLessEq(NumLessEq)
    visitNumEq(NumEq)
    visitNumGreater(NumGreater)
    visitNumGreaterEq(NumGreaterEq)
    visitNumNotEq(NumNotEq)
}

type NumComparison interface {
    Accept (NumComparisonVisitor)
}

type NumLess struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumLess) Accept(visitor NumComparisonVisitor) {
    visitor.visitNumLess(this)
}

type NumLessEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumLessEq) Accept(visitor NumComparisonVisitor) {
    visitor.visitNumLessEq(this)
}

type NumEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumEq) Accept(visitor NumComparisonVisitor) {
    visitor.visitNumEq(this)
}

type NumGreater struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumGreater) Accept(visitor NumComparisonVisitor) {
    visitor.visitNumGreater(this)
}

type NumGreaterEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumGreaterEq) Accept(visitor NumComparisonVisitor) {
    visitor.visitNumGreaterEq(this)
}

type NumNotEq struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumNotEq) Accept(visitor NumComparisonVisitor) {
    visitor.visitNumNotEq(this)
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
    visitIntLiteralExpr(IntLiteralExpr)
    visitFloatLiteralExpr(FloatLiteralExpr)
    visitStreamNameExpr(StreamNameExpr)
    visitNumMulExpr(NumMulExpr)
    visitNumDivExpr(NumDivExpr)
    visitNumPlusExpr(NumPlusExpr)
    visitNumMinusExpr(NumMinusExpr)
    visitIntPathExpr(IntPathExpr)
}
type NumExpr interface {
	Sprint() string
    Accept (NumExprVisitor)
}

type IntPathExpr struct {
	Path dt.JSONPath
}
func (this IntPathExpr) Accept(visitor NumExprVisitor) {
    visitor.visitIntPathExpr(this)
}

type IntLiteralExpr struct {
	Num int
}
func (this IntLiteralExpr) Accept(visitor NumExprVisitor) {
    visitor.visitIntLiteralExpr(this)
}

type FloatLiteralExpr struct {
	Num float32
}
func (this FloatLiteralExpr) Accept(visitor NumExprVisitor) {
    visitor.visitFloatLiteralExpr(this)
}

type StreamNameExpr struct {
	StreamName striverdt.StreamName
}
func (this StreamNameExpr) Accept(visitor NumExprVisitor) {
    visitor.visitStreamNameExpr(this)
}

type NumMulExpr struct {
	Left  NumExpr
	Right NumExpr
}
func (this NumMulExpr) Accept(visitor NumExprVisitor) {
    visitor.visitNumMulExpr(this)
}

type NumDivExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumDivExpr) Accept(visitor NumExprVisitor) {
    visitor.visitNumDivExpr(this)
}

type NumPlusExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumPlusExpr) Accept(visitor NumExprVisitor) {
    visitor.visitNumPlusExpr(this)
}

type NumMinusExpr struct {
	Left NumExpr
	Right NumExpr
}
func (this NumMinusExpr) Accept(visitor NumExprVisitor) {
    visitor.visitNumMinusExpr(this)
}

func NewIntPathExpr(p interface{}) (IntPathExpr) {
	path := p.(PathName).Val
	return IntPathExpr{dt.JSONPath(path)}
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
func NewFloatLiteralExpr(a interface{}) FloatLiteralExpr {
	return FloatLiteralExpr{a.(float32)}
}

func (i IntPathExpr) Sprint() string {
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
	return strconv.FormatFloat(float64(e.Num),'f',4,32)
}
