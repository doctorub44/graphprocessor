package graphproc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestIOC(t *testing.T) {
	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)
	iocs := make([]byte, 0, 2048)

	payload.Raw = append(payload.Raw, []byte("Abc 123 00236A2Ae558018ed13b5222ef1bd98700236A2Ae558018ed13b5222ef1bd987 bunch 12345600236A2Ae558018ed13b5222ef1bd987 of text 00236A2Ae558018ed13b5222ef1bd987 and some 123ABC")...)
	Md5IOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	fmt.Println(string(payload.Raw))
	payload.Raw = append(payload.Raw[:0], []byte("Abc 123 bunch 1234567800236A2Ae558018ed13b5222ef1bd987 of text 00236A2Ae558018ed13b5222ef1bd987 and some 123ABC")...)
	Sha1IOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	fmt.Println(string(payload.Raw))
	payload.Raw = append(payload.Raw[:0], []byte("Abc 123 10.11.12.13 00236A2Ae558018ed13b5222ef1bd98700236A2Ae558018ed13b5222ef1bd987 192.168.100.100 bunch 12345600236A2Ae558018ed13b5222ef1bd987 of text 00236A2Ae558018ed13b5222ef1bd987 and some 123ABC 1.1.1.1")...)
	Sha256IOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	fmt.Println(string(payload.Raw))
	payload.Raw = append(payload.Raw[:0], []byte("Abc 123 10.11.12.13 00236A2Ae558018ed13b5222ef1bd98700236A2Ae558018ed13b5222ef1bd987 192.168.100.100 bunch 12345600236A2Ae558018ed13b5222ef1bd987 of text 00236A2Ae558018ed13b5222ef1bd987 and some 123ABC 1.1.1.1")...)
	Ipv4IOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	fmt.Println(string(payload.Raw))
	payload.Raw = append(payload.Raw[:0], []byte("Abc and https://wwww.hackme.com pv6 is 2001:db8:3333:4444:5555:6666:7777:8888 and 2001:db8:3333:4444:CCCC:DDDD:EEEE:FFFF 123 10.11.12.13 00236A2Ae558018ed13b5222ef1bd98700236A2Ae558018ed13b5222ef1bd987 192.168.100.100 bunch 12345600236A2Ae558018ed13b5222ef1bd987 of text 00236A2Ae558018ed13b5222ef1bd987 and some 123ABC 1.1.1.1")...)
	Ipv6IOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	fmt.Println(string(payload.Raw))
	payload.Raw = append(payload.Raw[:0], []byte("Abc and https://wwww.hackme.com pv6 is 2001:db8:3333:4444:5555:6666:7777:8888 and 2001:db8:3333:4444:CCCC:DDDD:EEEE:FFFF 123 10.11.12.13 00236A2Ae558018ed13b5222ef1bd98700236A2Ae558018ed13b5222ef1bd987 192.168.100.100 bunch 12345600236A2Ae558018ed13b5222ef1bd987 of text 00236A2Ae558018ed13b5222ef1bd987 and some 123ABC 1.1.1.1")...)
	UrlIOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	fmt.Println(string(payload.Raw))
	payload.Raw = append(payload.Raw[:0], iocs...)
	IOCtoData(nil, payload)
	fmt.Println(payload.Data)
	IOCDataToJson(nil, payload)
	fmt.Println(string(payload.Raw))
}

func TestIOC2(t *testing.T) {
	var data []byte

	EventInit()
	graph := NewGraphline()

	graph.RegisterStage(NormalIOC)
	graph.RegisterStage(UrlIOC)
	graph.RegisterStage(Ipv4IOC)
	graph.RegisterStage(Ipv6IOC)
	graph.RegisterStage(Md5IOC)
	graph.RegisterStage(Sha1IOC)
	graph.RegisterStage(Sha256IOC)
	graph.RegisterStage(IOCDataToJson)
	graph.RegisterStage(IOCtoData)

	g, err := graph.Sequence("graph:NormalIOC|UrlIOC;NormalIOC|Ipv4IOC;NormalIOC|Ipv6IOC;NormalIOC|Md5IOC;NormalIOC|Sha1IOC;NormalIOC|Sha256IOC;UrlIOC|IOCtoData;Ipv4IOC|IOCtoData;Ipv6IOC|IOCtoData;Md5IOC|IOCtoData;Sha1IOC|IOCtoData;Sha256IOC|IOCtoData|IOCDataToJson")
	if err != nil {
		fatalError(t, err, "Unable to create graph")
	}
	if data, err = ioutil.ReadFile("sample.txt"); err != nil {
		fatalError(t, err, "Unable to read sample")
	}

	trace := NewAspect(Trace, time.Duration(10*time.Millisecond), 0, CONTINUE)
	graph.RegisterAspect(g[0], START, nil, &trace)

	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw[:0], data...)
	if err = graph.Execute(g[0], payload); err != nil {
		fatalError(t, err, "Unable to execute graph")
	}
	fmt.Println(string(payload.Raw))
	time.Sleep(2000 * time.Millisecond)
}

func TestNormalIOC(t *testing.T) {
	var data []byte
	var err error

	if data, err = ioutil.ReadFile("sample.txt"); err != nil {
		fatalError(t, err, "Unable to read sample")
	}
	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw[:0], data...)
	NormalIOC(nil, payload)
	fmt.Println(string(payload.Raw))
}

func TestJsonToIOC(t *testing.T) {
	payload := new(Payload)
	payload.Raw = make([]byte, 0, 2048)
	iocs := make([]byte, 0, 2048)

	payload.Raw = append(payload.Raw[:0], []byte("Abc and https://wwww.hackme.com pv6 is 2001:db8:3333:4444:5555:6666:7777:8888 and 2001:db8:3333:4444:CCCC:DDDD:EEEE:FFFF 123 10.11.12.13 00236A2Ae558018ed13b5222ef1bd98700236A2Ae558018ed13b5222ef1bd987 192.168.100.100 bunch 12345600236A2Ae558018ed13b5222ef1bd987 of text 00236A2Ae558018ed13b5222ef1bd987 and some 123ABC 1.1.1.1")...)
	Ipv4IOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	payload.Raw = append(payload.Raw[:0], []byte("Abc and https://wwww.hackme.com pv6 is 2001:db8:3333:4444:5555:6666:7777:8888")...)
	UrlIOC(nil, payload)
	iocs = append(iocs, payload.Raw...)
	payload.Raw = append(payload.Raw[:0], iocs...)
	IOCtoData(nil, payload)
	IOCDataToJson(nil, payload)
	fmt.Println(string(payload.Raw))
}

func TestSelectFields(t *testing.T) {
	var config interface{}

	payload := NewPayload()
	state := new(State)
	json.Unmarshal([]byte(`{"fields":"0 2 4"}`), &config)
	state.config = config
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw, []byte("1,2,3,4,5\n01,02,03,04,05\n001,002,003,004,005")...)
	SelectFields(state, payload)
	fmt.Println(string(payload.Raw))
}

func TestBulkSelectFields(t *testing.T) {
	var config interface{}
	data, _ := ioutil.ReadFile("bulkdata.txt")
	payload := NewPayload()
	state := new(State)
	json.Unmarshal([]byte(`{"fields":"10"}`), &config)
	state.config = config
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw, data...)
	SelectFields(state, payload)
	Ipv4IOC(state, payload)
	fmt.Println(string(payload.Raw))
}

func TestCutFields(t *testing.T) {
	var config interface{}

	payload := NewPayload()
	state := new(State)
	json.Unmarshal([]byte(`{"fields":"0 2 4"}`), &config)
	state.config = config
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw, []byte("1,2,3,4,5\n01,02,03,04,05\n001,002,003,004,005")...)
	CutFields(state, payload)
	fmt.Println(string(payload.Raw))
}

func TestStageArguments(t *testing.T) {

	graph := NewGraphline()
	graph.RegisterStage(SelectFields)
	g, err := graph.Sequence(`graph:SelectFields{"fields":"2"}`)
	if err != nil {
		fatalError(t, err, "Unable to create graph")
	}
	payload := NewPayload()
	payload.Raw = make([]byte, 0, 2048)
	payload.Raw = append(payload.Raw, []byte("1,2,3,4,5\n01,02,03,04,05\n001,002,003,004,005")...)
	if err = graph.Execute(g[0], payload); err != nil {
		fatalError(t, err, "Unable to execute graph")
	}
	fmt.Println(string(payload.Raw))
}
