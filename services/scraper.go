package services

import (
	"fmt"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"podscraper/models"
	"strconv"
	"strings"
)

const PodcastGenres string = "https://podcasts.apple.com/us/genre/podcasts/id26"

var categories = map[string]bool{"Politics": true, "Philosophy": true}

var (
	genreQuery       = "#genre-nav .grid3-column ul.list.column li ul.list.top-level-subgenres li"
	podcastQuery     = "#selectedcontent .column ul li"
	descriptionQuery = "section .l-row .l-column:first-of-type .product-hero-desc .product-hero-desc__section p"
)

func CrawlGenre(e *colly.HTMLElement) {
	genreUrl := e.ChildAttr("a", "href")
	genreName := e.ChildText("a")
	// Visit link found on page
	log.WithFields(log.Fields{
		"Genre": genreName,
	}).Info("Visited Link")

	genre := models.GetGenreByName(genreName)

	if genre == nil {
		genre = &models.Genre{
			Name: genreName,
			Url:  genreUrl,
		}

		err := genre.Create()
		if err != nil {
			log.Fatal(err)
		}
	}

	value, exists := categories[genreName]

	if exists && value == true {
		err := e.Request.Visit(e.Request.AbsoluteURL(genreUrl))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CrawlPodcast(element *colly.HTMLElement) {
	text := element.Text
	urlPart := strings.Split(element.Request.URL.Path, "/")
	itunesId := strings.TrimLeft(urlPart[4], "id")

	temp := models.GetPodcastByTrack(itunesId)

	log.Printf("------------- \n")

	if temp != nil {
		substr := len(text) / 4
		log.WithFields(log.Fields{
			"Name":        temp.Name,
			"Description": text[0:substr],
			"Genre":       temp.Genre.Name,
		}).Info("Pod Record Found")

		if len(temp.Description) == 0 {
			//fmt.Println("Updating Record")
			temp.Description = text

			models.GetDB().Save(temp)
		}
	} else {
		podJSON := LookUp(itunesId)
		substr := len(text) / 4

		log.WithFields(log.Fields{
			"Name":        podJSON.CollectionName,
			"Description": text[0:substr],
			"Genre":       podJSON.Genre,
		}).Info("Saving Pod Record")

		itunesIdParsed, _ := strconv.Atoi(itunesId)

		genre := models.GetGenreByName(podJSON.Genre)

		if genre == nil {
			genre = &models.Genre{
				Name: podJSON.Genre,
				Url:  "",
			}

			err := genre.Create()
			if err != nil {
				log.Fatal(err)
			}
		}

		pod := &models.Podcast{
			Name:          podJSON.CollectionName,
			Description:   text,
			FeedUrl:       podJSON.FeedUrl,
			Genre:         *genre,
			ItunesId:      itunesIdParsed,
			PublisherID:   podJSON.ArtistID,
			PublisherName: podJSON.ArtistName,
			ImageUrlXL:    podJSON.ArtworkUrlXL,
			ImageUrlLG:    podJSON.ArtworkUrlLG,
			ImageUrlMD:    podJSON.ArtworkUrlMD,
			ImageURlSM:    podJSON.ArtworkUrlSM,
			EpisodesCount: podJSON.TrackCount,
		}

		err := pod.Create()
		if err != nil {
			return
		}
	}
	log.Printf("------------- \n")
}

func ScrapePods() {
	var initUrl = PodcastGenres

	// Instantiate default collector
	c := colly.NewCollector(
		colly.MaxDepth(4),
		colly.CacheDir("./itunes_cache"),
	)

	detailCrawler := c.Clone()
	//Visit Genre
	c.OnHTML(genreQuery, CrawlGenre)
	//c.OnHTML("#genre-nav.main.nav:nth-child(2)  div.grid3-column  ul.list.column.last:nth-child(3) > li:nth-child(5) ", CrawlGenre)

	//// On every a element which has href attribute call callback
	c.OnHTML(podcastQuery, func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		//fmt.Printf("Podcast found: %q\n", e.Text)
		log.WithFields(log.Fields{
			"Pod Url": link,
		}).Info("Visited Link")

		err := detailCrawler.Visit(e.Request.AbsoluteURL(link))

		if err != nil {
			log.Fatal(err)
		}
	})

	detailCrawler.OnHTML(descriptionQuery, CrawlPodcast)

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(initUrl)
	if err != nil {
		log.Fatal(err)
	}

}
