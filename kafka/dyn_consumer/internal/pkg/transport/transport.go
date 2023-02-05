package transport

import (
	"sync"
	"time"
)

type Transport interface {
	Start()
	Stop()
}

func TransportController(transports ...Transport) (stopFn func(duration time.Duration)) {
	var wg sync.WaitGroup

	for _, t := range transports {
		wg.Add(1)
		go t.Start()
	}

	return func(duration time.Duration) {
		for _, t := range transports {
			go func(t Transport) {
				t.Stop()
				wg.Done()
			}(t)
		}

		select {
		case <-time.After(duration):
			return
		case <-func() chan struct{} {
			c := make(chan struct{})
			go func() {
				wg.Wait()
				close(c)
			}()
			return c
		}():
			return
		}
	}
}
