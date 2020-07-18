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

// OptionSetter lets you to modify the internal settings
type OptionSetter func(*BGG)

// SetClient allows you to modify the default client
func SetClient(client *http.Client) OptionSetter {
	return func(bgg *BGG) {
		bgg.client = client
	}
}

// SetHost changes the host, default is boardgamegeek.com
func SetHost(host string) OptionSetter {
	return func(bgg *BGG) {
		bgg.host = host
	}
}

// SetSchema changes the schema, default is https
func SetSchema(schema string) OptionSetter {
	return func(bgg *BGG) {
		bgg.scheme = schema
	}
}

// NewBGGClient returns a new client
func NewBGGClient(opt ...OptionSetter) *BGG {
	result := &BGG{
		host:   "boardgamegeek.com",
		scheme: "https",
		client: &http.Client{},
	}

	for i := range opt {
		opt[i](result)
	}

	return result
}
