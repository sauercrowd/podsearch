# Routes
|Route|Method|Purpose  |
|-----|------|---------|
|/api/v1/addpodcast|POST|Add a new podcast to the system|
|/api/v1/getpodcast|GET|Get all podcasts|
|/api/v1/getpodcast?url=x|GET|Get a podcast by its url|

### /api/v1/addpodcast
POST Body:
```
{
    "url": "http://mypodcast.com/feed.rss"
}