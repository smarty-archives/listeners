package listeners

import (
	"testing"
	"time"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestSerialListenerFixture(t *testing.T) {
	gunit.Run(new(SerialListenerFixture), t)
}

type SerialListenerFixture struct {
	*gunit.Fixture
}

func (this *SerialListenerFixture) TestListenCallInOrder() {
	items := []Listener{
		&FakeForSerialListener{listened: utcNow().Add(time.Second)},
		&FakeForSerialListener{listened: utcNow()},
		&FakeForSerialListener{listened: utcNow().Add(-time.Second)},
	}

	NewSerialListener(items...).Listen()

	times := []time.Time{}
	for _, item := range items {
		fake := item.(*FakeForSerialListener)
		times = append(times, fake.listened)

		this.So(fake.calls, should.Equal, 1)
	}

	this.So(times, should.BeChronological)

}

func (this *SerialListenerFixture) TestNilListenersAreIgnored() {
	this.So(NewSerialListener(nil).Listen, should.NotPanic)
	this.So(NewSerialListener(nil).Close, should.NotPanic)
}

func (this *SerialListenerFixture) TestCloseCalledOnInnerListeners() {
	items := []Listener{&FakeForSerialListener{}, &FakeForSerialListener{}}

	NewSerialListener(items...).Close()

	for _, item := range items {
		this.So(item.(*FakeForSerialListener).closeCalls, should.Equal, 1)
	}
}

type FakeForSerialListener struct {
	calls      int
	closeCalls int
	listened   time.Time
}

func (this *FakeForSerialListener) Listen() {
	this.calls++
	this.listened = utcNow()
	time.Sleep(time.Microsecond)
}

func (this *FakeForSerialListener) Close() {
	this.closeCalls++
}
