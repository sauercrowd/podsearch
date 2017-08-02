package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sauercrowd/podsearch/pkg/algolia"
	"github.com/sauercrowd/podsearch/pkg/flags"
	"github.com/sauercrowd/podsearch/pkg/storage"
	"github.com/sauercrowd/podsearch/pkg/webapi"
)

func main() {
	flags := flags.ParseFlags()

	dbstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", flags.Puser, flags.Ppassword, flags.Phost, flags.Pport, flags.Pdb)
	db, err := sql.Open("postgres", dbstr)
	if err != nil {
		log.Fatalf("Could not open connection to database: %v", err)
	}
	if err := storage.Setup(db); err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}

	alg := algolia.New(flags.AlgoliaID, flags.AlgoliaKey)
	context := &webapi.WebContext{DBConn: db, Algolia: alg}
	webapi.RegisterRoutes(context)

	//http.Handle("/", http.FileServer(http.Dir("./static"))) //TODO: Decide if handled by go server or by seperate nginx instance
	log.Printf("Server listening on :%d", flags.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flags.Port), nil))
}
