package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:    false,
		DisableTimestamp: false,
		PadLevelText:     true,
	})

	//services.ScrapePods()
	//services.ScrapeCharts()

	/*
		//c.OnHTML("section.l-row .l-column .product-hero-desc .product-hero-desc__section p", func(element *colly.HTMLElement) {
		//	text := element.Text
		//	urlPart := strings.Split(element.Request.URL.Path, "/")
		//	itunesId := strings.TrimLeft(urlPart[4], "id")
		//
		//	podJSON := services.LookUp(itunesId)
		//
		//	fmt.Printf("------------- \n")
		//
		//	fmt.Printf("Name:  %q\n", podJSON.CollectionName)
		//	fmt.Printf("Description:  %q\n", text[0:75])
		//	fmt.Printf("Genre:  %q\n", podJSON.Genre)
		//
		//	fmt.Printf("------------- \n")
		//
		//	temp := models.GetPodcastByTrack(itunesId)
		//	if temp != nil && len(temp.Description) == 0 {
		//		fmt.Println("Updating Record")
		//		temp.Description = text
		//		models.GetDB().Save(temp)
		//	}
		//})

		////Visit Podcast link
		//c.OnResponse(func(r *colly.Response) {
		//	if strings.HasPrefix(r.Request.URL.String(), "https://podcasts.apple.com/us/podcast/") {
		//		fmt.Println("Visited", r.Request.URL)
		//	}
		//})

		// Start scraping on itunes.apple.com
		//err := c.Visit("https://itunes.apple.com/us/genre/podcasts/id26?mt=2")
		//if err != nil {
		//	return
		//}
	*/
}
