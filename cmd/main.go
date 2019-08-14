package main

import (
	"context"
	"go.uber.org/automaxprocs/maxprocs"
	"stash.tutu.ru/golang/context_os"
	"stash.tutu.ru/golang/envs"
	"stash.tutu.ru/golang/http-server"
	"stash.tutu.ru/golang/log"
	_ "stash.tutu.ru/golang/opentracing"
	"stash.tutu.ru/golang/readiness"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/handlers"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/resources"
)

func main() {

	_, err := maxprocs.Set(maxprocs.Min(2), maxprocs.Logger(log.Logger.Info().Msgf))
	if err != nil {
		log.Logger.Error().Err(err).Msg("CPU quota error")
	}

	envs.UpdateDotenv()
	//envs.PrintKubernetesEnv()

	ctx := context_os.Context(context.Background())

	res := resources.Get(ctx)

	ready := readiness.New()
	ready.AddProbe(func() {
		//add some initialization
	})

	h := handlers.New(res)

	s := server.NewServer(ready)
	s.HandleFunc("/test", h.Test)

	if err := s.Start(ctx); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error http server")
	}
}
