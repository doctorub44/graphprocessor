package graphproc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	deepcopy "github.com/barkimedes/go-deepcopy"
)

// Stage : graph execution stage type
type Stage struct {
	name    string
	sfunc   func(*State, *Payload) error
	state   *State
	preasp  map[string]*Aspect
	postasp map[string]*Aspect
}

// State : graph stage state type
type State struct {
	config    any
	appstate  any
	selectcfg *SelectCfg
}

// Payload : graph stage payload type
type Payload struct {
	Raw  []byte
	Data interface{}
}

// Graphline : map with the set of graphs and map of application stages used in the graphs
type Graphline struct {
	graphs   map[string]*Graph
	appstage map[string]func(*State, *Payload) error
}

var publicgraphs *Graphline

// NewGraphline : create a new graphline
func NewGraphline() Graphline {
	var g Graphline
	g.graphs = make(map[string]*Graph)
	g.appstage = make(map[string]func(*State, *Payload) error)
	return g
}

// PublishGraphs : publish the graphline if needed elsewhere (subgraph execution stage)
func (g *Graphline) PublishGraphs() {
	publicgraphs = g
}

// PublicGraphs : return the graphline that was published
func PublicGraphs() *Graphline {
	return publicgraphs
}

// NewState : create a state type
func NewState() *State {
	s := new(State)
	s.config = make(map[string]interface{})
	return s
}

// NewPayload : create a new payload type
func NewPayload() *Payload {
	p := new(Payload)
	p.Raw = make([]byte, 0, 2048)
	return p
}

// NewStage : create and new graph execution stage type
func NewStage() *Stage {
	s := new(Stage)
	s.state = new(State)
	s.preasp = make(map[string]*Aspect)
	s.postasp = make(map[string]*Aspect)
	return s
}

// Config : get a field value when using the convention of map of interface values
func (s *State) Config(field string) (string, bool) {
	if f, ok := s.config.(map[string]interface{})[field].(string); ok {
		return f, true
	}
	return "", false
}

// SetConfig : set a field value when using the convention of map of interface values
func (s *State) SetConfig(field string, value interface{}) {
	s.config.(map[string]interface{})[field] = value
}

// funcName : returns a string for the name of a func in form <package>.<func name>
func funcName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

//callerName : return name of the calling function
//func callerName(index int) string {
//	pc, _, _, _ := runtime.Caller(index)
//	details := runtime.FuncForPC(pc)
//	return details.Name()
//}

// Copy : deep copy the src payload to the destination payload
func (p *Payload) Copy(dest, src *Payload) (*Payload, error) {
	var err error
	err = nil
	dest.Raw = append(dest.Raw[:0], src.Raw...)
	if src.Data != nil {
		dest.Data, err = deepcopy.Anything(src.Data)
	}
	return dest, err
}

// Append : deep append the src payload to the destination payload. Data field must be map or slice
func (p *Payload) Append(dest, src *Payload) (*Payload, error) {
	var err error
	dest.Raw = append(dest.Raw, src.Raw...)
	if src.Data != nil {
		if reflect.TypeOf(src.Data).Kind() == reflect.Slice {
			dest.Data = reflect.AppendSlice(reflect.ValueOf(dest.Data), reflect.ValueOf(src.Data))
		} else if reflect.TypeOf(src.Data).Kind() == reflect.Map {
			dest.Data, err = appendMap(dest.Data, src.Data)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("Append: source type invalid - must be slice or map")
		}
	}
	return dest, nil
}

// SetData : set the data field when using convention of map of interface values
func (p *Payload) SetData(key string, value interface{}) error {
	if p.Data == nil {
		p.Data = make(map[string]interface{})
	}
	p.Data.(map[string]interface{})[key] = value
	return nil
}

// GetData : get the data field when using convention of map of interface values
func (p *Payload) GetData(key string) (interface{}, error) {
	if value, ok := p.Data.(map[string]interface{})[key]; ok {
		return value, nil
	}
	return nil, errors.New("GetData: key not found in payload data")
}

// appendMap : does a deep append of the source map to destination map
func appendMap(dest, src interface{}) (interface{}, error) {
	newmap := make(map[string]interface{})

	iter := reflect.ValueOf(src).MapRange()
	for iter.Next() {
		newmap[iter.Key().String()] = iter.Value()
	}

	iter = reflect.ValueOf(dest).MapRange()
	for iter.Next() {
		newmap[iter.Key().String()] = iter.Value()
	}

	return newmap, nil
}

// RegisterStage : register an application stage function for use in a graph
func (g *Graphline) RegisterStage(sfunc func(*State, *Payload) error) error {
	fields := strings.Split(funcName(sfunc), ".")
	aname := fields[len(fields)-1]
	if _, aok := g.appstage[aname]; !aok {
		g.appstage[aname] = sfunc
	}
	return nil
}

// RegisterAspect : register an application aspect function for use in a graph
func (g *Graphline) RegisterAspect(graphid string, action AspectAction, stage func(*State, *Payload) error, newasp *Aspect) error {
	graph := g.graphs[graphid]
	for _, vertex := range graph.Path {
		s := vertex.Vstage
		if stage == nil || funcName(s.sfunc) == funcName(stage) {
			aspname := funcName(newasp.aspect)
			if action == START {
				if a, ok := s.preasp[s.name]; !ok {
					s.preasp[aspname] = newasp
				} else if funcName(a.aspect) != aspname {
					s.preasp[aspname] = newasp
				}
			} else if action == END {
				if a, ok := s.postasp[s.name]; !ok {
					s.postasp[aspname] = newasp
				} else if funcName(a.aspect) != aspname {
					s.postasp[aspname] = newasp
				}
			} else {
				return errors.New("RegisterAspect: illegal action passed, must be START or END : " + strconv.Itoa(int(action)))
			}

			if stage != nil { //if nil stage not specified, return, otherwise apply aspect to all stages
				return nil
			}
		}
	}

	return nil
}

// Sequence : parse and generate the execution sequence for a graph specification
func (g *Graphline) Sequence(graphspec string) ([]string, error) {
	if strings.TrimSpace(graphspec) == "" {
		return nil, errors.New("Sequence: no graph specification provided ")
	}
	gnames, graphs, err := (NewParser(strings.NewReader(graphspec))).Parse()
	if err == nil {
		for i, graph := range graphs {
			graph.BuildPath()
			g.graphs[gnames[i]] = graph
			for _, vertex := range graph.Path {
				if sfunc, aok := g.appstage[vertex.Name]; aok {
					vertex.Vstage.sfunc = sfunc
					vertex.Vstage.name = vertex.Name
				} else {
					return nil, errors.New("Sequence: no registered stage found for " + vertex.Name)
				}
			}
		}
	} else {
		return nil, err
	}
	return gnames, nil
}

// PrintPath : print the path for diagnostic purposes
func (g *Graphline) PrintPath(graphid string) {
	graph := g.graphs[graphid]
	for _, v := range graph.Path {
		fmt.Println(v.Name)
	}
}

// Execute : execute a graph based on its name
func (g *Graphline) Execute(gname string, payload *Payload) error {
	var err, serr error
	graph := g.graphs[gname]
	scope := new(Scope)

	for _, vertex := range graph.Path {
		stage := vertex.Vstage
		name := funcName(stage.sfunc)

		if stage.sfunc != nil {
			//Execute pre aspects for this stage
			if err = g.RunAspects(name, payload.Raw, scope, stage.preasp); err != nil {
				return err
			}

			if vertex.joined == 0 { //First vertex in graph so no input data, execute the vertex
			} else if vertex.joined == 1 { // Only a single edge with input data, copy the input and execute the vertex
				payload, err = payload.Copy(payload, vertex.Prev[0].Epayload)
				if err != nil {
					return err
				}
			} else if vertex.joined == vertex.produced { //All edges have produced input data, aggregate the inputs
				payload.Raw = payload.Raw[:0]
				for _, e := range vertex.Prev {
					payload, err = payload.Append(payload, e.Epayload)
					if err != nil {
						return err
					}
				}
			} else if vertex.joined != vertex.produced { //Not all edges have input data so the vertex is not ready to execute, continue to next vertex in path
				continue
			}
			//Execute the stage function
			if serr = stage.sfunc(stage.state, payload); serr == nil {
				for _, edge := range vertex.Next { //Copy the output to each output edge selected, all edges by default
					if edge.Selected {
						edge.Epayload, err = payload.Copy(edge.Epayload, payload)
						if err != nil {
							return err
						}
					}
					edge.Out.produced++
				}
			} else {
				return serr
			}

			//Execute post aspects for this stage
			if err = g.RunAspects(name, payload.Raw, scope, stage.postasp); err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

// SetState : deprecated
func (g *Graphline) SetState(gname string, name string, state *State) (int, error) {
	graph := g.graphs[gname]
	for i, v := range graph.V {
		if name == v.Name {
			v.Vstage.state = state
			graph.V[i] = v
			return 0, nil
		}
	}

	return -1, errors.New("SetStage: no registered stage found for " + name)
}

// CallStage : deprecated
func (g *Graphline) CallStage(graphid string, name string, state *State) (int, error) {
	var err error
	graph := g.graphs[graphid]
	for i, v := range graph.V {
		if name == v.Name {
			err = v.Vstage.sfunc(state, nil)
			graph.V[i] = v
			return i, err
		}
	}
	return -1, errors.New("CallStage: no registered stage found for " + name)
}

// RunAspects : run the aspects for a graph stage
func (g *Graphline) RunAspects(name string, mesg []byte, scope *Scope, asp map[string]*Aspect) error {
	for _, aspect := range asp {
		if aspect != nil {
			if action := aspect.Execute(name, mesg, scope); action == ERROR {
				return errors.New("RunAspects: error returned by execution aspect")
			}
		} else {
			break
		}
	}
	return nil
}

// BuildState
func (s *Stage) BuildState(cfg interface{}) {
	if s.state == nil {
		s.state = new(State)
	}
	s.state.config = cfg
}

// GetConfig
func (s *Stage) GetConfig() any {
	if s.state != nil {
		return s.state.config
	}
	return nil
}

// GetSelect
func (s *Stage) GetSelect() *SelectCfg {
	if s.state != nil {
		return s.state.selectcfg
	}
	return nil
}

// SetSelect
func (s *Stage) SetSelect(selcfg *SelectCfg) {
	if s.state != nil {
		s.state.selectcfg = selcfg
	}
}

// SelectEdge
func (s *Stage) SelectEdge(e *Edge) {
	if s.state != nil {
		if s.state.selectcfg != nil {
			s.state.selectcfg.AddEdge(e)
		}
	}
}
