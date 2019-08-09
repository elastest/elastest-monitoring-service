package impl

import(
  "fmt"
  "encoding/json"
  striverdt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
  parsercommon "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
  dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
  "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
)

var TRIGGER_PREFIX striverdt.StreamName = "trigger::"

func makeTriggerStreams(visitor *MoMToStriverVisitor, predicate parsercommon.Predicate, triggerAction session.EmitAction) {
  var template map[string]interface{} = nil
  templatebytes := []byte(triggerAction.DaTemplate)
  json.Unmarshal(templatebytes, &template)
  extractVals := make(map[string]striverdt.StreamName)
  suffix := generateExtractVals(template, extractVals);
  extractFields := []string{}
  outputSignalName := TRIGGER_PREFIX + striverdt.StreamName(suffix)
  signalNamesVisitor := SignalNamesFromPredicateVisitor{visitor.Preds, []striverdt.StreamName{}, visitor}
  predicate.AcceptPred(&signalNamesVisitor)
  signalNames := signalNamesVisitor.SNames
  predFun := func (args...striverdt.EvPayload) striverdt.EvPayload {
    rawevent := args[0]
    if !rawevent.IsSet {
      panic("No raw event!")
    }
    theEvent := rawevent.Val.(striverdt.EvPayload).Val.(dt.Event)
    snames := args[len(extractVals)+1:]
    argsMap := map[striverdt.StreamName]interface{}{}
    for i,sname := range signalNames {
      if !snames[i].IsSet {
        return striverdt.NothingPayload
      }
      argsMap[sname] = snames[i].Val.(striverdt.EvPayload).Val
    }
    theEvalVisitor := EvalVisitor{false, theEvent, visitor.Preds, argsMap}
    predicate.AcceptPred(&theEvalVisitor)
    thens := args[1:len(extractVals)+1]
    var ret map[string]interface{} = nil
    json.Unmarshal(templatebytes, &ret)
    fmt.Println(string(templatebytes))
    fmt.Println(ret)
    if theEvalVisitor.Result {
      for i, field := range extractFields {
        fmt.Println(i)
        fmt.Println(field)
        then := thens[i]
        if then.IsSet {
          ret[field]=then.Val.(striverdt.EvPayload).Val
        } else {
          delete(ret, field)
        }
      }
      return striverdt.Some(dt.Event{
        dt.ChannelSet{triggerAction.TagName.Tag:nil},
        ret,
        theEvent.Timestamp})
    } else {
      return striverdt.NothingPayload
    }
  }

  argSignals := []striverdt.ValNode{
    &striverdt.PrevEqValNode{striverdt.TNode{}, visitor.InSignalName, []striverdt.Event{}},
  }
  for field,signame := range extractVals {
    argSignals = append(argSignals, &striverdt.PrevEqValNode{striverdt.TNode{}, signame, []striverdt.Event{}})
    extractFields = append(extractFields, field)
  }
  for _,sname := range signalNames {
    argSignals = append(argSignals,
    &striverdt.PrevEqValNode{striverdt.TNode{}, sname, []striverdt.Event{}})
  }

  predVal := striverdt.FuncNode{argSignals, predFun}
  predStream := striverdt.OutStream{outputSignalName, striverdt.SrcTickerNode{visitor.InSignalName}, predVal}
  visitor.OutStreams = append(visitor.OutStreams, predStream)
}

func generateExtractVals(template map[string]interface{}, extractVals map[string]striverdt.StreamName) (ret string) {
  ret = ""
  for field, val := range template {
    s, ok := val.(string)
    if !ok {continue}
    if len(s) < 2 || s[0] != '%' {continue}
    streamName := s[1:]
    ret = ret + streamName
    extractVals[field] = striverdt.StreamName(streamName)
  }
  return
}
