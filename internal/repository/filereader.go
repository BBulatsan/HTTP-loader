package repository

import (
	"bufio"
	"os"
)

type ProxiesReader struct {
}

func NewProxiesReader() ProxiesReader {
	return ProxiesReader{}
}

func (p *ProxiesReader) ReadProxiesFromFile() ([]string, error) {
	var res []string
	file, err := os.Open("socks5_proxies.txt")
	if err != nil {
		return res, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		res = append(res, line)
	}

	if err = scanner.Err(); err != nil {
		return res, err
	}

	if err = file.Close(); err != nil {
		return res, err
	}

	return res, nil
}
