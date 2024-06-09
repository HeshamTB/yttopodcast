package bouncer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gitea.hbanafa.com/hesham/yttopodcast/ytlinkprov"
	"github.com/lrstanley/go-ytdlp"
)

const CTX_LINKPROV = "linkprov"

type Bouncer struct {
    http.Server
    ytdlpInstall ytdlp.ResolvedInstall
    urlProvider ytlinkprov.YtLinkProvider
}

func NewBouncerHTTPServer(
    ctx context.Context, 
    listAddr string, 
    link_prov ytlinkprov.YtLinkProvider,
) (srv *Bouncer, err error) {

    ytInstall, err := ytdlp.Install(
        ctx,
        &ytdlp.InstallOptions{
            AllowVersionMismatch: true,
        },
    )
    if err != nil {
        return nil, err
    }

    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", handleGETBounce)

    var httpHandler http.Handler = mux
    httpHandler = UrlCache(mux, link_prov)

    return &Bouncer{
        urlProvider: link_prov,
        Server: http.Server{
            WriteTimeout: time.Second * 60,
            ReadTimeout:  time.Second * 60,
            Addr:         listAddr,
            Handler:      httpHandler,
        },
        ytdlpInstall: *ytInstall,
    }, nil
}

func handleGETBounce(w http.ResponseWriter, r *http.Request) {
    urlProv, ok := r.Context().Value(CTX_LINKPROV).(ytlinkprov.YtLinkProvider)
    if !ok {
        fmt.Fprintf(os.Stderr, "Could not get url provider from ctx!\n")
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    id := r.URL.Query().Get("id")
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    log.Printf("request for %s", id)
    // vidUrl := fmt.Sprintf("https://youtube.com/watch?v=%s", id)
    // ytCmd := ytdlp.New().ExtractAudio().GetURL()
    // ytRes, err := ytCmd.Run(r.Context(), vidUrl)
    link, err := urlProv.GetLink(id)
    if err != nil {
        _, ok := err.(*ytdlp.ErrExitCode)
        if ok {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        fmt.Fprintln(os.Stderr, err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "audio/mp3")
    http.Redirect(w, r, strings.Trim(link, "\n"), http.StatusFound)
}
