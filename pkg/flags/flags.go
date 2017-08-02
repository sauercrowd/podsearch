package flags

import "flag"

type flags struct {
	Phost, Puser, Ppassword, Pdb, AlgoliaID, AlgoliaKey string
	Port, Pport                                         int
}

func ParseFlags() *flags {
	phost := flag.String("phost", "127.0.0.1", "Postgres host")
	puser := flag.String("puser", "postgres", "Postgres user")
	ppassword := flag.String("ppassword", "postgres", "Postgres password")
	pdb := flag.String("pdatabase", "podsearch", "Postgres database")
	pport := flag.Int("pport", 5432, "Postgres port")
	port := flag.Int("port", 8080, "Port on the API should listen on")
	algoliaID := flag.String("algoliaid", "", "Algolia AppID")
	algoliaKey := flag.String("algoliakey", "", "Algolia Admin Key")
	flag.Parse()
	return &flags{
		Phost:      *phost,
		Puser:      *puser,
		Ppassword:  *ppassword,
		Pport:      *pport,
		Port:       *port,
		Pdb:        *pdb,
		AlgoliaID:  *algoliaID,
		AlgoliaKey: *algoliaKey,
	}
}
