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

type SelectCfg struct {
	edges []*Edge
	conds []*Condition
}

func NewCondition(expr string) *Condition {
	cond := new(Condition)
	cond.expr = expr
	cond.eval = goval.NewEvaluator()
	return cond
}

func NewArg() *Argument {
	arg := new(Argument)
	arg.args = make(map[string]any)
	return arg
}

func NewSelectCfg() *SelectCfg {
	cfg := new(SelectCfg)
	cfg.edges = make([]*Edge, 0, 8)
	cfg.conds = make([]*Condition, 0, 8)
	return cfg
}

func (a *Argument) AddArg(key, val string) {
	a.args[key] = val
}

func (c *SelectCfg) AddCond(cond string) {
	c.conds = append(c.conds, NewCondition(cond))
}
func (c *SelectCfg) AddEdge(e *Edge) {
	c.edges = append(c.edges, e)
}

func Select(st *State, payload *Payload) error {
	config := st.config.(*SelectCfg)
	a, err := payload.GetData("argument")
	if err != nil {
		return err
	}

	for i, c := range config.conds {
		if result, _ := c.Rule(a.(*Argument).args); result {
			config.edges[i].Selected = true
		} else {
			config.edges[i].Selected = false
		}
	}
	return nil
}

func (c *Condition) Rule(args map[string]any) (bool, error) {
	result, err := c.eval.Evaluate(c.expr, args, nil)
	return result.(bool), err
}
