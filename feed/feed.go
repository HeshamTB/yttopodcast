package feed

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"text/template"
	"time"

	"gitea.hbanafa.com/hesham/yttopodcast/templates"
	"github.com/mmcdole/gofeed"
)


const (
    YT_FEED_URL  = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"
    YT_VIDEO_URL = "https://youtube.com/watch?v=%s"
    __GENERATOR_NAME = "yttopodcast - H.B."
)

type RSSMetadata struct {
    Summary string    
    Languge string
    Copyright string
    BounceURL string
}

/* bounce_url in the format of http://domain/?id=%s */
func ConvertYtToRss(w io.Writer, channel_id string, meta RSSMetadata) error {

    channelUrl := fmt.Sprintf(YT_FEED_URL, channel_id)
    feed, err := getFeed(channelUrl)
    if err != nil {
	return feedErr(err)
    }
    return convertFeedToRSS(w, *feed, meta)
}

// Convert to Yt Atom to RSS given a Reader that provides xml
func ConvertAtomToRSS(w io.Writer, r io.Reader, meta RSSMetadata) error {

    feed, err := gofeed.NewParser().Parse(r)
    if err != nil {
	return feedErr(err)
    }
    return convertFeedToRSS(w, *feed, meta)
}

func convertFeedToRSS(w io.Writer, feed gofeed.Feed, meta RSSMetadata) error {

    var podFeed templates.FeedData

    t_now := time.Now().UTC()
    podFeed.Title = feed.Title
    podFeed.Summary = meta.Summary
    podFeed.BuildDateRfcEmail = t_now.Format(time.RFC1123Z)
    podFeed.CopyRight = meta.Copyright
    podFeed.PublishDateRfcEmail = t_now.Format(time.RFC1123Z)
    podFeed.PodcastPage = feed.Link
    podFeed.Lang = meta.Languge
    podFeed.GeneratorName = __GENERATOR_NAME

    for i, item := range feed.Items {

	subStrings := strings.Split(item.GUID, ":")
	id := subStrings[2]

	bounceURL, err := url.Parse(
	    fmt.Sprintf(meta.BounceURL, id),
	)
	if err != nil {
	    return err
	}

	// Check this out
	g := item.Extensions["media"]["group"]
	gg := g[len(g)-1]
	thumb := gg.Children["thumbnail"][0]

	de := gg.Children["description"]
	desc := de[0].Value

	coverArtUrl, err := url.Parse(thumb.Attrs["url"])
	if err != nil {
	    return errors.Join(err, errors.New(
		fmt.Sprintf(
		    "could not parse item cover art for %s GUID: %s\n",
		    item.Title,
		    item.GUID,
	    )))
	}

	if i == 0 {
	    podFeed.PodcastImageURL = coverArtUrl.String()
	}

	podFeed.Items = append(podFeed.Items,
	    templates.FeedItem{
		Title:               item.Title,
		CoverImageURL:       coverArtUrl.String(),
		Id:                  id,
		Duration:            "0",
		PublishDateRfcEmail: item.PublishedParsed.Format(time.RFC1123Z),
		Description:         desc,
		Length:              0,
		EnclosureURL:	     bounceURL.String(),
	})
    }

    rssTemplate, err := template.New("rss").Parse(templates.RSSTemplate)
    if err != nil {
	return err
    }

    err = rssTemplate.Execute(w, podFeed)
    if err != nil {
	return err
    }

    rssResult := bytes.Buffer{}
    rssTemplate.Execute(&rssResult, podFeed)
    _, err = gofeed.NewParser().ParseString(rssResult.String())
    if err != nil {
	return feedErr(err)
    }

    return nil
}

func feedErr(err error) error {
    httpErr, ok := err.(gofeed.HTTPError)
    if ok {
	switch httpErr.StatusCode {
	case 404:
	    return errors.Join(err, errors.New("yt: could not find channel id"))
	}
    }
    return err
}

func getFeed(url string) (*gofeed.Feed, error) {
    parser := gofeed.NewParser()
    return parser.ParseURL(url)
}


