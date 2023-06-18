package gobgg

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Limiter is a rate limiter interface from the go.uber.org/ratelimit
// the package itself is not needed, but can be used with this client
type Limiter interface {
	// Take should block to make sure that the RPS is met.
	Take() time.Time
}

type noOpLimiter struct {
}

func (noOpLimiter) Take() time.Time {
	return time.Time{}
}

// BGG is the client for the boardgame geek site
type BGG struct {
	host    string
	scheme  string
	client  *http.Client
	limiter Limiter

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
	u, err := url.Parse(path)
	if err != nil {
		u = &url.URL{
			Scheme: bgg.scheme,
			Host:   bgg.host,
			Path:   path,
		}
	}

	if u.Host == "" {
		u.Host = bgg.host
	}

	if u.Scheme == "" {
		u.Scheme = bgg.scheme
	}

	q := u.Query()
	for i := range args {
		q.Set(i, args[i])
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (bgg *BGG) do(req *http.Request) (*http.Response, error) {
	bgg.limiter.Take()
	return bgg.client.Do(req)
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

// SetLimiter can use tos et a limiter to limit the api call to the BGG
func SetLimiter(limiter Limiter) OptionSetter {
	return func(bgg *BGG) {
		bgg.limiter = limiter
	}
}

// NewBGGClient returns a new client
func NewBGGClient(opt ...OptionSetter) *BGG {
	result := &BGG{
		host:    "boardgamegeek.com",
		scheme:  "https",
		client:  &http.Client{},
		lock:    sync.RWMutex{},
		limiter: noOpLimiter{},
	}

	for i := range opt {
		opt[i](result)
	}

	return result
}
