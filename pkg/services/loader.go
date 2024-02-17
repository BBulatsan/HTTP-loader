package services

import (
	"HTTP-loader/internal/config"
	"HTTP-loader/internal/repository"
	"HTTP-loader/pkg/controller"
	"log"
	"sync"
	"sync/atomic"
)

type Loader struct {
	cfg            config.LoaderConfig
	loadController controller.Controller
	proxiesReader  repository.ProxiesReader
	success        atomic.Uint64
	blocked        atomic.Uint64

	proxySemaphore  chan struct{}
	workerSemaphore chan struct{}
}

func NewLoader(
	config config.LoaderConfig,
	loadController controller.Controller,
	proxiesReader repository.ProxiesReader) Loader {
	return Loader{
		cfg:            config,
		loadController: loadController,
		proxiesReader:  proxiesReader,

		proxySemaphore:  make(chan struct{}, config.MaxNumProxiesRead),
		workerSemaphore: make(chan struct{}, config.MaxNumRequests),
	}
}

func (l *Loader) LoaderStart() error {
	proxies, err := l.proxiesReader.ReadProxiesFromFile()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	if l.cfg.TargetHTTP != "" {
		for _, proxyTarget := range proxies {
			l.proxySemaphore <- struct{}{}
			wg.Add(1)
			go func(proxyTarget string) {
				log.Printf("HTTP:Start with ProxyTarget %s\n", proxyTarget)
				l.workerHTTP(1, proxyTarget)
				wg.Done()
			}(proxyTarget)
		}
	}

	if l.cfg.TargetTCP != "" {
		for _, proxyTarget := range proxies {
			l.proxySemaphore <- struct{}{}
			wg.Add(1)
			go func(proxyTarget string) {
				log.Printf("TCP:Start with ProxyTarget %s\n", proxyTarget)
				l.workerTCP(1, proxyTarget)
				wg.Done()
			}(proxyTarget)
		}
	}

	wg.Wait()
	log.Printf("All proxies has finished! \n Request success: %d \n Request blocked: %d \n",
		l.success.Load(), l.blocked.Load())

	return nil
}

func (l *Loader) workerHTTP(scale uint64, proxyTarget string) error {
	var err error
	var wg sync.WaitGroup

	for i := 0; i < int(scale); i++ {
		go func() {
			wg.Add(1)
			l.workerSemaphore <- struct{}{}
			err = l.loadController.LoadHTTP(proxyTarget, l.cfg.TargetHTTP)
			if err == nil {
				l.success.Add(1)
			}
			<-l.workerSemaphore
			wg.Done()
		}()
	}
	wg.Wait()
	if err != nil {
		log.Printf("HTTP:ProxyTarget %s has been blocked with err: %s\n", proxyTarget, err)
		<-l.proxySemaphore
		l.blocked.Add(1)

		return err
	}

	return l.workerHTTP(scale*2, proxyTarget)
}

func (l *Loader) workerTCP(scale uint64, proxyTarget string) error {
	var err error
	var wg sync.WaitGroup

	for i := 0; i < int(scale); i++ {
		go func() {
			wg.Add(1)
			l.workerSemaphore <- struct{}{}
			err = l.loadController.LoadTCP(proxyTarget, l.cfg.TargetHTTP)
			if err == nil {
				l.success.Add(1)
			}
			<-l.workerSemaphore
			wg.Done()
		}()
	}
	wg.Wait()
	if err != nil {
		log.Printf("HTTP:ProxyTarget %s has been blocked with err: %s\n", proxyTarget, err)
		<-l.proxySemaphore
		l.blocked.Add(1)

		return err
	}

	return l.workerTCP(scale*2, proxyTarget)
}
