# yttopodcast
Tools to convert YouTube video feeds to standard Podcast RSS feeds. A workaround is done, bouncer,
to serve files from youtube directly, rather than storing the audio files and serving them
independantly.

A web server is required to serve the RSS feed as a file, and an HTTP boucner server is used to
fetch valid links from yt. Since getting the content URLs is quite slow, the implemented bouncer
has a basic cache.

# Usage
The tools can be used standalone, imported and invoked in go code, or run with the helpers as a collection. 

## yttopodcast
Generates single or multiple feeds on an interval, serves the feed files (RSS/XML), and launches a yt 
link bouncer. This can be used as a complete service.
```sh
cd cmd/yttopodcast
go build .
./yttopodcast -id CHANNEL_ID
```

## genfeed
Generate a feed given a channel id
```sh
cd cmd/genfeed
go build .
./genfeed -id CHANNEL_ID > feed.xml
```
`feed.xml` file can be used as an RSS feed.

## ytbouncer
The bouncer uses standard go http and can be embedded. To run it standalone
```sh
cd cmd/ytbouncer
go build .
./ytbouncer
```
starts a server.

**Note:** the url resolution is quite slow and heavy computaionally, due to yt-dlp backend. It can be used
as a DoS without the use of rate limiting or other measures.


