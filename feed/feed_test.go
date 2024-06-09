package feed

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestParseYTAtom(t *testing.T) {
    fnames := []string{"a.xml", "b.xml", "c.xml"}
    for _, n := range fnames {
	nn := "testdata/" + n
	runTestYTAtom(t, nn)
    }
}

func runTestYTAtom(t *testing.T, fname string) {
    t.Run(fname, func(t *testing.T) {

	f, err := os.Open(fname)
	if err != nil {
	    t.Error(err.Error())
	}

	feed, err := io.ReadAll(f)
	if err != nil {
	    t.Error(err.Error())
	}

	parser := gofeed.NewParser()
	parser.ParseString(string(feed))

    })
}

func TestFetchRemoteYTFeed(t *testing.T) {
    f, err := os.Open("testdata/feedlink.txt")
    if err != nil {
	t.Log(err.Error())
	t.FailNow()
    }

    linkb, err := io.ReadAll(f)
    if err != nil {
	t.Log(err.Error())
	t.FailNow()
    }

    link := strings.Trim(string(linkb), "\n")
    
    _, err = url.Parse(link)
    if err != nil {
	t.Log(err.Error())
	t.FailNow()
    }

    resp, err := http.Get(link)
    if err != nil {
	t.Log(err.Error())
	t.FailNow()
    }
    defer resp.Body.Close()

    _, err = gofeed.NewParser().Parse(resp.Body)
    if err != nil {
	t.Log(err.Error())
	t.FailNow()
    }

}

func TestAtomToRSS(t *testing.T) {
    f, err := os.Open("testdata/a.xml")
    if err != nil {
	t.Error(err.Error())
	t.FailNow()
    }

    buf := bytes.Buffer{}
    
    err = convertAtomToRSS(&buf, f, RSSMetadata{BounceURL: "http://localhost:8081/q=%s"})
    if err != nil {
	t.Error(err.Error())
	t.FailNow()
    }

    _, err = gofeed.NewParser().Parse(&buf)
    if err != nil {
	t.Error(err.Error())
	t.FailNow()
    }

}

