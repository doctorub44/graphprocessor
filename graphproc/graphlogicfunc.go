package graphproc

import (
	"github.com/maja42/goval"
)

//  "graph1: a | b; b | d | k; b | c | f | h | k; b | e | g; g | i | l; g | j | l; graph2:k|l|m{\"k1\":\"v1\", \"k2\":\"v2\"}"
//
// a > b > c
//		 > d
//		 > e

type Condition struct {
	expr string
	eval *goval.Evaluator
}

type Argument struct {
	args map[string]any
}

type SwitchCfg struct {
	edges []*Edge
	conds []*Condition
}

func Switch(st *State, payload *Payload) error {
	config := st.config.(*SwitchCfg)
	key := "argument"
	a, err := payload.GetData(key)
	if err != nil {
		return err
	}

	for i, c := range config.conds {
		if result, _ := c.Rule(a.(Argument).args); result {
			config.edges[i].Selected = true
		} else {
			config.edges[i].Selected = false
		}
	}
	return nil
}

func NewCondition(expr string) *Condition {
	cond := new(Condition)
	cond.expr = expr
	cond.eval = goval.NewEvaluator()
	return cond
}

func (c *Condition) Rule(args map[string]any) (bool, error) {
	result, err := c.eval.Evaluate(c.expr, args, nil)
	return result.(bool), err
}
