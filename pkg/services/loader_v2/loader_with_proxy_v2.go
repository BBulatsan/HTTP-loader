package loader_v2

import (
	"context"
	"log"
	"sync"
	"time"
)

const checkWaitTime = 10 * time.Second

func (l *Loader) loaderStartWithProxy(ctx context.Context) error {
	for i := 0; i <= int(l.cfg.MaxNumProxiesRead); i++ {
		go l.workerReadProxy(ctx)
	}

	proxies, err := l.proxiesReader.ReadProxiesFromFile()
	if err != nil {
		return err
	}

	l.timeStart = time.Now()

	for _, proxy := range proxies {
		l.proxyChan <- proxy
	}

	ticker := time.NewTicker(checkWaitTime)
breakWait:
	for {
		select {
		case <-ticker.C:
			if len(l.workerSemaphore) == 0 {
				ticker.Stop()

				break breakWait
			}
		}
	}

	l.LoaderStat()

	return nil
}

func (l *Loader) workerReadProxy(ctx context.Context) {
	for {
		select {
		case proxyTarget := <-l.proxyChan:
			if l.cfg.TargetHTTP != "" {
				l.scalableWorkerHTTP(proxyTarget)
			}
			if l.cfg.TargetTCP != "" {
				l.scalableWorkerTCP(proxyTarget)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (l *Loader) scalableWorkerHTTP(proxyTarget string) {
	var wg sync.WaitGroup
	var err error
	scale := 1
	for {
		for i := 0; i <= scale; i++ {
			l.workerSemaphore <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = l.loadController.LoadHTTPWithProxy(proxyTarget, l.cfg.TargetHTTP)
				if err == nil {
					l.successHTTP.Add(1)
				}
				<-l.workerSemaphore
			}()
		}
		wg.Wait()
		if err != nil {
			log.Printf("HTTP:ProxyTarget %s has been blocked with err: %s\n", proxyTarget, err)
			l.blockedHTTP.Add(1)

			return
		}

		scale = scale * 2
	}
}

func (l *Loader) scalableWorkerTCP(proxyTarget string) {
	var wg sync.WaitGroup
	var err error
	scale := 1
	for {
		for i := 0; i <= scale; i++ {
			l.workerSemaphore <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = l.loadController.LoadTCPWithProxy(proxyTarget, l.cfg.TargetTCP)
				if err == nil {
					l.successTCP.Add(1)
				}
				<-l.workerSemaphore
			}()
		}
		wg.Wait()
		if err != nil {
			log.Printf("TCP:ProxyTarget %s has been blocked with err: %s\n", proxyTarget, err)
			l.blockedTCP.Add(1)

			return
		}

		scale = scale * 2
	}
}
