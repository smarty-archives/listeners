package listeners

import (
	"sync"
)

type CascadingWaitListener struct {
	waiter    WaitGroup
	listeners []ListenCloser
}

func NewCascadingWaitListener(listeners ...Listener) *CascadingWaitListener {
	this := &CascadingWaitListener{waiter: &sync.WaitGroup{}}

	for _, listener := range listeners {
		this.listeners = append(this.listeners, NewCloseWaitListener(listener))
	}

	return this
}

func (this *CascadingWaitListener) Listen() {
	this.waiter.Add(len(this.listeners))
	for _, listener := range this.listeners {
		go this.listen(listener)
	}
	this.waiter.Wait()
}
func (this *CascadingWaitListener) listen(listener Listener) {
	listener.Listen()
	this.waiter.Done()
}

func (this *CascadingWaitListener) Close() {
	for _, listener := range this.listeners {
		listener.Close()
	}
}
