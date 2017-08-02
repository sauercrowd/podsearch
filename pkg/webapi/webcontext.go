package webapi

import (
	"database/sql"

	"github.com/sauercrowd/podsearch/pkg/algolia"
)

// WebContext delivers environment for web handlers
type WebContext struct {
	DBConn  *sql.DB
	Algolia *algolia.PodcastAlgolia
}
