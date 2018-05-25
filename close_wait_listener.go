package listeners

import "sync"

type CloseWaitListener struct {
	inner  Listener
	waiter WaitGroup
}

func NewCloseWaitListener(listener Listener) *CloseWaitListener {
	return &CloseWaitListener{inner: listener, waiter: &sync.WaitGroup{}}
}

func (this *CloseWaitListener) Listen() {
	if this.inner == nil {
		return
	}

	this.waiter.Add(1)
	this.inner.Listen()
	this.waiter.Done()
}

func (this *CloseWaitListener) Close() {
	if closer, ok := this.inner.(ListenCloser); ok {
		closer.Close()
	}
	this.waiter.Wait()
}
