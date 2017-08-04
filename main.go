package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sauercrowd/podsearch/pkg/flags"
	"github.com/sauercrowd/podsearch/pkg/search"
	"github.com/sauercrowd/podsearch/pkg/storage"
	"github.com/sauercrowd/podsearch/pkg/webapi"
)

func main() {
	flags := flags.ParseFlags()

	db, err := storage.Setup(flags)
	if err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}

	s, err := search.New(flags)
	if err != nil {
		log.Fatalf("Could not initialize search: %v", err)
	}
	context := &webapi.Context{DBConn: db, Search: s}
	webapi.RegisterRoutes(context)

	//http.Handle("/", http.FileServer(http.Dir("./static"))) //TODO: Decide if handled by go server or by seperate nginx instance
	log.Printf("Server listening on :%d", flags.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flags.Port), nil))
}
