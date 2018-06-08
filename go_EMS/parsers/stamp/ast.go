package stamp

import(
	"fmt"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/common"
)

type Filter struct {
	Pred common.Predicate
	Tag common.Tag
}

type Filters struct {
	Defs []Filter
}

type Monitor Filters

var Verbose bool = false

func Print(mon Monitor) {

	if Verbose {
		fmt.Printf("There are %d stampers\n",len(mon.Defs))
	}
	for _,v := range mon.Defs {
		fmt.Printf("when %s do %s\n", v.Pred.Sprint(), v.Tag.Tag)
	}
}

func newFiltersNode(defs interface{}) (Filters) {
	parsed_defs := common.ToSlice(defs)
	ds := make([]Filter, len(parsed_defs))
	for i,v := range parsed_defs {
		ds[i] = v.(Filter)
	}
	return Filters{ds}
}
