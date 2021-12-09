package graphproc

import (
	"time"
)

//Aspect :
type Aspect struct {
	aspect func(name string, mesg []byte, scope *Scope) error
	time   time.Duration
	value  int64
	action AspectAction
}

//Scope :
type Scope struct {
	err   error
	start time.Time
	end   time.Time
}

//AspectAction :
type AspectAction int

//Aspect action :
const (
	RETRY AspectAction = iota + 1
	CANCEL
	CONTINUE
	DEALDLINE
	TIMEOUT
	ERROR
	START
	END
)

//NewAspect :
func NewAspect(f func(string, []byte, *Scope) error, t time.Duration, v int64, action AspectAction) Aspect {
	var a Aspect
	a.aspect = f
	a.time = t
	a.value = v
	a.action = action
	return a
}

//Execute :
func (a *Aspect) Execute(name string, mesg []byte, scope *Scope) AspectAction {
	if a.aspect == nil {
		return ERROR
	}
	scope.err = a.aspect(name, mesg, scope)

	return a.action
}

//Trace : trace aspect that sends trace events
func Trace(name string, mesg []byte, scope *Scope) error {
	var message string
	var errmesg string

	if scope.err != nil {
		errmesg = scope.err.Error()
	}
	if errmesg != "" {
		message = " : error = " + errmesg
	}
	var event Gevent
	event.id = MESSAGE
	event.strval = "Trace: " + name + message

	SendEvent(event)

	return nil
}

//StartTime : aspect that sends start time as an event
func StartTime(name string, mesg []byte, scope *Scope) error {
	var event Gevent

	scope.start = time.Now()
	event.id = MESSAGE
	event.strval = "Start: " + scope.start.String()
	SendEvent(event)

	return nil
}

//EndTime : aspect that sends the end time as an event
func EndTime(name string, mesg []byte, scope *Scope) error {
	var event Gevent

	scope.end = time.Now()
	event.id = MESSAGE
	event.strval = "End: " + scope.end.String()
	SendEvent(event)

	return nil
}

//DurationTime : aspect to calculate the duration and send an event. Uses the times from StartTime and EndTime
func DurationTime(name string, mesg []byte, scope *Scope) error {
	var event Gevent

	event.id = MESSAGE
	event.strval = "Duration: " + scope.end.Sub(scope.start).String()
	SendEvent(event)

	return nil
}
