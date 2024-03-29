package listeners

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestCompositeWaitListenerFixture(t *testing.T) {
	gunit.Run(new(CompositeWaitListenerFixture), t)
}

type CompositeWaitListenerFixture struct {
	*gunit.Fixture

	completed time.Time
	listener  *CompositeWaitListener
	items     []Listener
}

func (this *CompositeWaitListenerFixture) Setup() {
	log.SetOutput(ioutil.Discard)
	this.items = []Listener{&FakeListener{}, &FakeListener{}}
	this.listener = NewCompositeWaitListener(this.items...)
}

func (this *CompositeListenerFixture) Teardown() {
	log.SetOutput(os.Stderr)
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestAllListenersAreCalledAndWaitedFor() {
	this.listener.Listen()

	this.completed = utcNow()

	for _, item := range this.items {
		if item == nil {
			continue
		}

		this.So(item.(*FakeListener).calls, should.Equal, 1)
		this.So(this.completed.After(item.(*FakeListener).instant), should.BeTrue)
	}
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestNilListenersDontCausePanic() {
	this.listener = NewCompositeWaitListener(nil, nil, nil)
	this.So(this.listener.Listen, should.NotPanic)
	this.So(this.listener.Close, should.NotPanic)
}

//////////////////////////////////////////

func (this *CompositeWaitListenerFixture) TestCloseCallsInnerListeners() {
	this.listener.Close()

	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 1)
	}
}

func (this *CompositeWaitListenerFixture) TestMultipleCloseCallInnerListenersExactlyOnce() {
	this.listener.Close()
	this.listener.Close()

	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 1)
	}
}

func (this *CompositeWaitListenerFixture) TestCloseDoesntInvokeInfiniteLoop() {
	this.listener = NewCompositeWaitShutdownListener(this.items...)

	go this.listener.Close()
	this.listener.Listen()

	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 1)
	}
}

func (this *CompositeWaitListenerFixture) TestDelayedCloseDoesDelay() {
	delay := 50 * time.Millisecond
	this.listener = NewCompositeWaitDelayedShutdownListener(delay, this.items...)

	start := time.Now()
	go this.listener.Listen()
	time.Sleep(time.Millisecond) //Listen() needs to execute before kill or we'll get instant fail (empty WaitGroup)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	//pre-delay
	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 0)
	}

	this.listener.waiter.Wait()
	elapsed := time.Since(start)
	this.So(elapsed.Nanoseconds(), should.BeGreaterThanOrEqualTo, delay.Nanoseconds())

	//post-delay
	for _, item := range this.items {
		this.So(item.(*FakeListener).closeCalls, should.Equal, 1)
	}
}

//////////////////////////////////////////

type FakeListener struct {
	calls      int
	closeCalls int
	instant    time.Time
}

func (this *FakeListener) Listen() {
	this.instant = utcNow()
	time.Sleep(time.Millisecond)
	this.calls++
}

func (this *FakeListener) Close() {
	this.closeCalls++
}

//////////////////////////////////////////
