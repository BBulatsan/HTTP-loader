package main

import (
	"HTTP-loader/internal/config"
	"HTTP-loader/internal/repository"
	"HTTP-loader/pkg/controller"
	"HTTP-loader/pkg/services"
	"HTTP-loader/pkg/services/loader"
	"HTTP-loader/pkg/services/loader_v2"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	lConfig := config.NewConfig()

	proxiesReader := repository.NewProxiesReader(lConfig.ProxyFile)
	loadController := controller.NewLoadController()

	var loaderService services.Loader
	switch lConfig.Version {
	case 1:
		loaderService = loader.NewLoader(lConfig, &loadController, proxiesReader)
	case 2:
		loaderService = loader_v2.NewLoader(lConfig, &loadController, proxiesReader)
	default:
		log.Fatalf("Err: uncorrect version!")
	}

	err := loaderService.LoaderStart(ctx)
	if err != nil {
		log.Fatalf("Err: %s\n", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	if !lConfig.UseProxy {
		loaderService.LoaderStat()
	}

	cancel()
}
