package resources

import (
	"context"
	//"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"stash.tutu.ru/golang/log"
	//"stash.tutu.ru/golang/resources/db"
	//"stash.tutu.ru/golang/resources/kafka"
	//"stash.tutu.ru/golang/resources/mongo"
)

type Resources struct {
	Env *Env
	//Db            *db.Database
	//KafkaConsumer *kafka.Kafka
	//KafkaProducer *kafka.SyncProducer
	//Mongo         *mongo.Collection
}

func Get(ctx context.Context) *Resources {

	r := &Resources{}

	if err := r.initEnv(); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error init env")
	}
	group, ctx := errgroup.WithContext(ctx)
	//group.Go(func() error {
	//	return errors.Wrap(r.initDB(), "init db")
	//})
	//group.Go(func() error {
	//	return errors.Wrap(r.initKafka(ctx), "init kafka")
	//})
	//group.Go(func() error {
	//	return errors.Wrap(r.initMongo(ctx), "init mongodb")
	//})
	if err := group.Wait(); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error init resources")
	}
	return r
}
