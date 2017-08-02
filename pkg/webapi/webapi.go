package webapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// WebHandler is a special handler for http.Handle to pass additional context to handler funtions
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

// RegisterRoutes creates all routes and applies it to http.Handle
func RegisterRoutes(ctx *WebContext) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/podcast", WebHandler{WebContext: ctx, Handler: addPodcastHandler}.ServeHTTP).Methods("POST")
	r.HandleFunc("/api/v1/podcast", WebHandler{WebContext: ctx, Handler: getPodcastHandler}.ServeHTTP).Methods("GET").Queries("url", "{url:.}")
	r.HandleFunc("/api/v1/podcast", WebHandler{WebContext: ctx, Handler: getPodcastsHandler}.ServeHTTP).Methods("GET")

	http.Handle("/", r)
	listRoutes(r)
}

func listRoutes(r *mux.Router) {
	log.Println("Available Routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		p, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		log.Println(strings.Join(m, ","), t, p)
		return nil
	})
}

type podResponse struct {
	Err bool   `json:"error"`
	Msg string `json:"msg"`
}

func sendSimpleResponse(rw http.ResponseWriter, msg string, success bool) (int, error) {
	r := podResponse{!success, msg}
	b, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("Error while marshaling podcast: %v", err)
		return http.StatusInternalServerError, nil
	}

	rw.Header().Add("content-type", "application/json")
	rw.Write(b)
	return http.StatusOK, nil
}
