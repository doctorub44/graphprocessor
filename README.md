# graphprocessor
Graphprocessor allows applications to be assembled dynamically from a set of Go functions without recompilation. An analogy is Linux shell commands where output is piped between commands for a desired result. An directed graph is specified with a textual definition allowing different graphs to be constructed and executed without additional code. New capabilities are added by adding a new graph function to the available set. Aspect functions can be assigned to graph nodes for error handling, logging, metrics, and observability, removing common logic from individual graph functions. 

# Graph Definition 
Create a simple pipeline, **graph1**.

`graph1:Func1|Func2|Func3`

Create multiple graphs, **graph1** and **graph2**.

`graph1:Func1|Func2|Func3;graph2:Func3|Func4`

Create a graph with JSON configuration data for a graph function.

`graph1:Func1|Func2{"keyname":"value"}`

Create a graph with a subgraph.

`graph1:Func1|Func2|Func3`
`graph2:FuncA|SubGraph{"name":"graph1"}|FuncB`

**graph2** has an execution path of:

FuncA > Func1 > Func2 > Func3 > FuncB

Create an directed graph with branches.

`graph1:Func1|Func2|Func3;Func2|Func4;Func2|Func5;Func2|Func6;Func3|Func7;Func4|Func7;Func5|Func7;Func6|Func7`
![Example Graph](sample_graph.jpg)

# Graph Functions
Each graph function has a calling convention with two parameters: state and payload.  Payload is data passed between graph function during processing. That state parameter is local to the graph function and contains configuration data from when the graph was defined and any run time state the function requires.
```
func Func1(s *State, payload *Payload) error {
	//Code...
	return nil
}
```
## Payload
The **Payload** type provides two options for passing data between graph functions, a raw byte stream or an application defined Data type.  Application functions need to agree on which field to use. For example, the initial graph function unmarshalls JSON in the **Raw** field into a structure in the **Data** field. Subsequent graph functions in the processing stream operate on the **Data** field with the final graph function marshalling the JSON Data back to **Raw**.
```
type Payload struct {
	Raw  []byte
	Data interface{}
}
```
One current limitation on the **Data** field is if multiple nodes in the graph converge to a single node the action is to aggregate the inputs using a deep copy and append. In this scenario, the **Data** field is limited to slice or map data types.
## State
The **State** type associated with each graph function has two fields: **config** and **appstate**. The **config** field contains unmarshalled JSON configuration data from the graph definition such as **Func2{"keyname":"value"}**. The **appstate** field is for a graph function to persist state between repeated calls of the function. A typical example is a data stream with the graph function writing a data stream to a file: the open file information is stored in **appstate**.
```
type State struct {
	config   interface{}
	appstate interface{}
}
```
# Aspects
Aspects are cross cutting functions assigned as a pre-aspect prior to execution of a graph function or a post-aspect after the execution of a graph function. Aspects can be assigned to individual graph functions or all of them. Aspects are used to metrics collection, diagnostics, and error handling.  Individual graph functions are not aware of or require any code to apply aspects.
```
func Aspect1(name string, mesg []byte, scope *Scope) error {
	//Code ...
	return nil
}
```
# Sample Graph Processor
This example graph processor is for processing Indicators of Compromise (IOC).  The graph inputs source data then forks three nodes with each recieving a copy of the input and search for URL IOCs, IPv4 address IOCs, and file MD5 hash IOCs with the results aggregated and converted into a data structure that is marshalled into JSON.  A trace aspect is assigned to each node in the graph and a start timer aspect to the first node, an end timer aspect is assigned to the last node along with a duration aspect that calculates the duration.

A simple logging mechanism is initialized by **EventInit()** with a goroutine receiving logs over a buffered channel from the **Trace()** aspect .
```
EventInit()
```
A new graph is graph is created.
```
graph := NewGraphline()
```
Each Go graph function used in a graph is registered with **RegisterStage()**.
```
err := graph.RegisterStage(InputIOC)
```
The graph is created with **Sequence()** that returns a list graph names and computes the execution path of the graph.
```
graphnames, err := graph.Sequence(`graph1:InputIOC|UrlIOC|IOCToData;InputIOC|Ipv4IOC|IOCToData;InputIOC|Md5IOC|IOCToData;IOCToData|IOCDataToJson`
```
The aspects are created and assigned and then the graph executed with **Execute()**.  Registering an aspect against **nil** applies it to all nodes in the graph.
```
err = graph.Execute(graphid[0], NewPayload()
```
A sleep at the end allows any remaining logging data to be written to the log file.

Full example code below.

```
func main() {
	EventInit()
	graph := NewGraphline()

	err := graph.RegisterStage(InputIOC)
	err  = graph.RegisterStage(UrlIOC)
	err  = graph.RegisterStage(Ipv4IOC)
	err  = graph.RegisterStage(Md5IOC)
	err  = graph.RegisterStage(IOCToData)
	err  = graph.RegisterStage(IOCDataToJson)

	graphnames, err := graph.Sequence(`graph1:InputIOC|UrlIOC|IOCToData;InputIOC|Ipv4IOC|IOCToData;InputIOC|Md5IOC|IOCToData;IOCToData|IOCDataToJson`)

	start    := NewAspect(StartTime, 	time.Duration(10*time.Millisecond), 0, CONTINUE)
	end      := NewAspect(EndTime, 		time.Duration(10*time.Millisecond), 0, CONTINUE)
	duration := NewAspect(DurationTime,	time.Duration(10*time.Millisecond), 0, CONTINUE)
	trace    := NewAspect(Trace, 		time.Duration(10*time.Millisecond), 0, CONTINUE)

	graph.RegisterAspect(graphnames[0], START,	nil, 		&trace)
	graph.RegisterAspect(graphnames[0], START,	InputIOC, 	&start)
	graph.RegisterAspect(graphnames[0], END, 	IOCDataToJson, 	&end)
	graph.RegisterAspect(graphnames[0], END, 	IOCDataToJson, 	&duration)

	err      = graph.Execute(graphnames[0], NewPayload())

	time.Sleep(2000 * time.Millisecond)
}
```
