package search

import (
	"fmt"
	"log"
	"time"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/sauercrowd/podsearch/pkg/flags"
	"gopkg.in/olivere/elastic.v5"
)

type Config struct {
	AlgoliaClient *algoliasearch.Client
	ElasticClient *elastic.Client
}

func New(flags flags.Flags) (*Config, error) {
	ac := algoliasearch.NewClient(flags.AlgoliaID, flags.AlgoliaKey)
	elasticURL := elastic.SetURL(fmt.Sprintf("http://%s:%s@%s:%d", flags.ElasticUser, flags.ElasticPassword, flags.ElasticHost, flags.ElasticPort))
	ec, err := elastic.NewClient(elasticURL)
	//if err and err is node not available and wait, try to reach elasticsearch as long as needed
	for err != nil && err.Error() == "health check timeout: no Elasticsearch node available" && flags.Wait {
		log.Println("Trying to reach elasticsearch...")
		ec, err = elastic.NewClient(elasticURL)
		time.Sleep(time.Second * 3)
	}
	//if it should wait or there's another error, return it
	if err != nil {
		return nil, err
	}
	return &Config{AlgoliaClient: &ac, ElasticClient: ec}, nil
}
