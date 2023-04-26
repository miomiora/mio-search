package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"mio-search/config"
)

var (
	es  *elastic.Client
	ctx context.Context
)

func InitElasticSearch(cfg *config.ElasticSearchConfig) (err error) {
	es, err = elastic.NewSimpleClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)),
		elastic.SetBasicAuth(cfg.User, cfg.Password),
	)
	ctx = context.Background()
	return
}
