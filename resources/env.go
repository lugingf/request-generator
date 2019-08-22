package resources

import (
	"github.com/kelseyhightower/envconfig"
	"stash.tutu.ru/golang/envs"
	"stash.tutu.ru/golang/log"
)

type Env struct {
	ServiceName 		string `envconfig:"APP_SERVICENAME" default:"request-generator"`
	PodName     		string `envconfig:"APP_PODNAME" default:"podname"`
	Etcd      			string `envconfig:"APP_ETCD" default:"http-01.dock.cats.devel.tutu.ru:8000"`
	ElasticSearchUrlEnv string `envconfig:"ELASTICSEARCH_URL" default:"http://elastic-logs.devel.tutu.ru:9200"`
}

func init() {
	envs.Register("application", Env{})
}

func (r *Resources) initEnv() error {
	var s Env
	err := envconfig.Process("app", &s)
	if err != nil {
		return err
	}

	r.Env = &s
	log.Logger.Info().Msg("initEnv success")
	return nil
}
