package jsonrw

import (
    "errors"
    "strings"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
  "encoding/json"
    //"fmt"
)

func ExtractFromMap(themap map[string]interface{}, strpath dt.JSONPath) (interface{}, error) {
  path := strings.Split(string(strpath), ".")
  if (len(path) == 0) {
    panic("empty path")
  }
  var ok bool
  var extracted interface{} = themap
  for _,key := range path {
    themap = extracted.(map[string]interface{})
    if len(key) == 0 {
      return nil, errors.New("Incorrect path, empty key")
    }
    if key[len(key)-1:] == "*" {
      key = key[:len(key)-1]
      ok = false
      for k,v := range themap {
        if strings.HasPrefix(k, key) {
          extracted,ok = v, true
        }
      }
    } else {
      extracted, ok = themap[key]
    }
    if (!ok) {
      return nil, errors.New("Incorrect path")
    }
  }
  return extracted,nil
}

func ExtractFromMap2(themap map[string]interface{}, strpaths []dt.JSONPath) (interface{}, error) {
  var extractedif interface{}
  e, _ := json.Marshal(themap)
  extractedif = string(e)
  var err error
  var amap map[string]interface{}
  for _,strpath := range strpaths {
    err = json.Unmarshal([]byte(extractedif.(string)), &amap)
    if err != nil {
      return nil, err
    }
    extractedif, err = ExtractFromMap(amap, strpath)
    if err != nil {
      return nil, err
    }
  }
  return extractedif, nil
}
