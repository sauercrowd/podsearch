package storage

import (
	"database/sql"
	"log"

	"github.com/sauercrowd/podsearch/pkg/podcast"
)

func setupPodcastTables(db *sql.DB) error {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS podcasts(title text NOT NULL, url text PRIMARY KEY NOT NULL, language text, description text, imageurl text)")
	if err != nil {
		return err
	}
	_, err = db.Query("CREATE TABLE IF NOT EXISTS podcastepisodes(podcasturl text references podcast(url) NOT NULL, title text NOT NULL, url text NOT NULL, pubdate date, description text, audiourl text)")
	if err != nil {
		return err
	}
	return nil
}

func insertPodcastToDB(db *sql.DB, p *podcast.PodcastChannel) error {
	err := db.QueryRow("INSERT INT podcast(title, url, language, description, imageurl) VALUES($1,$2,$3,$4,$5)",
		p.Title, p.Link, p.Language, p.Description, p.ImageURL).Scan()
	if err != nil {
		return err
	}
	for _, e := range p.Episodes {
		if err := insertPodcastEpisodeToDB(db, &e, p.Link); err != nil {
			return err
		}
	}
	return nil
}

func insertPodcastEpisodeToDB(db *sql.DB, p *podcast.PodcastEpisode, podcasturl string) error {
	err := db.QueryRow("INSERT INT podcastepisodes(podcasturl, title, url, pubdate, description, audiourl) VALUES ($1, $2, $3, $4, $5, $6",
		podcasturl, p.Title, p.Link, p.PubDate, p.Description, p.AudioURL).Scan()
	return err
}

//AddPodcast adds a new podcast to the database if it doesn't exist
func AddPodcast(db *sql.DB, p *podcast.PodcastChannel) (int, error) {
	//check if podcast already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM podcast WHERE url=$1", p.Link).Scan(&count)
	if err != nil {
		return -1, err
	}
	if count > 0 {
		log.Printf("Podcast %s already exists", p.Title)
		return -1, nil
	}
	return 0, insertPodcastToDB(db, p)
}

func getEpisodesForPodcastURL(db *sql.DB, url string) ([]podcast.PodcastEpisode, error) {
	pepisodes := make([]podcast.PodcastEpisode, 0)
	rows, err := db.Query("SELECT title, url, pubdate, description, audiourl FROM podcastepisodes WHERE podcasturl=$1", url)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e podcast.PodcastEpisode
		rows.Scan(&e)
		pepisodes = append(pepisodes, e)
	}
	return pepisodes, nil
}

//GetPodcast returns a podcast based on it's url if it does exist
func GetPodcast(db *sql.DB, url string) (*podcast.PodcastChannel, error) {
	var p podcast.PodcastChannel
	err := db.QueryRow("SELECT title, url, language, description, imageurl FROM podcast WHERE url=$1", url).Scan(&p)
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
func GetPodcasts(db *sql.DB) (*[]podcast.PodcastChannel, error) {
	pods := make([]podcast.PodcastChannel, 0)
	rows, err := db.Query("SELECT title, url, language, description, imageurl FROM podcast")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	for rows.Next() {
		var p podcast.PodcastChannel
		rows.Scan(&p)
		pepisodes, err := getEpisodesForPodcastURL(db, p.Link)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		p.Episodes = pepisodes
		pods = append(pods, p)
	}
	return &pods, nil
}
