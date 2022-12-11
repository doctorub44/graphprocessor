package graphproc

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func fatalError(t *testing.T, err error, s string) {
	if err != nil {
		t.Errorf(s)
	}
}

func normalIOC(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".normalioc")...)
	return nil
}

func filterIOC(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".filterioc")...)
	return nil
}

func urlIOC(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".urlioc")...)
	return nil
}

func ipIOC(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".ipioc")...)
	return nil
}

func hashIOC(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".hashioc")...)
	return nil
}

func jsonout(s *State, payload *Payload) error {
	text := string(payload.Raw)
	text = "{" + text + "}"
	payload.Raw = append(payload.Raw[:0], []byte(text)...)
	return nil
}

func func1(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".func1")...)
	return nil
}

func func2(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".func2")...)
	return nil
}

func func3(s *State, payload *Payload) error {
	var event Gevent

	payload.Raw = append(payload.Raw, []byte(".func3")...)
	event.id = MESSAGE
	event.strval = "func3: " + string(payload.Raw)
	SendEvent(event)

	time.Sleep(100 * time.Millisecond)

	return nil
}

func func4(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".func4")...)
	return nil
}

func func5(s *State, payload *Payload) error {
	payload.Raw = append(payload.Raw, []byte(".func5")...)
	return nil
}

func TestRegister(t *testing.T) {
	graph := NewGraphline()
	fatalError(t, graph.RegisterStage(func1), "Unable to register func 1")
	fatalError(t, graph.RegisterStage(func2), "Unable to register func 2")
	fatalError(t, graph.RegisterStage(func3), "Unable to register func 3")
}

func TestSequence0(t *testing.T) {
	graph := NewGraphline()

	fatalError(t, graph.RegisterStage(func1), "Unable to register func 1")
	fatalError(t, graph.RegisterStage(func2), "Unable to register func 2")
	fatalError(t, graph.RegisterStage(func3), "Unable to register func 3")

	_, err := graph.Sequence("")
	if err == nil {
		fatalError(t, errors.New("Should return error"), "Unable to create sequence 0")
	}
	_, err = graph.Sequence("graph:f")
	if err == nil {
		fatalError(t, errors.New("Should return error"), "Unable to create sequence 1")
	}
	_, err = graph.Sequence("graph:func1")
	fatalError(t, err, "Unable to create sequence 2")
	_, err = graph.Sequence("graph:func1|func2|func3")
	fatalError(t, err, "Unable to create sequence 3")
}
func TestPayloadCopy(t *testing.T) {
	s := NewPayload()
	d := NewPayload()

	s.Raw = []byte("abcdefg")
	s.Data = []string{"0", "1", "3"}
	d.Data = []string{"0000"}
	d, _ = s.Copy(d, s)
	fmt.Println(d.Raw)
	fmt.Println(d.Data)

	s.Data = map[string]string{"zero": "0", "one": "1", "three": "3"}
	d.Data = make(map[string]string)
	d, _ = s.Copy(d, s)
	fmt.Println(d.Data)
}
func TestPayloadAppend(t *testing.T) {
	s := NewPayload()
	d := NewPayload()

	s.Raw = []byte("abcdefg")
	s.Data = []string{"0", "1", "3"}
	d.Data = []string{"0000", "1111", "2222", "3333"}
	d, _ = s.Append(d, s)
	fmt.Println(d.Raw)
	fmt.Println(d.Data)

	s.Data = map[string]string{"zero": "0", "one": "1", "three": "3"}
	d.Data = map[string]string{"aaaa": "a", "bbbb": "b"}
	d, _ = s.Append(d, s)
	fmt.Println(d.Data)
}

func TestArguments(t *testing.T) {
	graph := NewGraphline()

	fatalError(t, graph.RegisterStage(normalIOC), "Unable to register normalIOC")
	fatalError(t, graph.RegisterStage(urlIOC), "Unable to register urlIOC")
	fatalError(t, graph.RegisterStage(ipIOC), "Unable to register ipIOC")
	fatalError(t, graph.RegisterStage(hashIOC), "Unable to register hashIOC")
	fatalError(t, graph.RegisterStage(jsonout), "Unable to register jsonout")
	seq, err := graph.Sequence(`graph1:normalIOC{"a":"123", "b":"456"}|ipIOC{"a":"vala" , "b":"valb"}|hashIOC|jsonout`)
	fatalError(t, err, "Unable to create graph with arguments")
	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)
	fatalError(t, graph.Execute(seq[0], payload), "Unable to execute graph")
	fmt.Println(string(payload.Raw))
}

func TestMatrix(t *testing.T) {
	graph := NewGraphline()

	fatalError(t, graph.RegisterStage(SelectFields), "Unable to register SelectFields")
	seq, err := graph.Sequence(`graph1:SelectFields{"fields":"0 2 4"}`)
	fatalError(t, err, "Unable to create graph with arguments")
	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw, []byte("1,2,3,4,5\n01,02,03,04,05\n001,002,003,004,005")...)
	err = graph.Execute(seq[0], payload)
	fmt.Println(string(payload.Raw))
}

func TestGraph(t *testing.T) {
	EventInit()
	graph := NewGraphline()

	fatalError(t, graph.RegisterStage(normalIOC), "Unable to register normalIOC")
	fatalError(t, graph.RegisterStage(urlIOC), "Unable to register urlIOC")
	fatalError(t, graph.RegisterStage(ipIOC), "Unable to register ipIOC")
	fatalError(t, graph.RegisterStage(hashIOC), "Unable to register hashIOC")
	fatalError(t, graph.RegisterStage(jsonout), "Unable to register jsonout")
	fatalError(t, graph.RegisterStage(func1), "Unable to register func1")
	fatalError(t, graph.RegisterStage(func2), "Unable to register func2")
	fatalError(t, graph.RegisterStage(filterIOC), "Unable to register filterIOC")
	graphid, err := graph.Sequence(`graph1:normalIOC|urlIOC|filterIOC;normalIOC|ipIOC|filterIOC;normalIOC|hashIOC|filterIOC;filterIOC|func1;filterIOC|func2;func1|jsonout;func2|jsonout`)
	fatalError(t, err, "Unable to create graph")

	start := NewAspect(StartTime, time.Duration(10*time.Millisecond), 0, CONTINUE)
	end := NewAspect(EndTime, time.Duration(10*time.Millisecond), 0, CONTINUE)
	duration := NewAspect(DurationTime, time.Duration(10*time.Millisecond), 0, CONTINUE)
	trace := NewAspect(Trace, time.Duration(10*time.Millisecond), 0, CONTINUE)
	graph.RegisterAspect(graphid[0], START, nil, &trace)
	graph.RegisterAspect(graphid[0], START, normalIOC, &start)
	graph.RegisterAspect(graphid[0], END, jsonout, &end)
	graph.RegisterAspect(graphid[0], END, jsonout, &duration)

	if graphid != nil {
		graph.PrintPath(graphid[0])
		payload := new(Payload)
		payload.Raw = make([]byte, 0, 2048)
		fatalError(t, graph.Execute(graphid[0], payload), "Unable to execute graph")
		fmt.Println(string(payload.Raw))

		time.Sleep(2000 * time.Millisecond)
	}
}

func TestSubgraph(t *testing.T) {
	EventInit()
	graph := NewGraphline()
	fatalError(t, graph.RegisterStage(func1), "Unable to register func1")
	fatalError(t, graph.RegisterStage(func2), "Unable to register func2")
	fatalError(t, graph.RegisterStage(func3), "Unable to register func3")
	fatalError(t, graph.RegisterStage(func4), "Unable to register func4")
	fatalError(t, graph.RegisterStage(func5), "Unable to register func5")
	fatalError(t, graph.RegisterStage(SubGraph), "Unable to register SubGraph")

	graph1, err := graph.Sequence(`graph1:func1|func2|SubGraph{"name":"graph2"}|func3`)
	fatalError(t, err, "Unable to create graph 1")
	graph2, err := graph.Sequence(`graph2:func4|func5`)
	fatalError(t, err, "Unable to create graph 2")

	trace := NewAspect(Trace, time.Duration(10*time.Millisecond), 0, CONTINUE)
	graph.RegisterAspect(graph1[0], START, nil, &trace)
	graph.RegisterAspect(graph2[0], START, nil, &trace)

	graph.PublishGraphs()

	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw, []byte("Peter Piper Picked a peck of pickled peppers")...)

	fatalError(t, graph.Execute(graph1[0], payload), "Unable to execute graph1")
	time.Sleep(1000 * time.Millisecond)

}

func TestFile(t *testing.T) {
	EventInit()
	graph := NewGraphline()
	fatalError(t, graph.RegisterStage(func1), "Unable to register func 1")
	fatalError(t, graph.RegisterStage(func2), "Unable to register func 2")
	fatalError(t, graph.RegisterStage(func3), "Unable to register func 3")
	fatalError(t, graph.RegisterStage(FileWrite), "Unable to register func 4")
	fatalError(t, graph.RegisterStage(FileRead), "Unable to register func 5")

	seq1, err := graph.Sequence(`graph1:FileRead{"name":"filein.txt"}|func1|func2|func3|FileWrite{"name":"fileout.txt"}`)
	fatalError(t, err, "Unable to create sequence 1")

	trace := NewAspect(Trace, time.Duration(10*time.Millisecond), 0, CONTINUE)
	start := NewAspect(StartTime, time.Duration(10*time.Millisecond), 0, CONTINUE)
	end := NewAspect(EndTime, time.Duration(10*time.Millisecond), 0, CONTINUE)
	duration := NewAspect(DurationTime, time.Duration(10*time.Millisecond), 0, CONTINUE)

	graph.RegisterAspect(seq1[0], START, nil, &trace)
	graph.RegisterAspect(seq1[0], START, FileRead, &start)
	graph.RegisterAspect(seq1[0], END, FileWrite, &end)
	graph.RegisterAspect(seq1[0], END, FileWrite, &duration)

	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)

	fatalError(t, graph.Execute(seq1[0], payload), "Unable to execute sequence 1")
	time.Sleep(1000 * time.Millisecond)
}

func TestFileStream(t *testing.T) {
	EventInit()
	graph := NewGraphline()
	fatalError(t, graph.RegisterStage(FileTextReader), "Unable to register func 1")
	fatalError(t, graph.RegisterStage(FileTextWriter), "Unable to register func 2")
	g, err := graph.Sequence(`graph1:FileTextReader{"name":"filein.txt"}|FileTextWriter{"name":"fileout.txt"}`)
	fatalError(t, err, "Unable to create sequence 1")
	trace := NewAspect(Trace, time.Duration(10*time.Millisecond), 0, CONTINUE)
	graph.RegisterAspect(g[0], START, nil, &trace)
	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)

	for {
		if err = graph.Execute(g[0], payload); err != nil {
			break
		}
	}

	time.Sleep(2000 * time.Millisecond)
}

func TestLambdaGraph(t *testing.T) {
	EventInit()
	graph := NewGraphline()

	fatalError(t, graph.RegisterStage(AWSDownloadS3Bucket), "Unable to register AWSDownloadS3Bucket")
	fatalError(t, graph.RegisterStage(NormalIOC), "Unable to register NormalIOC")
	fatalError(t, graph.RegisterStage(UrlIOC), "Unable to register urlIOC")
	fatalError(t, graph.RegisterStage(Ipv4IOC), "Unable to register Ipv4IOC")
	fatalError(t, graph.RegisterStage(Ipv6IOC), "Unable to register Ipv4IOC")
	fatalError(t, graph.RegisterStage(Md5IOC), "Unable to register Md5IOC")
	fatalError(t, graph.RegisterStage(Sha1IOC), "Unable to register Sha1IOC")
	fatalError(t, graph.RegisterStage(Sha256IOC), "Unable to register Sha256IOC")
	fatalError(t, graph.RegisterStage(IOCtoData), "Unable to register IOCtoData")
	fatalError(t, graph.RegisterStage(IOCDataToJson), "Unable to register IOCDataToJson")
	graphid, err := graph.Sequence(`lambda:AWSDownloadS3Bucket{"region":"us-east-2","bucket":"lambdapipeprocessor","file":"bulkdata.txt"}|NormalIOC|UrlIOC|IOCtoData;NormalIOC|Ipv4IOC|IOCtoData;NormalIOC|Ipv6IOC|IOCtoData;NormalIOC|Md5IOC|IOCtoData;NormalIOC|Sha1IOC|IOCtoData;NormalIOC|Sha256IOC|IOCtoData;IOCtoData|IOCDataToJson`)
	fatalError(t, err, "Unable to create graph")

	if graphid != nil {
		graph.PrintPath(graphid[0])
		payload := new(Payload)
		payload.Raw = make([]byte, 0, 2048)
		fatalError(t, graph.Execute(graphid[0], payload), "Unable to execute graph")
		fmt.Println(string(payload.Raw))

		time.Sleep(2000 * time.Millisecond)
	}
}
