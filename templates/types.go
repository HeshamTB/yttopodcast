package templates

import (
	"embed"
)

//go:embed *.templ
var TemplatesFS embed.FS

//go:embed base.rss.templ
var RSSTemplate string

type FeedData struct {
	Title               string
	PublishDateRfcEmail string
	BuildDateRfcEmail   string
	GeneratorName       string
	PodcastPage         string
	Lang                string
	CopyRight           string
	Summary             string
	PodcastImageURL     string
	FeedURL             string
	Items               []FeedItem
}

type FeedItem struct {
	Title               string
	PublishDateRfcEmail string
	Id                  string
	CoverImageURL       string
	Description         string
	Length              int
	EnclosureURL        string
	Duration            string
}
