package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"gitea.hbanafa.com/hesham/yttopodcast/feed"
)

var CHAN_ID *string = flag.String("id", "", "youtube channel id")

func main() {

    flag.Parse()

    if *CHAN_ID == "" {
        perr("provide channel id with -id <id>\n")
        os.Exit(1)
    }

    url := fmt.Sprintf(feed.YT_FEED_URL, *CHAN_ID)
    perr(url + "\n")
    resp, err := http.Get(url)

    if err != nil {
        perr(err.Error() + "\n")
        os.Exit(1)
    }

    if resp.StatusCode != http.StatusOK {
        perr("http: endpoint returned %s\n", resp.Status)
        os.Exit(1)
    }

    sc := bufio.NewReader(resp.Body)
    for {
        newl, err := sc.ReadString('\n')
        if err != nil {

            if errors.Is(err, io.EOF) {
                break
            }
            perr(err.Error() + "\n")
            os.Exit(1)
        }
        fmt.Print(newl)
    }
}

func perr(f string, args... any) {
    fmt.Fprintf(os.Stderr, f, args...)
}
