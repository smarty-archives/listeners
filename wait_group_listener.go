package listeners

type WaitGroupListener struct {
	inner  Listener
	waiter WaitGroup
}

func NewWaitGroupListener(listener Listener, waiter WaitGroup) *WaitGroupListener {
	waiter.Add(1)
	return &WaitGroupListener{
		inner:  listener,
		waiter: waiter,
	}
}

func (this *WaitGroupListener) Listen() {
	if this.inner != nil {
		this.inner.Listen()
	}
	this.waiter.Done()
}

func (this *WaitGroupListener) Close() {
	if closer, ok := this.inner.(ListenCloser); ok {
		closer.Close()
	}
}
