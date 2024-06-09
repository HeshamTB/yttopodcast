# yttopodcast
Tools to convert YouTube video feeds to standard Podcast RSS feeds. A workaround is done, bouncer,
to serve files from youtube directly, rather than storing the audio files and serving them
independantly.

A web server is required to serve the RSS feed as a file, and an HTTP boucner server is used to
fetch valid links from yt. Since getting the content URLs is quite slow, the implemented bouncer
has a basic cache.

## genfeed
Generate a feed given a channel id
```sh
cd cmd/genfeed
go build .
./genfeed -id CHANNEL_ID > feed.xml
```

## ytbouncer
The bouncer uses standard go http and can be embedded
```sh
cd cmd/ytbouncer
go build .
./ytbouncer
```
starts a server.

**Note:** the url resolution is quite slow and heavy computaionally, due to yt-dlp backend. It can be used
as a DoS without the use of rate limiting or other measures.


