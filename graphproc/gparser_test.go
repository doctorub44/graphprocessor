package graphproc

import (
	"fmt"
	"strings"
	"testing"
)

func TestGraphParse(t *testing.T) {
	gnames, graphs, err := (NewParser(strings.NewReader(`graph1:a | b; b|d|k;b|c|f|h|k;b|e|g ;g|i|l; g|j|l;graph2:k|l|m`))).Parse()
	fmt.Println(graphs)
	fmt.Println(gnames)
	fatalError(t, err, "Unable to parse json")
}
func TestJsonParse1(t *testing.T) {
	gnames, graphs, err := (NewParser(strings.NewReader(`graph2:k|l|m{"k1":"value"}`))).Parse()
	fmt.Println(graphs)
	fmt.Println(gnames)
	fatalError(t, err, "Unable to parse json")
}
func TestJsonParse2(t *testing.T) {
	gnames, graphs, err := (NewParser(strings.NewReader(`graph2:k|l|m{"k1":{"v1":"vv1"}}`))).Parse()
	fmt.Println(graphs)
	fmt.Println(gnames)
	fatalError(t, err, "Unable to parse json")
}

func TestJsonParse3(t *testing.T) {
	gnames, graphs, err := (NewParser(strings.NewReader(`graph2:k|l|m[{"k1":{"v1":"vv1"}}, {"k2":{"v2":"vv2"}}]`))).Parse()
	fmt.Println(graphs)
	fmt.Println(gnames)
	fatalError(t, err, "Unable to parse json")
}
