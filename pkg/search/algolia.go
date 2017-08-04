package search

import (
	"log"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/sauercrowd/podsearch/pkg/podcast"
)

func (c *Config) episodeToAlgolia(p podcast.Channel) error {
	algIndex := (*c.AlgoliaClient).InitIndex("podcasts")
	algoliaepisodes := make([]algoliasearch.Object, 0)
	for _, e := range p.Episodes {
		ae := algoliasearch.Object{
			"podcast":            p.Title,
			"imageurl":           p.ImageURL,
			"feedurl":            p.FeedURL,
			"language":           p.Language,
			"podcastDescription": p.Description,
			"title":              e.Title,
			"description":        e.Description,
			"pubdate":            e.PubDate,
			"audiourl":           e.AudioURL,
			"link":               e.Link,
		}
		algoliaepisodes = append(algoliaepisodes, ae)
	}
	//batchRes not needed
	_, err := algIndex.AddObjects(algoliaepisodes)
	if err != nil {
		log.Println("Could not add podcast to algolia: ", err)
		return err
	}
	return nil
}
