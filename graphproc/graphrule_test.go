package graphproc

import (
	"testing"
)

func TestRule(t *testing.T) {
	cond := NewCondition("a == 5")
	args := map[string]any{"a": 5}
	result, err := cond.Rule(args)
	if !result || err != nil {
		fatalError(t, nil, "Unable to evaluate integer rule")
	}

	cond2 := NewCondition("a == `peterpiper`")
	args2 := map[string]any{"a": "peterpiper"}
	result2, err2 := cond2.Rule(args2)
	if !result2 || err2 != nil {
		fatalError(t, nil, "Unable to evaluate string rule")
	}

	cond3 := NewCondition("a == `peterpiper` && b == `paid the boss`")
	args3 := map[string]any{"a": "peterpiper", "b": "paid the boss"}
	result3, err3 := cond3.Rule(args3)
	if !result3 || err3 != nil {
		fatalError(t, nil, "Unable to evaluate string rule")
	}
}

func TestSelect(t *testing.T) {

	state := new(State)
	cfg := NewSelectCfg()
	for i := 0; i < 3; i++ {
		cfg.AddEdge(new(Edge))
	}

	cfg.AddCond("A == `String A`")
	cfg.AddCond("A == `String B`")
	cfg.AddCond("B == `String B`")
	state.selectcfg = cfg
	payload := new(Payload)
	arg := NewArg()
	arg.AddArg("A", "String A")
	arg.AddArg("B", "String B")
	payload.SetData("argument", arg)
	err := Select(state, payload)
	if err != nil {
		fatalError(t, nil, "Unable to switch")
	}
	if cfg.edges[0].Selected != true {
		fatalError(t, nil, "Edge not selected")
	}
	if cfg.edges[1].Selected != false {
		fatalError(t, nil, "Edge selected")
	}
	if cfg.edges[2].Selected != true {
		fatalError(t, nil, "Edge not selected")
	}
}
