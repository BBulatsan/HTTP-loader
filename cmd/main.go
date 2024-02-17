package main

import (
	"HTTP-loader/internal/config"
	"HTTP-loader/internal/repository"
	"HTTP-loader/pkg/controller"
	"HTTP-loader/pkg/services"
	"log"
)

func main() {
	lConfig := config.LoaderConfig{
		MaxNumProxiesRead: 100,
		MaxNumRequests:    5000,
		TargetHTTP:        "",
		TargetTCP:         "",
	}

	proxiesReader := repository.NewProxiesReader()
	loadController := controller.NewLoadController()
	loader := services.NewLoader(lConfig, &loadController, proxiesReader)
	err := loader.LoaderStart()
	if err != nil {
		log.Fatalf("Err: %s\n", err)
	}
}
