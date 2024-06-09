package ytlinkprov

import (
	"context"
	"fmt"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

type CacheLinkProv struct {
    cache map[string]TimedLink
    cacheWindow time.Duration
    ytInstall ytdlp.ResolvedInstall
}

type TimedLink struct {
    Link string
    Time time.Time
}

func NewCachedLinkProvider(expiration time.Duration) (*CacheLinkProv, error) {
    p := new(CacheLinkProv)
    p.cache = make(map[string]TimedLink)
    p.cacheWindow = expiration
    ctx := context.Background()
    ytInstall, err := ytdlp.Install(
        ctx,
        &ytdlp.InstallOptions{
            AllowVersionMismatch: true,
        },
    )
    if err != nil {
        return p, err
    }
    p.ytInstall = *ytInstall

    return p, nil
}

func (c *CacheLinkProv) GetLink(id string) (link string, err error) {

    cc, ok := c.cache[id]
    if ok && c.validCache(cc) {
        return cc.Link, nil
    }
    link, err = getRemoteLink(id)
    if err != nil {
        return "", err
    }
    t_now := time.Now().UTC()
    c.cache[id] = TimedLink{
        Link: link,
        Time: t_now,
    }

    // INFO: This is a vary, vary slow leak

    return link, nil
}

func (c *CacheLinkProv) validCache(l TimedLink) bool {
    t_exp := time.Now().UTC().Add(c.cacheWindow)
    if t_exp.Before(l.Time) {
        // expired
        return false
    }
    return true
    
}

func getRemoteLink(id string) (string, error) {
    vidUrl := fmt.Sprintf("https://youtube.com/watch?v=%s", id)
    ytCmd := ytdlp.New().ExtractAudio().GetURL()
    ytRes, err := ytCmd.Run(context.Background(), vidUrl)
    if err != nil {
        return "", err
    }
    return ytRes.Stdout, nil
}
