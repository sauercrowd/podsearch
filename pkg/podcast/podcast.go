package podcast

import (
	"io/ioutil"
	"log"
	"net/http"
)

//Episode is a episode belonging to a PodcastChannel
type Episode struct {
	Title       string `json:"title"`
	Link        string `json:"url"`
	PubDate     string `json:"date"`
	Description string `json:"description"`
	AudioURL    string `json:"audioUrl"`
}

//Channel is a Podcast with episodes
type Channel struct {
	Title       string    `json:"title"`
	Link        string    `json:"url"`
	Language    string    `json:"lang"`
	Description string    `json:"description"`
	FeedURL     string    `json:"feedUrl"`
	ImageURL    string    `json:"imageUrl"`
	Episodes    []Episode `json:"episodes"`
}

//AddPodcastFromURL takes a feed url and creates a podcast with episodes from it, if possible
func AddPodcastFromURL(url string) (*Channel, error) {
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
	return createChannelFromXMLFeed(url, feed), nil
}

func createChannelFromXMLFeed(url string, feed *podcastXMLFeed) *Channel {
	episodes := make([]Episode, len(feed.Items), len(feed.Items))
	//transfer items to episodes
	for i := 0; i < len(episodes); i++ {
		episodes[i] = Episode{
			feed.Items[i].Title,
			feed.Items[i].Link,
			feed.Items[i].PubDate,
			feed.Items[i].Description,
			feed.Items[i].AudioURL.URL,
		}
	}
	return &Channel{
		feed.Title,
		feed.Link,
		feed.Language,
		feed.Description,
		url,
		feed.ImageURL,
		episodes,
	}
}
