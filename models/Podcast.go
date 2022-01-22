package models

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type Podcast struct {
	GormModel
	Name          string `json:"name"`
	Description   string `sql:"type:longtext"`
	ItunesId      int    `json:"itunesId"`
	PublisherName string `json:"publisherName"`
	PublisherID   int    `json:"publisherId"`
	ImageURlSM    string `json:"imageUrlSM"`
	ImageUrlMD    string `json:"imageUrlMD"`
	ImageUrlLG    string `json:"imageUrlLG"`
	ImageUrlXL    string `json:"imageUrlXL"`
	Genre         Genre  `json:"genre"`
	GenreID       string `json:"-"`
	FeedUrl       string `json:"feedUrl"`
	EpisodesCount int    `json:"episodesCount"`
}

type Rank struct {
	GormModel
	Score     int
	Podcast   Podcast `gorm:"foreignkey:PodcastID"`
	PodcastID string
	ItunesId  int `gorm:"unique" json:"itunesId"`
}

func GetRankByItunesId(itunesId string) *Rank {
	rank := &Rank{}
	trackId, _ := strconv.Atoi(itunesId)

	err := GetDB().Table("ranks").Where("itunes_id = ?", trackId).First(rank).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil
	}

	return rank
}

func (podcast *Podcast) Create() error {
	err := GetDB().Create(&podcast).Error
	return err
}

func GetPodcastByTrack(itunesId string) *Podcast {
	podcast := &Podcast{}
	trackId, _ := strconv.Atoi(itunesId)

	err := GetDB().Table("podcasts").Preload("Genre").Where("itunes_id = ?", trackId).First(podcast).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil
	}

	return podcast
}
