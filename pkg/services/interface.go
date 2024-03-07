package services

import "context"

type Loader interface {
	LoaderStart(ctx context.Context) error
	LoaderStat()
}
