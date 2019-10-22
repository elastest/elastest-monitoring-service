package common

import(
    dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
)

type Tag struct {
	Tag dt.Channel
}

type PathName struct {
	Val string
}
type QuotedString struct {
	Val string
}
type BackQuotedString struct {
	Val string
}

type Version struct {
	Num string
}

type Identifier struct{
	Val string
}

func (id Identifier) AcceptComparableStringVisitor(visitor ComparableStringVisitor) {
    visitor.VisitIdentifier(id)
}

type Alphanum struct {
	Val string
}

type Keyword struct {
	Val string
}

type JSONExpr struct {
  paths []dt.JSONPath
}

func NewIdentifier(s string) (Identifier) {
	return Identifier{s}
}
func NewPathName(s string) (PathName) {
	return PathName{s}
}
func NewQuotedString(s string) (QuotedString) {
    s = s[1:len(s)-1]
	return QuotedString{s}
}
func NewBackQuotedString(s string) (BackQuotedString) {
    s = s[1:len(s)-1]
	return BackQuotedString{s}
}

func (qs QuotedString) AcceptComparableStringVisitor(visitor ComparableStringVisitor) {
    visitor.VisitQuotedString(qs)
}

type ComparableStringVisitor interface {
    VisitQuotedString(QuotedString)
    VisitIdentifier(Identifier)
}

type ComparableString interface {
    AcceptComparableStringVisitor(ComparableStringVisitor)
}
