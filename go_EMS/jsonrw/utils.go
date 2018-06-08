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
	for _,key := range path[:len(path)-1] {
		themap, ok = themap[key].(map[string]interface{})
		if (!ok) {
            return nil, errors.New("Incorrect path")
		}
	}
	ret, ok := themap[path[len(path)-1]]
	if (!ok) {
        return nil, errors.New("Incorrect path")
	}
	return ret,nil
}
