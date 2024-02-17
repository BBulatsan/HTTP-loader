package httpclient

import (
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

func NewHTTPClientWithProxy(proxyTarget string) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyTarget, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		Dial:               dialer.Dial,
		MaxIdleConns:       10,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: true,
	}

	return &http.Client{Transport: tr, Timeout: 60 * time.Second}, nil
}
