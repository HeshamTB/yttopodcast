package bouncer

import (
	"context"
	"net/http"

	"gitea.hbanafa.com/hesham/yttopodcast/ytlinkprov"
)

func UrlCache(next http.Handler, url_prov ytlinkprov.YtLinkProvider) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rr := r.WithContext(context.WithValue(r.Context(), CTX_LINKPROV, url_prov)) 
        next.ServeHTTP(w, rr)
    })
}
