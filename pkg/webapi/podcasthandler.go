package webapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sauercrowd/podsearch/pkg/storage"

	"github.com/sauercrowd/podsearch/pkg/podcast"
)

type addPodcastBody struct {
	URL string `json:"url"`
}

// addPodcastHandler creates a new podcast, if none exists
func addPodcastHandler(ctx *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("AddPocastHandler: %v", err)
		return http.StatusInternalServerError, nil
	}
	var body addPodcastBody
	if err := json.Unmarshal(bytes, &body); err != nil {
		log.Printf("AddPocastHandler: %v", err)
		return http.StatusBadRequest, nil
	}
	p, err := podcast.AddPodcastFromURL(body.URL)
	if err != nil {
		log.Printf("AddPocastHandler: %v", err)
		return sendSimpleResponse(w, "invalid content at url", false)
	}
	code, err := storage.AddPodcast(ctx.DBConn, p)
	if err != nil {
		log.Printf("AddPocastHandler: %v", err)
		return http.StatusInternalServerError, nil
	}
	if code == -1 {
		return sendSimpleResponse(w, "podcast already exists", false)
	}
	ctx.Search.AddPodcast(*p)
	return sendSimpleResponse(w, "podcast added", true)
}

// getPodcastHandler returns a saved podcast based on its url, if available
func getPodcastHandler(ctx *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	return 1, nil
}

// getPodcastsHandler returns a saved podcast based on its url, if available
func getPodcastsHandler(ctx *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	podcasts, err := storage.GetPodcasts(ctx.DBConn)
	if err != nil {
		log.Printf("AddPocastHandler: %v", err)
		return http.StatusInternalServerError, nil
	}
	return sendJSON(podcasts, w)
}
