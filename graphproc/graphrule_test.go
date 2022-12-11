package graphproc

import (
	"testing"
)

func TestRule(t *testing.T) {
	cond := NewCondition("a == 5")
	args := map[string]any{"a": 5}
	result, err := cond.Rule(args)
	if !result || err != nil {
		fatalError(t, nil, "Unable to evaluate rule")
	}
}
