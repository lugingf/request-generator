package resources

import (
	"github.com/kelseyhightower/envconfig"
	"stash.tutu.ru/golang/envs"
	"stash.tutu.ru/golang/log"
)

type Env struct {
	ServiceName string `envconfig:"APP_SERVICENAME" default:"request-generator"`
	PodName     string `envconfig:"APP_PODNAME" default:"podname"`
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
