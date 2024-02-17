package controller

type Controller interface {
	LoadTCP(proxyTarget, target string) error
	LoadHTTP(proxyTarget, target string) error
}
