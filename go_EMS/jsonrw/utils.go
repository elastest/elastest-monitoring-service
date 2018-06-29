package jsonrw

import (
    "errors"
    "strings"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
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
