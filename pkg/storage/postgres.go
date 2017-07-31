package storage

import "database/sql"

//Setup creates everything necessary(like tables), so the program is ready to operate
func Setup(db *sql.DB) error {
	if err := setupPodcastTables(db); err != nil {
		return err
	}
	return nil
}
