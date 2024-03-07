package controller

type Controller interface {
	LoadTCPWithProxy(proxyTarget, target string) error
	LoadTCP(target string) error
	LoadHTTPWithProxy(proxyTarget, target string) error
	LoadHTTP(target string) error
}
