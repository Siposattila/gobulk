package kill

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	KillCtx context.Context
	cancel  context.CancelFunc
	once    sync.Once
)

func ListenForKill() {
	once.Do(func() {
		KillCtx, cancel = context.WithCancel(context.Background())
		killChan := make(chan os.Signal, 1)
		signal.Notify(killChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-killChan
			cancel()
		}()
	})
}
