package services

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"podscraper/models"
	"strconv"
	"time"
)

func ScrapeCharts() {
	//pod := &models.Podcast{}
	//rankingList := &[]models.Rank{}

	//err := models.GetDB().Table("ranks").Preload("Podcast").Order("score asc").Find(&rankingList).Error

	//if err != nil && err == gorm.ErrRecordNotFound {
	//	fmt.Println("Record not found")
	//} else {
	//	for _, rank := range *rankingList {
	//		fmt.Println("Record found: ", rank.Podcast.Name)
	//	}
	//}

	c := colly.NewCollector(
		colly.CacheDir("./itunescharts_cache"),
	)

	c.OnHTML("#chart li", func(element *colly.HTMLElement) {
		scoreStr := element.DOM.Find("span:nth-of-type(1)").Text()
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			panic(err)
		}

		itunesAttr, itunesErr := element.DOM.Find("p:nth-of-type(2) a").Attr("href")
		if itunesErr == false {
			panic(itunesErr)
		}

		itunesId := GetIDFromUrl(itunesAttr)

		podcast := models.GetPodcastByTrack(itunesId)

		if podcast == nil {
			podJSON := LookUp(itunesId)

			podcast = &models.Podcast{}

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

			podcast.Name = podJSON.CollectionName
			podcast.FeedUrl = podJSON.FeedUrl
			podcast.Genre = *genre
			podcast.ItunesId, _ = strconv.Atoi(itunesId)
			podcast.PublisherID = podJSON.ArtistID
			podcast.PublisherName = podJSON.ArtistName
			podcast.ImageUrlXL = podJSON.ArtworkUrlXL
			podcast.ImageUrlLG = podJSON.ArtworkUrlLG
			podcast.ImageUrlMD = podJSON.ArtworkUrlMD
			podcast.ImageURlSM = podJSON.ArtworkUrlSM
			podcast.EpisodesCount = podJSON.TrackCount

			err := podcast.Create()
			if err != nil {
				return
			}

			fmt.Println("Collection Name: " + podJSON.CollectionName)
		}

		fmt.Printf("%d. Podcast: %q\n,", score, podcast.Name)

		tempRank := models.GetRankByItunesId(itunesId)

		if tempRank != nil {
			if tempRank.Score != score {
				tempRank.Score = score
				tempRank.PodcastID = podcast.ID
				models.GetDB().Save(tempRank)
			}

		} else {
			itunesID, _ := strconv.Atoi(itunesId)

			rank := &models.Rank{
				Score:     score,
				PodcastID: podcast.ID,
				ItunesId:  itunesID,
			}

			models.GetDB().Create(rank).Related(&podcast)
		}

	})

	c.SetRequestTimeout(270 * time.Second)

	err := c.Visit("http://www.itunescharts.net/us/charts/podcasts/current")
	if err != nil {
		fmt.Println("Colly Error")
		fmt.Println(err.Error())
		return
	}
}
