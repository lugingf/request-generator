package resources

import (
	"github.com/coreos/etcd/client"
	"stash.tutu.ru/golang/log"
	"time"
)

type Config struct {
	Etcd client.KeysAPI
}


func (r *Resources) initEtcd() error {

	cfg := client.Config{
		Endpoints:               []string{"http://" + r.Env.Etcd},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		return err
	}

	kapi := client.NewKeysAPI(c)

	config := &Config{Etcd: kapi}

	if r.Config == nil {
		log.Logger.Info().Msg("initEtcd success")
	}
	r.Config = config
	return nil
}
