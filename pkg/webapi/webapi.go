package webapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type WebHandler struct {
	*WebContext
	Handler func(*WebContext, http.ResponseWriter, *http.Request) (int, error)
}

func (wh WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := wh.Handler(wh.WebContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func sendJSON(data interface{}, w http.ResponseWriter) (int, error) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error while marshaling podcast: %v", err)
	}
	w.Header().Add("content-type", "application/json")
	w.Write(b)
	return 200, nil
}

func RegisterRoutes(ctx *WebContext) {
	http.HandleFunc("/api/v1/addpodcast", WebHandler{WebContext: ctx, Handler: AddPodcastHandler}.ServeHTTP)
}
