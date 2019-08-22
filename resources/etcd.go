package resources

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/pkg/errors"
	"stash.tutu.ru/golang/log"
	"strings"
	"time"
)

type Config struct {
	etcd client.KeysAPI
	Db   ConfigDB
}

type ConfigDB struct {
	Driver string
	Host   string
	Port   string
	DBName string
	User   string
	Pass   string
}

func (c *Config) InitDB(key string) error {
	return nil;
	resp, err := c.etcd.Get(context.Background(), key, nil)
	if err != nil {
		return err
	}


	if resp.Node == nil {
		return errors.New(fmt.Sprintf("Not found node %s", key))
	}

	for _, node := range resp.Node.Nodes {
		switch {
		case strings.HasSuffix(node.Key, "user"):
			c.Db.User = node.Value
		case strings.HasSuffix(node.Key, "password"):
			c.Db.Pass = node.Value
		case strings.HasSuffix(node.Key, "dsn"):
			dns := strings.Split(node.Value, ":")
			c.Db.Driver = dns[0]
			params := strings.Split(dns[1], ";")

			for _, param := range params {
				str := strings.Split(param, "=")
				switch str[0] {
				case "host":
					c.Db.Host = str[1]
				case "port":
					c.Db.Port = str[1]
				case "dbname":
					c.Db.DBName = str[1]
				}
			}
		}
	}
	return nil
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

	config := &Config{etcd: kapi}

	//kapi.Watcher()
	if err = config.InitDB("/config-tutu/" + r.Env.ServiceName + "/service/databases/default"); err != nil {
		return errors.New(fmt.Sprintf("get config db: %s", err.Error()))
	}

	if r.config == nil {
		log.Logger.Info().Msg("initEtcd success")
	}
	r.config = config
	return nil
}
