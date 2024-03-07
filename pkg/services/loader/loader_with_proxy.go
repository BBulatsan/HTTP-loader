package loader

import (
	"context"
	"log"
	"sync"
	"time"
)

func (l *Loader) loaderStartWithProxy(_ context.Context) error {
	proxies, err := l.proxiesReader.ReadProxiesFromFile()
	if err != nil {
		return err
	}

	l.timeStart = time.Now()

	var wg sync.WaitGroup
	if l.cfg.TargetHTTP != "" {
		for _, proxyTarget := range proxies {
			l.proxySemaphore <- struct{}{}
			wg.Add(1)
			go func(proxyTarget string) {
				defer wg.Done()
				log.Printf("HTTP:Start with ProxyTarget %s\n", proxyTarget)
				l.scalableWorkerHTTP(1, proxyTarget)
			}(proxyTarget)
		}
	}

	if l.cfg.TargetTCP != "" {
		for _, proxyTarget := range proxies {
			l.proxySemaphore <- struct{}{}
			wg.Add(1)
			go func(proxyTarget string) {
				defer wg.Done()
				log.Printf("TCP:Start with ProxyTarget %s\n", proxyTarget)
				l.scalableWorkerTCP(1, proxyTarget)
			}(proxyTarget)
		}
	}

	wg.Wait()

	l.LoaderStat()

	return nil
}

func (l *Loader) scalableWorkerHTTP(scale uint64, proxyTarget string) error {
	var err error
	var wg sync.WaitGroup

	for i := 0; i < int(scale); i++ {
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
		<-l.proxySemaphore
		l.blockedHTTP.Add(1)

		return err
	}

	return l.scalableWorkerHTTP(scale*2, proxyTarget)
}

func (l *Loader) scalableWorkerTCP(scale uint64, proxyTarget string) error {
	var err error
	var wg sync.WaitGroup

	for i := 0; i < int(scale); i++ {
		l.workerSemaphore <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = l.loadController.LoadTCPWithProxy(proxyTarget, l.cfg.TargetHTTP)
			if err == nil {
				l.successTCP.Add(1)
			}
			<-l.workerSemaphore
		}()
	}
	wg.Wait()
	if err != nil {
		log.Printf("TCP:ProxyTarget %s has been blocked with err: %s\n", proxyTarget, err)
		<-l.proxySemaphore
		l.blockedTCP.Add(1)

		return err
	}

	return l.scalableWorkerTCP(scale*2, proxyTarget)
}
