package gobgg

import (
	"net/http"
	"net/url"
)

// BGG is the client for the boardgame geek site
type BGG struct {
	host   string
	scheme string
	client *http.Client
}

func (bgg BGG) buildURL(path string, args map[string]string) string {
	u := url.URL{
		Scheme: bgg.scheme,
		Host:   bgg.host,
		Path:   path,
	}

	q := u.Query()
	for i := range args {
		q.Set(i, args[i])
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// NewBGGClient returns a new client
func NewBGGClient() *BGG {
	return &BGG{
		host:   "boardgamegeek.com",
		scheme: "https",
		client: &http.Client{},
	}
}
