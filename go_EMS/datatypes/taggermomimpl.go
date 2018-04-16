package data

import (
    "errors"
)

func (andNode AndNode) Eval(payload map[string]interface{}) bool {
    for _,node := range andNode.Members {
        if !node.Eval(payload) {
            return false
        }
    }
    return true
}

func (orNode OrNode) Eval(payload map[string]interface{}) bool {
    for _,node := range orNode.Members {
        if node.Eval(payload) {
            return true
        }
    }
    return false
}

func (pathfunNode PathFunNode) Eval(payload map[string]interface{}) bool {
    theVal, err := extractFromMap(payload, pathfunNode.Path)
    if err != nil {
        return false
    }
    return pathfunNode.Fun(theVal)
}

func extractFromMap(themap map[string]interface{}, path []string) (interface{}, error) {
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
