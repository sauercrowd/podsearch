package webapi

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func searchHandler(ctx *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	term := mux.Vars(r)["term"]
	results, err := ctx.Search.SearchEpisodes(term)
	if err != nil {
		log.Println("Could not search for epsiodes: ", err)
		return sendSimpleResponse(w, "Internal server error", false)
	}
	return sendJSON(results, w)
}

func searchWebsocketHandler(ctx *Context, w http.ResponseWriter, r *http.Request) (int, error) {
	return 0, nil
}
