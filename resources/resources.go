package resources

import (
	"context"
	"golang.org/x/sync/errgroup"
	"stash.tutu.ru/golang/log"
)

type Resources struct {
	Env    *Env
	Config *Config
}

func Get(ctx context.Context) *Resources {

	r := &Resources{}

	if err := r.initEnv(); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error init env")
	}
	var err error
	err = r.initEtcd()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Init Etcd failed")
	}

	group, ctx := errgroup.WithContext(ctx)
	if err := group.Wait(); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error init resources")
	}
	return r
}
