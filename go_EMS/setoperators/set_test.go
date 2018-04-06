package setoperators

import (
	"testing"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
)

var setempty = SetFromList([]string{})
var seta = SetFromList([]string{"a"})
var setb = SetFromList([]string{"b"})
var setab = SetFromList([]string{"a","b"})

func TestIn(t *testing.T) {
    if !SetIn(dt.Channel("a"), setab) {
        t.Errorf("a is said not to be in ab")
    }
}

func TestMinus(t *testing.T) {
    if !SetsAreEqual(seta, SetMinus(setab, setb)) {
        t.Errorf("wrong minus: ab-b <> a")
    }
}

func TestEmptyness(t *testing.T) {
    if !SetIsEmpty(setempty) {
        t.Errorf("set is empty, but was not detected")
    }
    if SetIsEmpty(setab) {
        t.Errorf("set is not empty, but was not detected")
    }
}

func TestAdd(t *testing.T) {
    if !SetsAreEqual(setab, SetAdd(seta, dt.Channel("b"))) {
        t.Errorf("wrong add: a+b <> ab")
    }
}

func TestUnion(t *testing.T) {
    if !SetsAreEqual(setab, SetUnion(seta, setb)) {
        t.Errorf("wrong union: aUb <> ab")
    }
}
