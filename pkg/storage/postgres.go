package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	//necessary for postgres
	_ "github.com/lib/pq"
	"github.com/sauercrowd/podsearch/pkg/flags"
)

//Setup creates everything necessary(like tables), so the program is ready to operate
func Setup(f flags.Flags) (*sql.DB, error) {
	dbstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", f.Puser, f.Ppassword, f.Phost, f.Pport, f.Pdb)
	db, err := sql.Open("postgres", dbstr)
	if err != nil {
		log.Fatalf("Could create database connection: %v", err)
	}
	err = db.Ping()
	//if the wait flag is present, try until it works
	for err != nil && f.Wait {
		log.Printf("Waiting for database @%s:%d...", f.Phost, f.Pport)
		err = db.Ping()
		time.Sleep(time.Second * 3)
	}
	// otherwise return the error
	if err != nil {
		return nil, err
	}
	if err := createDatabaseIfNotExist(f); err != nil {
		return nil, err
	}
	if err := setupPodcastTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

func createDatabaseIfNotExist(f flags.Flags) error {
	dbstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", f.Puser, f.Ppassword, f.Phost, f.Pport, f.Puser)
	db, err := sql.Open("postgres", dbstr)
	if err != nil {
		log.Fatalf("Could create database connection: %v", err)
		return err
	}
	//check if database exists
	var count int64
	err = db.QueryRow("SELECT COUNT(1) FROM pg_database WHERE datname = $1", f.Pdb).Scan(&count)
	//return if database exists or error happend
	if err != nil || count == 1 {
		if err == nil {
			err = db.Close()
		}
		return err
	}
	err = db.QueryRow(fmt.Sprintf("CREATE DATABASE %s", f.Pdb)).Scan()
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Could not create database %s: %v", f.Pdb, err)
		return err
	}
	err = db.Close()
	return err
}
