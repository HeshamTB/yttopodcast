package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gitea.hbanafa.com/hesham/yttopodcast/bouncer"
	"gitea.hbanafa.com/hesham/yttopodcast/feed"
	"gitea.hbanafa.com/hesham/yttopodcast/ytlinkprov"
)

var (
    l log.Logger

    listenAddr  = flag.String("addr", "0.0.0.0:8086", "bouncer listen address")
    fileServAdd  = flag.String("fs-addr", "0.0.0.0:8087", "http file server listen address")
    interval = flag.Int("interval", 30, "update interval for feed in minutes")
    chan_id = flag.String("id", "", "YouTube channel ID")
    bounc_url = flag.String("bouncer", "http://localhost:8081/?id=%s", "bouncer url as format string")
    lang = flag.String("lang", "en", "content language")
)


func main() {
    
    flag.Parse()
    l = *log.Default()

    if *chan_id == "" {
        l.Println("no channel id provided")
        os.Exit(1)
    }
    
    cache, err := ytlinkprov.NewCachedLinkProvider(time.Minute * time.Duration(*interval))
    if err != nil {
        l.Println(err.Error())
        os.Exit(1)
    }

    bouncer, err := bouncer.NewBouncerHTTPServer(context.Background(), *listenAddr, cache)
    if err != nil {
        l.Println(err.Error())
        os.Exit(1)
    }

    // This goro is sleeping until interval or SIGINT
    // Another one is running the bouncer

    go func() {
        l.Printf("http bouncer server starting on %s\n", bouncer.Addr)
        l.Println(bouncer.ListenAndServe())
    }()

    os.Mkdir("feeds", 0700)

    mux := http.NewServeMux()
    mux.Handle("GET /", http.FileServer(http.Dir("feeds")))
    fileServer := http.Server{
        Addr: *fileServAdd,
        Handler: mux,
        ErrorLog: &l,
        ReadTimeout: time.Second * 20,
        WriteTimeout: time.Second * 20,
    }

    go func() {
        l.Printf("http file server starting on %s\n", fileServer.Addr)
        l.Println(fileServer.ListenAndServe())
    }()

    err = genFeed()
    if err != nil {
        l.Println(err.Error())
        os.Exit(1)
    }
    
    sig := make(chan os.Signal)

    signal.Notify(sig, os.Interrupt, os.Kill)

    l:
    for {
        select {
        case s := <-sig:
            l.Println("got ", s.String())
            fileServer.Shutdown(context.Background())
            bouncer.Shutdown(context.Background())
            break l
        case <-time.NewTicker(time.Minute * time.Duration(*interval)).C:
            l.Println("tick")
            genFeed()
    }}

}

func genFeed() error {

    l.Println("generating feed")
    file, err := os.Create("./feeds/f.xml")
    if err != nil {
        return err
    }
    defer file.Close()

    return feed.ConvertYtToRss(file, *chan_id, feed.RSSMetadata{
        Languge: *lang,
        BounceURL: *bounc_url,
    })
}
