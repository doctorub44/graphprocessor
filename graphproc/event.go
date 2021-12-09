package graphproc

import (
	"time"
)

//Gevent : graphline events for logging, observability, metrics, etc.
type Gevent struct {
	id     peventid
	strval string
	//ival int
	//fval float
	tstamp time.Time
}

type peventid int64

//Graphline event IDs
const (
	MESSAGE peventid = iota + 1
	FMETRIC
	IMETRIC
)

var events chan Gevent

//EventInit : initiate the event channel and start a go routine to process events
func EventInit() {
	events = make(chan Gevent, 10000)
	go ProcessEvent()
}

//SendEvent : send event the the event channel
func SendEvent(e Gevent) {
	e.tstamp = time.Now()
	events <- e
}

//ProcessEvent : every second, process events from the event channel
func ProcessEvent() {
	for {
		time.Sleep(1 * time.Second)
		for e := range events {
			switch e.id {
			case MESSAGE:
				FastLogger(e.strval)
			case IMETRIC:
			case FMETRIC:
			}
		}
	}
}
