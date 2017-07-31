package podcast

import (
	"bytes"
	"encoding/xml"
)

type enclosure struct {
	URL string `xml:"url,attr"`
}

type podcastXMLFeedItem struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"RssDefault link"`
	PubDate     string    `xml:"pubDate"`
	Description string    `xml:"description"`
	AudioURL    enclosure `xml:"enclosure"`
}

type podcastXMLFeed struct {
	Title       string               `xml:"channel>title"`
	Link        string               `xml:"RssDefault channel>link"`
	Language    string               `xml:"channel>language"`
	Description string               `xml:"channel>description"`
	ImageURL    string               `xml:"channel>image>url"`
	Items       []podcastXMLFeedItem `xml:"channel>item"`
}

func parseFromBody(content []byte) (*podcastXMLFeed, error) {
	var feed podcastXMLFeed
	d := xml.NewDecoder(bytes.NewReader(content))
	d.DefaultSpace = "RssDefault"

	err := d.Decode(&feed)
	if err != nil {
		return nil, err
	}
	return &feed, err
}
