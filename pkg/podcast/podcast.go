package podcast

import (
	"io/ioutil"
	"log"
	"net/http"
)

type PodcastEpisode struct {
	Title       string `json:"title"`
	Link        string `json:"url"`
	PubDate     string `json:"date"`
	Description string `json:"description"`
	AudioURL    string `json:"audioUrl"`
}

type PodcastChannel struct {
	Title       string           `json:"title"`
	Link        string           `json:"url"`
	Language    string           `json:"lang"`
	Description string           `json:"description"`
	ImageURL    string           `json:"imageUrl"`
	Episodes    []PodcastEpisode `json:"episodes"`
}

func AddPodcastFromURL(url string) (*PodcastChannel, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Could not load RSS: %v", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Could not read Response Body: %v", err)
		return nil, err
	}
	feed, err := parseFromBody(body)
	if err != nil {
		log.Fatalf("Could not parse podcast feed: %v", err)
		return nil, err
	}
	return createChannelFromXMLFeed(feed), nil
}

func createChannelFromXMLFeed(feed *podcastXMLFeed) *PodcastChannel {
	episodes := make([]PodcastEpisode, len(feed.Items), len(feed.Items))
	//transfer items to episodes
	for i := 0; i < len(episodes); i++ {
		episodes[i] = PodcastEpisode{
			feed.Items[i].Title,
			feed.Items[i].Link,
			feed.Items[i].PubDate,
			feed.Items[i].Description,
			feed.Items[i].AudioURL.URL,
		}
	}
	return &PodcastChannel{
		feed.Title,
		feed.Link,
		feed.Language,
		feed.Description,
		feed.ImageURL,
		episodes,
	}
}
