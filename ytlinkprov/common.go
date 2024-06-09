package ytlinkprov

type YtLinkProvider interface {
    GetLink(id string) (link string, err error)
}

