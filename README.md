# Podsearch
[![CircleCI](https://circleci.com/gh/sauercrowd/podsearch.svg?style=svg)](https://circleci.com/gh/sauercrowd/podsearch)

Podsearch is a backend to make podcast descriptions searchable and probably more in the future.

Podsearch uses algolia and/or elasticsearch for searching and stores podcast information also in a postgres database


# Usage

Recommended: docker-compose

```
docker-compose up
```

# Command line options
```
▶ podsearch --help
Usage of podsearch:
  -algoliaid string
        Algolia AppID
  -algoliakey string
        Algolia Admin Key
  -elastichost string
        Elasticsearch host (default "127.0.0.1")
  -elasticpassword string
        Elasticsearch pass (default "changeme")
  -elasticport int
        Elasticsearch port (default 9200)
  -elasticuser string
        Elasticsearch user (default "elastic")
  -env
        Use Environment Variables instead of cmd line parameters; gets every other option from corresponding uppercase environment variables
  -noalgolia
        Do not use algolia
  -pdatabase string
        Postgres database (default "podsearch")
  -phost string
        Postgres host (default "127.0.0.1")
  -port int
        Port on the API should listen on (default 8080)
  -ppassword string
        Postgres password (default "postgres")
  -pport int
        Postgres port (default 5432)
  -puser string
        Postgres user (default "postgres")
  -wait
        Do not exit if database or search are not available, wait instead
```

# API
|Method|URI|
|------|----------|
|POST|[/api/v1/podcast](#/api/v1/podcast)|
|GET|[/api/v1/podcast](#/api/v1/podcast)|
|GET|[/api/v1/search](#/api/v1/search)|

## /api/v1/podcast
### GET
Returns saved informations to a podcast

Query: url=...
Example 

```
GET /api/v1/podcast?url=http://feeds.feedburner.com/pod-save-america?format=xml
```

### POST
Add a new podcast 

Body:
```
{
    "url":"..."
}
```

Example:

POST /api/v1/podcast
Body:
```
{
    "url":"http://feeds.feedburner.com/pod-save-america?format=xml"
}
```

## /api/v1/search
GET /api/v1/search

Query: q=...

Example:
```
GET /api/v1/search?q=politics
```