package webapi

import (
	"database/sql"
)

// WebContext delivers environment for web handlers
type WebContext struct {
	DBConn *sql.DB
}
