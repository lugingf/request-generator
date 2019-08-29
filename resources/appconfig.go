package resources

import (
	"context"
	"errors"
	"github.com/coreos/etcd/client"
	"stash.tutu.ru/golang/log"
	"strings"
)

type ChangedNode struct {
	Node 	string
	Value 	string
}

type GeneratorConfig struct {
	Rpm       string
	IsEnabled string
	GenRequestTimeout string
	LogResponsesEnabled string
	UrlTargetList []string
}


func (r *Resources) GetConfig() *GeneratorConfig {
	etcdResp := r.getEtcdNodes()

	genConf := GeneratorConfig{}
	log.Logger.Info().Msg("Initialization ETCD config started")
	for _, node := range etcdResp.Node.Nodes {
		if strings.HasSuffix(node.Key, "isEnabled") {
			genConf.IsEnabled = node.Value
		}
		if strings.HasSuffix(node.Key, "LogResponsesEnabled") {
			genConf.LogResponsesEnabled = node.Value
		}
		if strings.HasSuffix(node.Key, "GenRequestTimeout") {
			genConf.GenRequestTimeout = node.Value
		}
		if strings.HasSuffix(node.Key, "rpm") {
			genConf.Rpm = node.Value
		}
		if strings.HasSuffix(node.Key, "UrlTargetList") {
			genConf.UrlTargetList = strings.Split(node.Value, ",")
		}
	}
	log.Logger.Info().Msg("Initialization ETCD config finished")

	return &genConf
}

func (r *Resources) UpdateConfig() chan ChangedNode {
	ch := make(chan ChangedNode)
	log.Logger.Info().Msg("ETCD config watcher starting")
	go func() {
		key := r.getConfigKey()
		watcher := r.Config.Etcd.Watcher(key, &client.WatcherOptions{Recursive: true})

		for {
			resp, err := watcher.Next(context.Background())

			if err != nil {
				log.Logger.Fatal().Err(err).Msg(error.Error(err))
			}

			splitedNodeKey := strings.Split(resp.Node.Key, "/")
			nodeKey := splitedNodeKey[len(splitedNodeKey)-1]
			log.Logger.Info().Msg("Action " + resp.Action + " Node: " + nodeKey + " Value: " + resp.Node.Value)
			ch <- ChangedNode{Node: nodeKey, Value: resp.Node.Value}
		}
	}()

	log.Logger.Info().Msg("ETCD config watcher started")
	return ch
}

func (r *Resources) getConfigValue(key string) (string, error) {
	etcdResp := r.getEtcdNodes()

	for _, node := range etcdResp.Node.Nodes {
		if strings.HasSuffix(node.Key, key) {
			return node.Value, nil
		}
	}

	return "", errors.New("Not found key: " + key)
}

func (r *Resources) getEtcdNodes() *client.Response {
	key:= r.getConfigKey()
	etcdResponse, err := r.Config.Etcd.Get(context.Background(), key, nil)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg(error.Error(err))
	}

	return etcdResponse
}

func (r *Resources) getConfigKey() string {
	return "/config-tutu/" + r.Env.ServiceName + "/service/generator"
}
