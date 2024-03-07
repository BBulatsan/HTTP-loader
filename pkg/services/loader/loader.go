package loader

import (
	"HTTP-loader/internal/config"
	"HTTP-loader/internal/repository"
	"HTTP-loader/pkg/controller"
	"context"
	"log"
	"sync/atomic"
	"time"
)

type Loader struct {
	cfg            config.Config
	loadController controller.Controller
	proxiesReader  repository.ProxiesReader
	successHTTP    atomic.Uint64
	blockedHTTP    atomic.Uint64
	successTCP     atomic.Uint64
	blockedTCP     atomic.Uint64

	timeStart time.Time

	proxySemaphore  chan struct{}
	workerSemaphore chan struct{}
}

func NewLoader(
	config config.Config,
	loadController controller.Controller,
	proxiesReader repository.ProxiesReader) *Loader {
	return &Loader{
		cfg:            config,
		loadController: loadController,
		proxiesReader:  proxiesReader,

		proxySemaphore:  make(chan struct{}, config.MaxNumProxiesRead),
		workerSemaphore: make(chan struct{}, config.MaxNumRequests),
	}
}

func (l *Loader) LoaderStart(ctx context.Context) error {
	if l.cfg.UseProxy {
		err := l.loaderStartWithProxy(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (l *Loader) LoaderStat() {
	spendTime := time.Now().Sub(l.timeStart)

	if l.cfg.TargetHTTP != "" {
		log.Printf("HTTP:Request success: %d Request blocked: %d \n",
			l.successHTTP.Load(), l.blockedHTTP.Load())
		log.Printf("HTTP:Spend time: %v Average requests: %f rps",
			spendTime, float64(l.successHTTP.Load()+l.blockedHTTP.Load())/spendTime.Seconds())
	}

	if l.cfg.TargetTCP != "" {
		log.Printf("TCP:Request success: %d Request blocked: %d \n",
			l.successTCP.Load(), l.blockedTCP.Load())
		log.Printf("TCP:Spend time: %v Average requests: %f rps",
			spendTime, float64(l.successTCP.Load()+l.blockedTCP.Load())/spendTime.Seconds())
	}
}
