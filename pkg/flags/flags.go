package flags

import "flag"
import "os"
import "strconv"

type Flags struct {
	Phost, Puser, Ppassword, Pdb, AlgoliaID, AlgoliaKey, ElasticHost, ElasticUser, ElasticPassword string
	Port, Pport, ElasticPort                                                                       int
	Wait, NoAlgolia                                                                                bool
}

func ParseFlags() Flags {
	var r Flags

	flag.StringVar(&r.Phost, "phost", "127.0.0.1", "Postgres host")
	flag.StringVar(&r.Puser, "puser", "postgres", "Postgres user")
	flag.StringVar(&r.Ppassword, "ppassword", "postgres", "Postgres password")
	flag.StringVar(&r.Pdb, "pdatabase", "podsearch", "Postgres database")
	flag.IntVar(&r.Pport, "pport", 5432, "Postgres port")
	flag.IntVar(&r.Port, "port", 8080, "Port on the API should listen on")
	flag.StringVar(&r.AlgoliaID, "algoliaid", "", "Algolia AppID")
	flag.StringVar(&r.AlgoliaKey, "algoliakey", "", "Algolia Admin Key")
	flag.StringVar(&r.ElasticHost, "elastichost", "127.0.0.1", "Elasticsearch host")
	flag.IntVar(&r.ElasticPort, "elasticport", 9200, "Elasticsearch port")
	flag.StringVar(&r.ElasticUser, "elasticuser", "elastic", "Elasticsearch user")
	flag.StringVar(&r.ElasticPassword, "elasticpassword", "changeme", "Elasticsearch pass")

	useEnv := flag.Bool("env", false, "Use Environment Variables instead of cmd line parameters; gets every other option from corresponding uppercase environment variables")
	flag.BoolVar(&r.Wait, "wait", false, "Do not exit if database or search are not available, wait instead")
	flag.BoolVar(&r.NoAlgolia, "noalgolia", false, "Do not use algolia (e.g. for testing)")

	flag.Parse()

	//if useEnv, grab variables from environment
	if *useEnv {
		f := ParseFromEnv()
		f.Wait = r.Wait
		f.NoAlgolia = r.NoAlgolia
		return f
	}
	return r
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
