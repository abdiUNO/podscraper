package services

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"log"
	"net/url"
	"podscraper/utils"
	"strings"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

type ItunesPodcast struct {
	CollectionName string   `json:"collectionName"`
	CollectionId   int      `json:"collectionId"`
	ArtistName     string   `json:"artistName"`
	ArtistID       int      `json:"artistId"`
	ArtworkUrlSM   string   `json:"artworkUrl30"`
	ArtworkUrlMD   string   `json:"artworkUrl60"`
	ArtworkUrlLG   string   `json:"artworkUrl100"`
	ArtworkUrlXL   string   `json:"artworkUrl600"`
	Genre          string   `json:"primaryGenreName"`
	Genres         []string `json:"genres"`
	FeedUrl        string   `json:"feedUrl"`
	TrackCount     int      `json:"trackCount"`
}

type LookUpResponse struct {
	ResultCount int             `json:"resultCount"`
	Results     []ItunesPodcast `json:"results"`
}

func LookUp(id string) *ItunesPodcast {
	podcasts := LookUpResponse{}

	res, ok := utils.HttpGet("https://itunes.apple.com/lookup?id=" + id)

	if ok != nil {
		log.Fatal(ok)
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	if err := json.Unmarshal(body, &podcasts); err != nil {
		log.Fatal(err)
	}

	pod := &podcasts.Results[0]

	return pod
}

func ParseFeed(feedUrl string) *gofeed.Feed {

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedUrl)

	out, err := json.Marshal(feed)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = client.Set(feedUrl, out, 0).Err()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return feed
}

func GetIDFromUrl(urlString string) string {
	parsedUrl, _ := url.Parse(urlString)
	urlPart := strings.Split(parsedUrl.Path, "/")
	id := strings.TrimLeft(urlPart[4], "id")

	return id
}
