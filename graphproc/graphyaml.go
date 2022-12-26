package graphproc

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

type VertexY struct {
	Name      string    `json:"Name"`
	Condition []string  `json:"Condition"`
	Param     string    `json:"Param"`
	To        []VertexY `json:"To"`
}

type GraphY struct {
	Name   string    `json:"Name"`
	Vertex []VertexY `json:"Vertex"`
}

func NewGraphY(name string) *GraphY {
	gy := new(GraphY)
	gy.Name = name
	return gy
}

func NewVertexY(name string) *VertexY {
	vy := new(VertexY)
	vy.Name = name
	return vy
}

func GraphYaml(file string) *GraphY {
	gyaml := new(GraphY)
	data, err := os.ReadFile(file)
	if err == nil {
		err := yaml.Unmarshal(data, gyaml)
		if err == nil {
		} else {
			fmt.Println("error converting YAML")
		}
	} else {
		fmt.Println("error reading graph yaml")
	}
	return gyaml
}

func YamlToGraph(gy *GraphY) *Graph {
	g := NewGraph()
	g.SetName(gy.Name)
	vertexmap := make(map[string]*Vertex)
	for _, vy := range gy.Vertex {
		v1, _ := MakeVertex(g, vertexmap, &vy)
		for _, tov := range vy.To {
			ToNext(g, vertexmap, v1, &tov)
		}
	}

	return g
}

func ToNext(g *Graph, vertexmap map[string]*Vertex, v1 *Vertex, vy *VertexY) {
	v2, _ := MakeVertex(g, vertexmap, vy)
	g.Link(v1, v2)

	for _, tov := range vy.To {
		ToNext(g, vertexmap, v2, &tov)
	}
}
func MakeVertex(g *Graph, vertexmap map[string]*Vertex, vy *VertexY) (*Vertex, bool) {
	var v2 *Vertex
	var ok bool
	if v2, ok = vertexmap[vy.Name]; !ok {
		v2 = g.NewVertex(vy.Name)
		if vy.Param != "" {
			var config any
			json.Unmarshal([]byte(vy.Param), &config)
			v2.SetConfig(config)
		}
		if vy.Condition != nil {
			selcfg := NewSelectCfg()
			for _, c := range vy.Condition {
				selcfg.AddCond(c)
			}
			v2.SetSelect(selcfg)
		}
		vertexmap[vy.Name] = v2
	}
	return v2, ok
}

func GraphToYaml(g *Graph) *GraphY {
	vertexmap := make(map[string]*VertexY)
	gy := NewGraphY(g.Name)

	for _, v := range g.V {
		if vy, ok := MakeVertexY(vertexmap, v); !ok {
			for _, edge := range v.Next {
				ToNextY(vertexmap, vy, edge.Out)
			}
			gy.Vertex = append(gy.Vertex, *vy)
		}
	}
	return gy
}

func ToNextY(vertexmap map[string]*VertexY, vy *VertexY, v *Vertex) {
	newvy, _ := MakeVertexY(vertexmap, v)
	for _, edge := range v.Next {
		ToNextY(vertexmap, newvy, edge.Out)
	}
	vy.To = append(vy.To, *newvy)
}

func MakeVertexY(vertexmap map[string]*VertexY, v *Vertex) (*VertexY, bool) {
	var newvy *VertexY
	var ok bool
	if newvy, ok = vertexmap[v.Name]; !ok {
		newvy = NewVertexY(v.Name)
		if cfg := v.GetConfig(); cfg != nil {
			config, _ := json.Marshal(cfg)
			newvy.Param = string(config)
		}
		if selcfg := v.GetSelect(); selcfg != nil {
			for _, c := range selcfg.conds {
				newvy.Condition = append(newvy.Condition, c.expr)
			}
		}
		vertexmap[v.Name] = newvy
	}
	return newvy, ok
}

func (gy *GraphY) Marshal() string {
	yamlg, _ := yaml.Marshal(gy)
	return string(yamlg)
}
