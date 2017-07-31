package webapi

import (
	"log"
	"net/http"

	"github.com/sauercrowd/podsearch/pkg/podcast"
)

func AddPodcastHandler(ctx *WebContext, w http.ResponseWriter, r *http.Request) (int, error) {
	log.Printf("Context: %v", ctx)
	channel, err := podcast.AddPodcastFromURL("http://feeds.feedburner.com/GcpPodcast?format=xml")
	if err != nil {
		log.Fatal("Error")
	}
	return sendJSON(channel, w)
}
