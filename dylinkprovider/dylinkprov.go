package dylinkprovider

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "sync"
    "time"

    "gitea.hbanafa.com/hesham/yttopodcast/ytlinkprov"
    "github.com/lrstanley/go-ytdlp"
)

const (
    Q_EXPIRE = "expire"
    MAX_MAP_SZ = 1000
)

type DynamicCacheExpLinkProvider struct {
    cache map[string]url.URL
    l *log.Logger
    lock *sync.RWMutex
}

var _ = (ytlinkprov.YtLinkProvider)((*DynamicCacheExpLinkProvider)(nil))

func NewDynCacheExpLinkProv(l *log.Logger) *DynamicCacheExpLinkProvider {
    p := new(DynamicCacheExpLinkProvider)
    p.l = l
    p.lock = &sync.RWMutex{}
    p.cache = make(map[string]url.URL)
    ytdlp.MustInstall(context.Background(), &ytdlp.InstallOptions{})
    return p
}

// GetLink implements ytlinkprov.YtLinkProvider.
func (d *DynamicCacheExpLinkProvider) GetLink(id string) (link string, err error) {


    d.lock.RLock()
    cl1, ok := d.cache[id]
    d.lock.RUnlock()
    if ok && !isExpired(cl1) && is200(cl1) {
        d.l.Printf("[cache] hit on %s\n", id)
        return cl1.String(), nil              
    }

    d.l.Printf("[cache] miss on %s\n", id)
    newlink, err := getRemoteLink(id)
    if err != nil {
        return "", err
    }

    newlinkurl, err := url.Parse(newlink)
    if err != nil {
        return "", err
    }

    d.lock.Lock()
    if len(d.cache) >= MAX_MAP_SZ {
        var k string
        for kk := range d.cache { 
            k = kk
            break
        }
        delete(d.cache, k)
    }
    d.cache[id] = *newlinkurl
    d.lock.Unlock()

    d.l.Printf("[cache] new entry for %s\n", id)
    return newlinkurl.String(), nil

}

func getRemoteLink(id string) (string, error) {

    var link string

    vidUrl := fmt.Sprintf("https://youtube.com/watch?v=%s", id)
    ytCmd := ytdlp.New().ExtractAudio().GetURL()
    ytRes, err := ytCmd.Run(context.Background(), vidUrl)
    if err != nil {
        return "", err
    }
    linkFirst := strings.Split(ytRes.Stdout, "\n")[0]

    /* Get the last link in a chain of 3XX codes*/
    resp, err := http.Get(linkFirst)
    if err != nil {
        return "", err
    }

    if resp.StatusCode != http.StatusOK {
        return linkFirst, nil
    }

    link = resp.Request.URL.String()
    return link, nil

}

func isExpired(link url.URL) bool {

    exp := link.Query().Get("expire")
    if exp == "" {
        return true
    }

    tunixd, err := strconv.ParseInt(exp, 10, 64)
    if err != nil {
        return true
    }

    tunix := time.Unix(tunixd, 0)
    // If (tunix - now) is negative, the link is expired
    delta := tunix.Sub(time.Now())
    if delta <= 0 {
        return true
    }

    // Still not expired but we check delta against duration
    durd := 0.0

    dur := link.Query().Get("dur")
    if dur == "" {
        return true
    }

    durd, err = strconv.ParseFloat(exp, 64)
    if err != nil {
        return true
    }

    if delta < time.Duration(durd) + time.Second * 2400 {
        return true
    }

    return false
}

func is200(link url.URL) bool {
    resp, _ := http.Get(link.String())
    return resp.StatusCode == 200
}

