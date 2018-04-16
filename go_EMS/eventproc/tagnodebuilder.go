package eventproc

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    "errors"
    "strings"
)

func getNodeFromFilter(thejson map[string]interface{}) (dt.TagNode,error) {
    nodetype, ok := thejson["nodetype"].(string)
    if !ok {
        err := errors.New("No field 'nodetype' present in filter definition")
        return nil,err
    }
    switch nodetype {
    case "AND":
        return getAndNode(thejson)
    case "OR":
        return getOrNode(thejson)
    case "PATHFUN":
        return getPathFunNode(thejson)
    }
    return nil,errors.New("Unrecognized nodetype "+nodetype+" in filter definition")
}

func getAndNode(thejson map[string]interface{}) (dt.TagNode,error) {
    members, ok := thejson["members"].([]map[string]interface{})
    if !ok {
        err := errors.New("No valid field 'members' in AND node definition")
        return nil,err
    }
    memberNodes := make([]dt.TagNode, len(members))
    for i,member := range members {
        memberNode, err := getNodeFromFilter(member)
        if err!=nil {
            return nil, err
        }
        memberNodes[i] = memberNode
    }
    return dt.AndNode{memberNodes},nil
}

func getOrNode(thejson map[string]interface{}) (dt.TagNode,error) {
    members, ok := thejson["members"].([]map[string]interface{})
    if !ok {
        err := errors.New("No valid field 'members' in OR node definition")
        return nil,err
    }
    memberNodes := make([]dt.TagNode, len(members))
    for i,member := range members {
        memberNode, err := getNodeFromFilter(member)
        if err!=nil {
            return nil, err
        }
        memberNodes[i] = memberNode
    }
    return dt.OrNode{memberNodes},nil
}

func getPathFunNode(thejson map[string]interface{}) (dt.TagNode,error) {
    funName, ok := thejson["function"].(string)
    if !ok {
        err := errors.New("No valid field 'function' in PATHFUN node definition")
        return nil,err
    }
    strpath, ok := thejson["path"].(string)
    if !ok {
        err := errors.New("No valid field 'path' in PATHFUN node definition")
        return nil,err
    }
	path := strings.Split(strpath, ".")
	if (len(path) == 0) {
		return nil, errors.New("Empty path in PATHFUN node definition")
	}
    switch funName {
    case "exists":
        return dt.PathFunNode{path, func(in interface{})bool{return true}},nil
    }
    return nil,errors.New("Unrecognized function "+funName+" in funpath definition")
}
