package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gitea.hbanafa.com/hesham/yttopodcast/feed"
)

const (
    EXIT_ERR_BAD_CLI = 64
)

var (
    chan_id = flag.String("id", "", "YouTube channel ID")
    bounc_url = flag.String("bouncer", "http://localhost:8081/?id=%s", "Bouncer url as format string")
    lang = flag.String("lang", "en", "Content Language")
)

func main() {
    flag.Parse()
    if err := validFlags(); err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(EXIT_ERR_BAD_CLI)
    }
    fmt.Fprintf(os.Stderr, "id: %s\nbouncer: %s\n", *chan_id, *bounc_url)
    
    err := feed.ConvertYtToRss(os.Stdout, *chan_id, *bounc_url, 
        feed.RSSMetadata{Languge: "en", Copyright: "N/A", Summary: "YouTube Channel as podcast"})
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }
}

func validFlags() error {
    if *chan_id == "" {
        return errors.New("flag: id flag missing")
    }
    return nil
}
