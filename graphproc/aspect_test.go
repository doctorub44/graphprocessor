package graphproc

import (
	"testing"
	"time"
)

func TestNewAspect(t *testing.T) {

	dur := time.Duration(10 * time.Millisecond)

	NewAspect(testaspect, dur, 0, RETRY)
}

func testaspect(n string, m []byte, scope *Scope) error {
	return nil
}

func TestTrace(t *testing.T) {
	Level1()
}

func Level1() {
	Level2()
}

func Level2() {
	Trace("name", nil, nil)
}
