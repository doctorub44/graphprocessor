package graphproc

import (
	"encoding/json"
	"fmt"
	"io"
)

// Graph description language
// name1:a|b;b|d|k;b|c|f|h|k;b|e|g;g|i|l;g|j|l;k|l|m{“k1”:”v1”, “k2”:”v2”};name2:...

//Parser : represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

//NewParser : returns a new instance of Parser
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

//Parse : parse graph definition and output the slice of graphs
func (p *Parser) Parse() ([]string, []*Graph, error) {
	var v1, v2, lastvert *Vertex
	var ok, brace bool
	var gnames []string
	var graphs []*Graph
	var g *Graph
	var vertexmap map[string]*Vertex

	depthcount := 0

	// Loop over the graph definition
	for {
		// Read a token
		tok, lit := p.scanIgnoreWhitespace()
		if tok == EOF {
			break
		}
		if tok == IDENT {
			tok2, _ := p.scanIgnoreWhitespace()
			//Check for new graph specification
			if tok2 == COLON {
				g = NewGraph()
				graphs = append(graphs, g)
				gnames = append(gnames, lit)
				vertexmap = make(map[string]*Vertex)
			} else {
				p.unscan()
				if v1, ok = vertexmap[lit]; !ok {
					v1 = g.NewVertex(lit)
					vertexmap[lit] = v1
				}
				if v2 != nil {
					g.Link(v2, v1)
					v2 = v1
				}
				lastvert = v1
			}
		} else if tok == PIPE { //If PIPE, expect a vertex name in next token
			tok, lit2 := p.scanIgnoreWhitespace()
			if tok != IDENT {
				return nil, nil, fmt.Errorf("found %q, expected vertex", lit2)
			} else {
				if v2, ok = vertexmap[lit2]; !ok {
					v2 = g.NewVertex(lit2)
					vertexmap[lit2] = v2
				}
				g.Link(v1, v2)
				v1 = v2
				lastvert = v2
			}
		} else if tok == SEMICOLON { //If SEMICOLON, reset the second vertex since we are starting a new path
			v2 = nil
			continue
		} else if tok == LEFTBRACE || tok == LEFTBRACKET { //If LEFTBRACE/LEFTBRACKET, Scan to the RIGHTBRACE/RIGHTBRACKET and convert to JSON object
			if tok == LEFTBRACE {
				brace = true
			}
			jsonlit := lit
			depthcount++

			for i := 0; i < 2048; i++ {
				tok, lit = p.scanIgnoreWhitespace()
				if tok != IDENT && tok != DOUBLEQUOTE && tok != COMMA && tok != COLON && tok != ESCAPE &&
					tok != LEFTBRACE && tok != RIGHTBRACE && tok != RIGHTBRACKET && tok != LEFTBRACKET &&
					tok != PERIOD && tok != UNDERSCORE && tok != DASH && tok != SPACE {
					return nil, nil, fmt.Errorf("found %q, expected valid json string", lit)
				}
				jsonlit = jsonlit + lit
				if depthcount == 1 {
					if (tok == RIGHTBRACE && brace) || (tok == RIGHTBRACKET && !brace) { //If the next token is RIGHTBRACE or RIGHTBRACKET and it is the outermost one, assign the configuration state
						var config interface{}
						err := json.Unmarshal([]byte(jsonlit), &config)
						if err != nil {
							return nil, nil, fmt.Errorf("found %q, error unmarshalling json string", jsonlit)
						}
						lastvert.Vstage.BuildState(config)
						depthcount--
						break
					}
				}

				if tok == RIGHTBRACE && brace {
					depthcount--
				} else if tok == LEFTBRACE && brace {
					depthcount++
				} else if tok == RIGHTBRACKET && !brace {
					depthcount--
				} else if tok == LEFTBRACKET && !brace {
					depthcount++
				}

				if i == 4096 {
					return nil, nil, fmt.Errorf("expecting right brace or bracket prior to maximum of 4096 tokens in json string: %q", jsonlit)
				}
			}
			continue
		} else {
			return nil, nil, fmt.Errorf("found %q, invalid token", lit)
		}
	}

	if depthcount != 0 {
		return nil, nil, fmt.Errorf("unable to scan json string - unmatched out brace or bracket")
	}

	return gnames, graphs, nil
}

// scan : returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

//unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.buf.n = 1
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}
