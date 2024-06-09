package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"gitea.hbanafa.com/hesham/yttopodcast/bouncer"
	"gitea.hbanafa.com/hesham/yttopodcast/ytlinkprov"
)

var listenAddr = flag.String("listen-addr", ":8081", "Address and port to listen on")

func main() {

	ctx := context.Background()
	flag.Parse()
	linkProv, err := ytlinkprov.NewCachedLinkProvider(time.Minute * 30)

	bouncer, err := bouncer.NewBouncerHTTPServer(ctx, *listenAddr, linkProv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	fmt.Printf("Starting server on %s\n", *listenAddr)
	err = bouncer.ListenAndServe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
