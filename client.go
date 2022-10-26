package gobgg

import (
	"net/http"
	"net/url"
	"sync"
)

// BGG is the client for the boardgame geek site
type BGG struct {
	host   string
	scheme string
	client *http.Client

	// I prefer not to use the cookie jar since this is simpler
	cookies  []*http.Cookie
	username string

	lock sync.RWMutex
}

// GetActiveCookies return the cookies if the cookies are available
func (bgg *BGG) GetActiveCookies() []*http.Cookie {
	bgg.lock.RLock()
	defer bgg.lock.RUnlock()

	return bgg.cookies
}

// GetActiveUsername return the username that the current cookie are based on
func (bgg *BGG) GetActiveUsername() string {
	bgg.lock.RLock()
	defer bgg.lock.RUnlock()

	return bgg.username
}

func (bgg *BGG) requestCookies(req *http.Request) {
	bgg.lock.RLock()
	defer bgg.lock.RUnlock()
	// If there is a cookie
	for i := range bgg.cookies {
		req.AddCookie(bgg.cookies[i])
	}
}

func (bgg *BGG) buildURL(path string, args map[string]string) string {
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

// OptionSetter modify the internal settings
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

// SetCookies set the cookies for this object, in case the user is logged in already
func SetCookies(username string, c []*http.Cookie) OptionSetter {
	return func(bgg *BGG) {
		bgg.cookies = c
		bgg.username = username
	}
}

// NewBGGClient returns a new client
func NewBGGClient(opt ...OptionSetter) *BGG {
	result := &BGG{
		host:   "boardgamegeek.com",
		scheme: "https",
		client: &http.Client{},
		lock:   sync.RWMutex{},
	}

	for i := range opt {
		opt[i](result)
	}

	return result
}
