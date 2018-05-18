package signals

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    //"fmt"
    striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
)

type SampledSignalDefinition struct {
    Name striverdt.StreamName
    InChannel dt.Channel
    ValuePath dt.JSONPath
}

type FuncSignalDefinition struct {
    Name striverdt.StreamName
    SourceName striverdt.StreamName
    FuncDef FunctionDefinition
}

type ConditionalSignalDefinition struct {
    Name striverdt.StreamName
    SourceSignal striverdt.StreamName
    Condition striverdt.StreamName
}

type WriteValue struct {
    Valtype string
    Valpayload string
}

type SignalWriteDefinition struct {
    SourceSignal striverdt.StreamName
    OutChannel dt.Channel
    //FieldsAndValues map[dt.JSONPath]dt.WriteValue
    // This is stringly typed. Actually, value could be either a Param, a
    // literal or "value". Maybe every Param should be present at least once.
}

// SignalDefVisitor

type SignalDefVisitor interface {
    visitSampled(SampledSignalDefinition)
    visitFuncSignal(FuncSignalDefinition)
}

type SignalDefinition interface {
    Accept(SignalDefVisitor)
}

func (this SampledSignalDefinition) Accept(visitor SignalDefVisitor) {
    visitor.visitSampled(this)
}

func (this FuncSignalDefinition) Accept(visitor SignalDefVisitor) {
    visitor.visitFuncSignal(this)
}

// FunDefs

type FunctionDefinition interface {
    getFunction() (func(interface{}) interface{})
}

type SignalEqualsLiteral struct {
    Literal string
}

func (def SignalEqualsLiteral) getFunction() (func(interface{}) interface{}) {
    return func(arg interface{}) interface{} { return arg.(string) == def.Literal}
}
