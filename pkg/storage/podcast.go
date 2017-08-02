package storage

import (
	"database/sql"
	"log"

	"github.com/sauercrowd/podsearch/pkg/podcast"
)

func setupPodcastTables(db *sql.DB) error {
	err := db.QueryRow("CREATE TABLE IF NOT EXISTS podcasts(title text NOT NULL, url text PRIMARY KEY NOT NULL, language text, description text, feedurl text, imageurl text)").Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	err = db.QueryRow("CREATE TABLE IF NOT EXISTS podcastepisodes(podcasturl text references podcasts(url) NOT NULL, title text NOT NULL, url text NOT NULL, pubdate date, description text, audiourl text)").Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func insertPodcastToDB(db *sql.DB, p *podcast.Channel) error {
	err := db.QueryRow("INSERT INTO podcasts(title, url, language, description, feedurl, imageurl) VALUES($1,$2,$3,$4,$5,$6)",
		p.Title, p.Link, p.Language, p.Description, p.FeedURL, p.ImageURL).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, e := range p.Episodes {
		if err := insertPodcastEpisodeToDB(db, &e, p.Link); err != nil && err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}

func insertPodcastEpisodeToDB(db *sql.DB, p *podcast.Episode, podcasturl string) error {
	err := db.QueryRow("INSERT INTO podcastepisodes(podcasturl, title, url, pubdate, description, audiourl) VALUES ($1, $2, $3, $4, $5, $6)",
		podcasturl, p.Title, p.Link, p.PubDate, p.Description, p.AudioURL).Scan()
	return err
}

//AddPodcast adds a new podcast to the database if it doesn't exist
func AddPodcast(db *sql.DB, p *podcast.Channel) (int, error) {
	//check if podcast already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM podcasts WHERE url=$1", p.Link).Scan(&count)
	if err != nil {
		return -1, err
	}
	if count > 0 {
		log.Printf("Podcast %s already exists", p.Title)
		return -1, nil
	}
	return 0, insertPodcastToDB(db, p)
}

func getEpisodesForPodcastURL(db *sql.DB, url string) ([]podcast.Episode, error) {
	pepisodes := make([]podcast.Episode, 0)
	rows, err := db.Query("SELECT title, url, pubdate, description, audiourl FROM podcastepisodes WHERE podcasturl=$1", url)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e podcast.Episode
		rows.Scan(&e.Title, &e.Link, &e.PubDate, &e.Description, &e.AudioURL)
		pepisodes = append(pepisodes, e)
	}
	return pepisodes, nil
}

//GetPodcast returns a podcast based on it's url if it does exist
func GetPodcast(db *sql.DB, url string) (*podcast.Channel, error) {
	var p podcast.Channel
	err := db.QueryRow("SELECT title, url, language, description, feedurl, imageurl FROM podcasts WHERE url=$1", url).Scan(&p)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	pepisodes, err := getEpisodesForPodcastURL(db, url)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	p.Episodes = pepisodes
	return &p, nil
}

//GetPodcasts returns all saved podcasts
func GetPodcasts(db *sql.DB) (*[]podcast.Channel, error) {
	pods := make([]podcast.Channel, 0)
	rows, err := db.Query("SELECT title, url, language, description, feedurl, imageurl FROM podcasts")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	for rows.Next() {
		var p podcast.Channel
		rows.Scan(&p.Title, &p.Link, &p.Language, &p.Description, &p.FeedURL, &p.ImageURL)
		pepisodes, err := getEpisodesForPodcastURL(db, p.Link)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		p.Episodes = pepisodes
		pods = append(pods, p)
	}
	return &pods, nil
}
