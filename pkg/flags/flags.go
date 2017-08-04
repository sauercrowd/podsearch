package flags

import "flag"
import "os"
import "strconv"

type Flags struct {
	Phost, Puser, Ppassword, Pdb, AlgoliaID, AlgoliaKey, ElasticHost, ElasticUser, ElasticPassword string
	Port, Pport, ElasticPort                                                                       int
	Wait                                                                                           bool
}

func ParseFlags() Flags {
	phost := flag.String("phost", "127.0.0.1", "Postgres host")
	puser := flag.String("puser", "postgres", "Postgres user")
	ppassword := flag.String("ppassword", "postgres", "Postgres password")
	pdb := flag.String("pdatabase", "podsearch", "Postgres database")
	pport := flag.Int("pport", 5432, "Postgres port")
	port := flag.Int("port", 8080, "Port on the API should listen on")
	algoliaID := flag.String("algoliaid", "", "Algolia AppID")
	algoliaKey := flag.String("algoliakey", "", "Algolia Admin Key")
	elasticHost := flag.String("elastichost", "127.0.0.1", "Elasticsearch host")
	elasticPort := flag.Int("elasticport", 9200, "Elasticsearch port")
	elasticUser := flag.String("elasticuser", "elastic", "Elasticsearch user")
	elasticPass := flag.String("elasticpassword", "changeme", "Elasticsearch pass")

	useEnv := flag.Bool("env", false, "Use Environment Variables instead of cmd line parameters; gets every other option from corresponding uppercase environment variables")
	wait := flag.Bool("wait", false, "Do not exit if database or search are not available, wait instead")
	flag.Parse()

	//if useEnv, grab variables from environment
	if *useEnv {
		f := ParseFromEnv()
		f.Wait = *wait
		return f
	}
	return Flags{
		Phost:           *phost,
		Puser:           *puser,
		Ppassword:       *ppassword,
		Pport:           *pport,
		Port:            *port,
		Pdb:             *pdb,
		AlgoliaID:       *algoliaID,
		AlgoliaKey:      *algoliaKey,
		ElasticHost:     *elasticHost,
		ElasticPort:     *elasticPort,
		ElasticUser:     *elasticUser,
		ElasticPassword: *elasticPass,
		Wait:            *wait,
	}
}

// ParseFromEnv parse the flags from environment
func ParseFromEnv() Flags {
	return Flags{
		Phost:           getFromEnv("PHOST", "127.0.0.1"),
		Puser:           getFromEnv("PUSER", "postgres"),
		Ppassword:       getFromEnv("PPASSWORD", "postgres"),
		Pdb:             getFromEnv("PDATABASE", "podsearch"),
		Pport:           getIntFromEnv("PPORT", 5432),
		Port:            getIntFromEnv("PORT", 8080),
		AlgoliaID:       getFromEnv("ALGOLIAID", ""),
		AlgoliaKey:      getFromEnv("ALGOLIAKEY", ""),
		ElasticHost:     getFromEnv("ELASTICHOST", "127.0.0.1"),
		ElasticPort:     getIntFromEnv("ELASTICPORT", 9200),
		ElasticUser:     getFromEnv("ELASTICUSER", "elastic"),
		ElasticPassword: getFromEnv("ELASTICPASSWORD", "changeme"),
	}
}

func getFromEnv(v string, def string) string {
	s := os.Getenv(v)
	if s == "" {
		return def
	}
	return s
}

func getIntFromEnv(v string, def int) int {
	s := os.Getenv(v)
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}
