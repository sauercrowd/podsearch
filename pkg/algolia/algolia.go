package algolia

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/sauercrowd/podsearch/pkg/podcast"
)

type PodcastAlgolia struct {
	AlgoliaClient *algoliasearch.Client
}

func New(appID string, adminKey string) *PodcastAlgolia {
	c := algoliasearch.NewClient(appID, adminKey)
	return &PodcastAlgolia{AlgoliaClient: &c}
}

func (pa *PodcastAlgolia) AddPodcast(podcast *podcast.Channel) {
	podcastIndex := (*pa.AlgoliaClient).InitIndex("podcasts")
	//	episodeIndex := (*pa.AlgoliaClient).InitIndex("episodes")
	//	algoliapodcast := podcastToAlgolia(podcast)
	//	podcastIndex.AddObject(algoliapodcast)
	algoliaepisodes := episodeToAlgolia(*podcast, podcast.Episodes)
	podcastIndex.AddObjects(algoliaepisodes)
}

func episodeToAlgolia(p podcast.Channel, episode []podcast.Episode) []algoliasearch.Object {
	m := make([]algoliasearch.Object, 0)
	for _, e := range episode {
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
		m = append(m, ae)
	}
	return m
}

func podcastToAlgolia(podcast *podcast.Channel) algoliasearch.Object {
	return algoliasearch.Object{
		"title":       podcast.Title,
		"description": podcast.Description,
		"imageurl":    podcast.ImageURL,
		"url":         podcast.Link,
		"feedurl":     podcast.FeedURL,
		"language":    podcast.Language,
	}
}
