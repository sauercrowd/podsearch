package webapi

import (
	"database/sql"

	"github.com/sauercrowd/podsearch/pkg/search"
)

// Context delivers environment for web handlers
type Context struct {
	DBConn *sql.DB
	Search *search.Config
}
