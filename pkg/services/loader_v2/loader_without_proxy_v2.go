package loader_v2

import (
	"context"
	"log"
	"time"
)

func (l *Loader) loaderStartWithoutProxy(ctx context.Context) error {
	for i := 0; i <= int(l.cfg.MaxNumRequests); i++ {
		if l.cfg.TargetHTTP != "" {
			go l.workerHTTP(ctx)
		}
		if l.cfg.TargetTCP != "" {
			go l.workerTCP(ctx)
		}
	}

	l.timeStart = time.Now()

	return nil
}

func (l *Loader) workerHTTP(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			go func() {
				err := l.loadController.LoadHTTP(l.cfg.TargetHTTP)
				if err != nil {
					log.Printf("HTTP: has been blocked with err: %s\n", err)
					l.blockedHTTP.Add(1)

					return
				}
				l.successHTTP.Add(1)
			}()
		case <-ctx.Done():
			ticker.Stop()

			return
		}
	}
}

func (l *Loader) workerTCP(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			go func() {
				err := l.loadController.LoadTCP(l.cfg.TargetTCP)
				if err != nil {
					log.Printf("TCP: has been blocked with err: %s\n", err)
					l.blockedTCP.Add(1)

					return
				}
				l.successTCP.Add(1)
			}()
		case <-ctx.Done():
			ticker.Stop()

			return
		}
	}
}
