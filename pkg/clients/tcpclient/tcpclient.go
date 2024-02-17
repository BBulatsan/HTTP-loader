package tcpclient

import (
	"net"
	"time"

	"golang.org/x/net/proxy"
)

func NewTCPClientWithProxy(proxyTarget, target string) (net.Conn, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyTarget, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	conn, err := dialer.Dial("tcp", target)
	if err != nil {
		return nil, err
	}

	err = conn.SetDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
