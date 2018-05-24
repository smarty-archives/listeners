package listeners

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/smartystreets/logging"
)

type ShutdownListener struct {
	logger *logging.Logger

	mutex    sync.Once
	channel  chan os.Signal
	shutdown func()
}

func NewShutdownListener(shutdown func()) *ShutdownListener {
	channel := make(chan os.Signal, 16)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	return &ShutdownListener{channel: channel, shutdown: shutdown}
}

func (this *ShutdownListener) Listen() {
	if message := <-this.channel; message != nil {
		this.logger.Printf("[INFO] Received application shutdown signal [%s].\n", message)
	}

	this.shutdown()
}

func (this *ShutdownListener) Close() {
	this.mutex.Do(this.close)
}

func (this *ShutdownListener) close() {
	signal.Stop(this.channel)
	close(this.channel)
	this.logger.Println("[INFO] Unsubscribed from OS shutdown signals.")
}
