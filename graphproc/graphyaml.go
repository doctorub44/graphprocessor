package graphproc

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

type VertexJ struct {
	Name      string    `json:"Name"`
	Condition []string  `json:"Condition"`
	Param     string    `json:"Param"`
	To        []VertexJ `json:"To"`
}

type GraphJ struct {
	Name   string    `json:"Name"`
	Vertex []VertexJ `json:"Vertex"`
}

func GraphYaml(file string) {
	gjson := new(GraphJ)
	data, err := os.ReadFile(file)
	if err == nil {
		err := yaml.Unmarshal(data, gjson)
		if err == nil {
		} else {
			fmt.Println("error converting YAML")
		}
	} else {
		fmt.Println("error reading graph yaml")
	}
}
