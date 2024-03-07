package repository

type ReadProxies interface {
	ReadProxiesFromFile() ([]string, error)
}
