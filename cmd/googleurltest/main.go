package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println(os.Args[0], " <url>")
        os.Exit(1)
    }

    url, err := url.Parse(os.Args[1])
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Printf("url.Query(): %+v\n", url.Query())

    resp, err := http.Get(url.String())
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if resp.StatusCode != 200 {
        fmt.Printf("http: got %d\n", resp.StatusCode)
    }

    if url != resp.Request.URL {
        fmt.Printf("resp.Request.URL: %+v\n", resp.Request.URL)
        fmt.Printf("resp.Request.URL.Query(): %+v\n", resp.Request.URL.Query())
    }

    time_s := url.Query().Get("expire")
    if time_s == "" {
        fmt.Println("url: expire key missing")
        os.Exit(1)
    }

    time_ss, err := strconv.ParseInt(time_s, 10, 64)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    unix := time.Unix(time_ss, 0)
    duration := unix.Sub(time.Now())
    fmt.Printf("duration.String(): %v\n", duration.String())

}
