package main

import (
    "context"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "time"

    "gitea.hbanafa.com/hesham/yttopodcast/bouncer"
    "gitea.hbanafa.com/hesham/yttopodcast/dylinkprovider"
    "gitea.hbanafa.com/hesham/yttopodcast/feed"
)

var (
    l log.Logger

    listenAddr  = flag.String("addr", "0.0.0.0:8086", "bouncer listen address")
    fileServAdd  = flag.String("fs-addr", "0.0.0.0:8087", "http file server listen address")
    interval = flag.Int("interval", 30, "update interval for feed in minutes")
    chan_id = flag.String("id", "", "YouTube channel ID")
    chanlist_file = flag.String("list-file", "", "file with newline seperated channel ids")
    bounc_url = flag.String("bouncer", "http://localhost:8081/?id=%s", "bouncer url as format string")
    lang = flag.String("lang", "en", "content language")
)


func main() {
    
    flag.Parse()
    l = *log.Default()


    var ids []string
    if *chan_id != "" {
        l.Println("adding id arg")
        ids = append(ids, *chan_id)
    }

    if *chanlist_file != "" {
        listids, err := func() ([]string, error) {
            f, err := os.Open(*chanlist_file)
            if err != nil {
                return nil, err
            }

            content, err := io.ReadAll(f)
            if err != nil {
                return nil, err
            }

            ids := strings.Split(string(content), "\n")
            for i := range ids {
                ids[i] = strings.ReplaceAll(strings.Join(strings.Fields(ids[i]), ""), "\r", "")
                ids[i] = strings.ReplaceAll(ids[i], "\t", "")
            }

            return ids[:len(ids)-1], nil
        }()         
        if err != nil {
            l.Println(err.Error())
            os.Exit(1)
        }
        ids = append(ids, listids...)
    }

    if len(ids) == 0 {
        l.Println("no channel id or list file provided")
        os.Exit(1)
    }
    l.Println("channel count: ", len(ids))
    l.Printf("channels: %v\n", ids)

    l.Println("[feed] initial feed generation")
    genFeeds(ids)

    cache := dylinkprovider.NewDynCacheExpLinkProv(&l)

    bouncer, err := bouncer.NewBouncerHTTPServer(context.Background(), *listenAddr, cache)
    if err != nil {
        l.Println(err.Error())
        os.Exit(1)
    }

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
        l.Printf("http bouncer server starting on %s\n", bouncer.Addr)
        l.Println(bouncer.ListenAndServe())
        l.Println("http bouncer stopped")
    }()

    go func() {
        l.Printf("http file server starting on %s\n", fileServer.Addr)
        l.Println(fileServer.ListenAndServe())
        l.Println("http file server stopped")
    }()

    
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
            genFeeds(ids)
    }}

}

func genFeeds(ids []string) {
    for _, id := range ids {
        l.Printf("[feed] generating feed for %s\n", id)
        err := genFeed(id)
        if err != nil {
            l.Printf(err.Error())
            continue
        }
    }
}

func genFeed(id string) error {

    tmpFilename := fmt.Sprintf("./feeds/%s.xml.t", id)
    finalFilename := fmt.Sprintf("./feeds/%s.xml", id)

    file, err := os.Create(tmpFilename)
    if err != nil {
        return err
    }
    defer file.Close()
    defer os.Remove(tmpFilename)

    err = feed.ConvertYtToRss(file, id, feed.RSSMetadata{
        Languge: *lang,
        BounceURL: *bounc_url,
    })

    if err != nil { return err }

    return os.Rename(tmpFilename, finalFilename)
}
