package config

type Config struct {
	Loader LoaderConfig
}

type LoaderConfig struct {
	TargetHTTP        string
	TargetTCP         string
	MaxNumProxiesRead uint64
	MaxNumRequests    uint64
}
