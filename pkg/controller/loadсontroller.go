package controller

import (
	"HTTP-loader/pkg/clients/httpclient"
	"HTTP-loader/pkg/clients/tcpclient"
	"log"
	"net/http"
)

type LoadController struct {
}

func NewLoadController() LoadController {
	return LoadController{}
}

func (l *LoadController) LoadTCP(proxyTarget, target string) error {
	client, err := tcpclient.NewTCPClientWithProxy(proxyTarget, target)
	if err != nil {
		return err
	}
	n, err := client.Write([]byte{})
	if err != nil {
		return err
	}

	log.Printf("Sent tcp bytes : %d\n", n)

	err = client.Close()
	if err != nil {
		return err
	}

	return nil
}

func (l *LoadController) LoadHTTP(proxyTarget, target string) error {
	client, err := httpclient.NewHTTPClientWithProxy(proxyTarget)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Printf("Sent http with status : %s\n", resp.Status)

	_ = resp.Close

	client.CloseIdleConnections()

	return nil
}
