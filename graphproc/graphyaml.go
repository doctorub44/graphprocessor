package graphproc

import (
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

func YamlToGraph(gyaml *GraphY) *Graph {
	var v1 *Vertex
	var ok bool
	g := NewGraph()
	g.SetName(gyaml.Name)
	vertexmap := make(map[string]*Vertex)
	for _, yamlv := range gyaml.Vertex {
		if v1, ok = vertexmap[yamlv.Name]; !ok {
			v1 = g.NewVertex(yamlv.Name)
			vertexmap[yamlv.Name] = v1
		}
		for _, tov := range yamlv.To {
			ToNext(g, vertexmap, v1, &tov)
		}
	}

	return g
}

func ToNext(g *Graph, vertexmap map[string]*Vertex, v1 *Vertex, yamlv *VertexY) {
	var v2 *Vertex
	var ok bool
	if v2, ok = vertexmap[yamlv.Name]; !ok {
		v2 = g.NewVertex(yamlv.Name)
		vertexmap[yamlv.Name] = v2
	}
	g.Link(v1, v2)
	for _, tov := range yamlv.To {
		ToNext(g, vertexmap, v2, &tov)
	}
}

func GraphToYaml(g *Graph) *GraphY {
	var vy *VertexY
	var ok bool
	vertexmap := make(map[string]*VertexY)
	gy := new(GraphY)
	gy.Name = g.Name
	for _, v := range g.V {
		if _, ok = vertexmap[v.Name]; !ok {
			vy = new(VertexY)
			vy.Name = v.Name
			vertexmap[v.Name] = vy
			for _, edge := range v.Next {
				ToNextY(vertexmap, vy, edge.Out)
			}
			gy.Vertex = append(gy.Vertex, *vy)
		}
	}
	return gy
}

func ToNextY(vertexmap map[string]*VertexY, vy *VertexY, v *Vertex) {
	var newvy *VertexY
	var ok bool
	if newvy, ok = vertexmap[v.Name]; !ok {
		newvy = new(VertexY)
		newvy.Name = v.Name
		vertexmap[v.Name] = newvy
	}
	for _, edge := range v.Next {
		ToNextY(vertexmap, newvy, edge.Out)
	}
	vy.To = append(vy.To, *newvy)
}

func (gy *GraphY) Marshal() string {
	yamlg, _ := yaml.Marshal(gy)
	return string(yamlg)
}
