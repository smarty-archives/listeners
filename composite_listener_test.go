package listeners

import (
	"testing"
	"time"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestCompositeListenerFixture(t *testing.T) {
	gunit.Run(new(CompositeListenerFixture), t)
}

type CompositeListenerFixture struct {
	*gunit.Fixture

	listeners []Listener
	composite *CompositeListener
}

func (this *CompositeListenerFixture) Setup() {
	for x := 0; x < 100; x++ {
		this.listeners = append(this.listeners, &FakeForCompositeListener{})
	}
	this.composite = NewCompositeListener(this.listeners...)
}

func (this *CompositeListenerFixture) TestCompositeListenerCallsInnerListenersConcurrently() {
	started := utcNow()
	this.composite.Listen()
	this.So(time.Since(started), should.BeLessThan, nap*5)
}

////////////////////////////////////////////////////////////////////////////////

func (this *CompositeListenerFixture) TestCompositeListenerDoesntFailWithNoListeners() {
	this.listeners = nil
	this.composite = NewCompositeListener(this.listeners...)
	this.So(this.composite.Listen, should.NotPanic)
}

////////////////////////////////////////////////////////////////////////////////

func (this *CompositeListenerFixture) TestCompositeListenerSkipNilListeners() {
	this.listeners = append(this.listeners, &FakeForCompositeListener{})
	this.listeners = append(this.listeners, nil)
	this.listeners = append(this.listeners, nil)
	this.composite = NewCompositeListener(this.listeners...)
	this.So(this.composite.Listen, should.NotPanic)
	this.So(this.composite.Close, should.NotPanic)
}

////////////////////////////////////////////////////////////////////////////////

func (this *CompositeListenerFixture) TestCloseCallsInnerListeners() {
	items := []Listener{&FakeForCompositeListener{}, &FakeForCompositeListener{}}

	NewCompositeListener(items...).Close()

	for _, item := range items {
		this.So(item.(*FakeForCompositeListener).closeCalls, should.Equal, 1)
	}
}

////////////////////////////////////////////////////////////////////////////////

var nap = time.Millisecond

type FakeForCompositeListener struct{ closeCalls int }

func (this *FakeForCompositeListener) Listen() {
	time.Sleep(nap)
}

func (this *FakeForCompositeListener) Close() {
	this.closeCalls++
}

// //////////////////////////////////////////////////////////////////////////////
func utcNow() time.Time {
	return time.Now().UTC()
}
