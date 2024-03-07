package repository

import (
	"bufio"
	"os"
)

type ProxiesReader struct {
	fileName string
}

func NewProxiesReader(fileName string) ProxiesReader {
	return ProxiesReader{fileName: fileName}
}

func (p *ProxiesReader) ReadProxiesFromFile() ([]string, error) {
	var res []string
	file, err := os.Open(p.fileName)
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
