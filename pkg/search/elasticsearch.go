package search

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/sauercrowd/podsearch/pkg/podcast"
)

func (c *Config) episodesToElasticsearch(p podcast.Channel, indexName string) error {
	for i, e := range p.Episodes {
		ee := SearchEpisode{
			Title:              e.Title,
			Description:        e.Description,
			Podcast:            p.Title,
			ImageURL:           p.ImageURL,
			Language:           p.Language,
			PodcastDescription: p.Description,
			PubDate:            e.PubDate,
			AudioURL:           e.AudioURL,
			Link:               e.Link,
			FeedURL:            p.FeedURL,
		}
		//do not care about created object
		_, err := c.ElasticClient.Index().
			Index(indexName).
			Type("episode").
			Id(strconv.Itoa(i)).
			BodyJson(ee).
			Do(context.Background())
		if err != nil {
			log.Println("Could not add podcast episode to elasticsearch: ", err)
			return err
		}
	}
	_, err := c.ElasticClient.Flush().Index("podcasts").Do(context.Background())
	if err != nil {
		log.Println("Could not flush episodes to elasticsearch: ", err)
		return err
	}
	return nil
}

func (c *Config) createElasticseachIndex(indexName string) error {
	exists, err := c.ElasticClient.IndexExists(indexName).Do(context.Background())
	if err != nil {
		log.Println("Could not check podcast index on elasticsearch: ", err)
		return err
	}
	if exists {
		return nil
	}
	createIndex, err := c.ElasticClient.CreateIndex(indexName).Body(elasticsearchIndexConfig).Do(context.Background())
	if err != nil {
		log.Println("Could not create podcast index on elasticsearch", err)
		return err
	}
	if !createIndex.Acknowledged {
		log.Println("Index creation not acknowledged")
		return fmt.Errorf("Index creation not acknowledged")
	}
	return nil
}

func (c *Config) elasticSearch(term string, indexName string) ([]SearchEpisode, error) {
	termQuery := elastic.NewTermQuery("description", term)
	searchResult, err := c.ElasticClient.Search().
		Index(indexName).
		Query(termQuery).
		//Sort("podcastTitle", true).      // sort by "user" field, ascending
		From(0).Size(100).       // max 100 documents
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		return nil, err
	}
	elasticEpisodes := make([]SearchEpisode, 0)
	var ea SearchEpisode
	for _, item := range searchResult.Each(reflect.TypeOf(ea)) {
		t := item.(SearchEpisode)
		elasticEpisodes = append(elasticEpisodes, t)
	}
	return elasticEpisodes, nil
}

var elasticsearchIndexConfig = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"_default_": {
			"_all": {
				"enabled": true
			}
		},
		"episode":{
			"properties":{
				"podsearch-title":{
					"type":"text"
				},
				"podsearch-description":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"podsearch-podcast":{
					"type":"text"
				},
				"podsearch-imageurl":{
					"type":"text"
				},
				"podsearch-feedurl":{
					"type":"text"
				},
				"podsearch-language":{
					"type":"text"
				},
				"podsearch-podcastDescription":{
					"type":"text"
				},
				"podsearch-pubdate":{
					"type":"text"
				},
				"podsearch-audiourl":{
					"type":"text"
				},
				"podsearch-link":{
					"type":"text"
				}
			}
		}
	}
}
`
