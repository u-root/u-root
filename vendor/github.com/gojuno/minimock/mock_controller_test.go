package minimock

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	c := NewController(t)
	assert.Equal(t, &safeTester{Tester: t}, c.Tester)
}

func TestController_RegisterMocker(t *testing.T) {
	c := &Controller{}
	c.RegisterMocker(nil)
	assert.Len(t, c.mockers, 1)
}

type dummyMocker struct {
	finishCounter int32
	waitCounter   int32
}

func (dm *dummyMocker) MinimockFinish() {
	atomic.AddInt32(&dm.finishCounter, 1)
}

func (dm *dummyMocker) MinimockWait(time.Duration) {
	atomic.AddInt32(&dm.waitCounter, 1)
}

func TestController_Finish(t *testing.T) {
	dm := &dummyMocker{}
	c := &Controller{
		mockers: []Mocker{dm, dm},
	}

	c.Finish()
	assert.Equal(t, int32(2), atomic.LoadInt32(&dm.finishCounter))
}

func TestController_Wait(t *testing.T) {
	dm := &dummyMocker{}
	c := &Controller{
		mockers: []Mocker{dm, dm},
	}

	c.Wait(0)
	assert.Equal(t, int32(2), atomic.LoadInt32(&dm.waitCounter))
}

func TestController_WaitConcurrent(t *testing.T) {
	um1 := &unsafeMocker{}
	um2 := &unsafeMocker{}

	c := &Controller{
		Tester:  newSafeTester(&unsafeTester{}),
		mockers: []Mocker{um1, um2},
	}

	um1.tester = c
	um2.tester = c

	c.Wait(0) //shouln't produce data races
}

type unsafeMocker struct {
	Mocker
	tester Tester
}

func (um *unsafeMocker) MinimockWait(time.Duration) {
	um.tester.FailNow()
}

type unsafeTester struct {
	Tester

	finished bool
}

func (u *unsafeTester) FailNow() {
	u.finished = true
}
