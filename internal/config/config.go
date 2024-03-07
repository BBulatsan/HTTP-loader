package config

import "github.com/alexflint/go-arg"

type Config struct {
	TargetHTTP        string `arg:"--targetHTTP" help:"http target"`
	TargetTCP         string `arg:"--targetTCP" help:"tcp target"`
	MaxNumProxiesRead uint64 `default:"100" arg:"--maxProxyRead" help:"max batch read of proxies"`
	MaxNumRequests    uint64 `default:"100" arg:"--rps" help:"max rps"`
	Version           uint8  `default:"2" arg:"-v" help:"version of loader 1 or 2"`
	UseProxy          bool   `default:"false" arg:"-p" help:"switch from two regime of work"`
	ProxyFile         string `default:"socks5_proxies.txt" arg:"--file" help:"way to file proxy socks5 format on field ip:port"`
}

func NewConfig() Config {
	cfg := Config{
		TargetHTTP: "",
		TargetTCP:  "",
	}

	arg.MustParse(&cfg)

	return cfg
}
