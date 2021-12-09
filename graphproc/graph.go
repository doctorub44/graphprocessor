package graphproc

//Vertex : graph vertex type
type Vertex struct {
	Name     string
	Prev     []*Edge
	Next     []*Edge
	Vstage   *Stage
	joined   int
	forked   int
	produced int
}

//Edge : graph edge type
type Edge struct {
	//In     *Vertex
	Out      *Vertex
	Epayload *Payload
	Estate   *State
}

//Graph : graph type of vertexes and edges and execution path
type Graph struct {
	V       []*Vertex
	E       []*Edge
	Path    []*Vertex
	current int

	step    *Vertex
	forking []*Vertex
}

//NewGraph : create a new graph
func NewGraph() *Graph {
	g := new(Graph)
	g.V = make([]*Vertex, 0, 32)
	g.E = make([]*Edge, 0, 32)
	g.Path = make([]*Vertex, 0, 32)
	g.forking = make([]*Vertex, 0, 32)
	return g
}

//NewVertex : create a new vertex
func (g *Graph) NewVertex(n string) *Vertex {
	v := new(Vertex)
	v.Name = n
	v.Prev = make([]*Edge, 0, 32)
	v.Next = make([]*Edge, 0, 32)
	v.Vstage = NewStage()
	g.V = append(g.V, v)
	return v
}

//NewEdge : create a new edge
func (g *Graph) NewEdge() *Edge {
	e := new(Edge)
	e.Epayload = new(Payload)
	e.Epayload.Raw = make([]byte, 0, 2048)
	g.E = append(g.E, e)
	return e
}

//NextVertex : get the next vertex in the path, returns nil if at the path end
func (g *Graph) NextVertex() *Vertex {
	if g.current == len(g.Path) {
		return nil
	}
	v := g.Path[g.current]
	g.current++
	return v
}

//Link : link two vertexes by creating an edge from vertex 1 to vertex 2
func (g *Graph) Link(v1 *Vertex, v2 *Vertex) {
	e := g.NewEdge()
	e.Out = v2
	v1.Next = append(v1.Next, e)
	v2.Prev = append(v2.Prev, e)
}

//Path : build the execution path by walking through the graph depth first until a vertex is found with incomplete inputs
//	go back to the last vertex with a fork in the path and start the next walk
func (g *Graph) BuildPath() {
	if g.step == nil { //First step in building the path
		g.step = g.V[0]
		g.Path = append(g.Path, g.step)
	}

	for {
		if len(g.step.Prev) == g.step.joined { //All inputs are present
			if len(g.step.Next) == 0 { // We are done
				break
			}
			if len(g.step.Next) == 1 { // Simple serial step in the path
				g.step = g.step.Next[0].Out
				g.step.joined++
				g.Path = append(g.Path, g.step)
			} else { //This vertex is a fork in the path
				if g.step.forked < len(g.step.Next) {
					if g.step.forked == 0 { //If first time, add it to the slice of vertexes forking
						g.forking = append(g.forking, g.step)
					}
					g.step.forked++
					g.step = g.step.Next[g.step.forked-1].Out
					g.step.joined++
					g.Path = append(g.Path, g.step)
				} else { //Finished with this vertex's forking, remove it from the forking slice
					g.forking = g.forking[1:]
					g.step = g.forking[0]
				}
			}
		} else if len(g.forking) > 0 { //Go back to a vertex that is forking since this vertex is waiting for inputs
			g.step = g.forking[0]
		} else {
			continue
		}
	}
}
