package resources

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	//"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"stash.tutu.ru/golang/log"
	//"stash.tutu.ru/golang/resources/db"
)

type Resources struct {
	Env 	*Env
	config  *Config
	Db 		*sql.DB
	//Db            *db.Database
}

func Get(ctx context.Context) *Resources {

	r := &Resources{}

	if err := r.initEnv(); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error init env")
	}
	var err error
	err = r.initEtcd()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Init etcd failed")
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return errors.Wrap(r.InitDb(), "init db")
	})
	if err := group.Wait(); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error init resources")
	}
	return r
}
