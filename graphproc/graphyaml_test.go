package graphproc

import (
	"fmt"
	"testing"
)

func TestGraphYaml(t *testing.T) {
	graphj := GraphYaml("samplegraph.yaml")
	g := YamlToGraph(graphj)

	gy := GraphToYaml(g)
	fmt.Println(gy.Marshal())
}
