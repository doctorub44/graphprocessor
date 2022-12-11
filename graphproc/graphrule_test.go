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
