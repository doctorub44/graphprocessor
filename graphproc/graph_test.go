package graphproc

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseAndPath(t *testing.T) {
	_, graphs, err := (NewParser(strings.NewReader("graph1:a | b; b|d|k;b|c|f|h|k;b|e|g ;g|i|l; g|j|l;graph2:k|l|m{\"k1\":\"v1\", \"k2\":\"v2\"}"))).Parse()
	if err == nil {
		for _, g := range graphs {
			g.BuildPath()
			fmt.Println(len(g.Path))
		}
	} else {
		fatalError(t, err, "Unable to build path")
	}
}
