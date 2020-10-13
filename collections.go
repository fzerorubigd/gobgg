package gobgg

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

const collectionPath = "/xmlapi2/collection"

type collectionItems struct {
	XMLName    xml.Name `xml:"items"`
	Text       string   `xml:",chardata"`
	Totalitems string   `xml:"totalitems,attr"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Pubdate    string   `xml:"pubdate,attr"`
	Item       []struct {
		Text       string `xml:",chardata"`
		Objecttype string `xml:"objecttype,attr"`
		Objectid   string `xml:"objectid,attr"`
		Subtype    string `xml:"subtype,attr"`
		Collid     string `xml:"collid,attr"`
		Name       struct {
			Text      string `xml:",chardata"`
			Sortindex string `xml:"sortindex,attr"`
		} `xml:"name"`
		Yearpublished string `xml:"yearpublished"`
		Image         string `xml:"image"`
		Thumbnail     string `xml:"thumbnail"`
		Stats         struct {
			Text        string `xml:",chardata"`
			Minplayers  string `xml:"minplayers,attr"`
			Maxplayers  string `xml:"maxplayers,attr"`
			Minplaytime string `xml:"minplaytime,attr"`
			Maxplaytime string `xml:"maxplaytime,attr"`
			Playingtime string `xml:"playingtime,attr"`
			Numowned    string `xml:"numowned,attr"`
			Rating      struct {
				Text         string       `xml:",chardata"`
				Value        string       `xml:"value,attr"`
				Usersrated   SimpleString `xml:"usersrated"`
				Average      SimpleString `xml:"average"`
				Bayesaverage SimpleString `xml:"bayesaverage"`
				Stddev       SimpleString `xml:"stddev"`
				Median       SimpleString `xml:"median"`
				Ranks        struct {
					Text string `xml:",chardata"`
					Rank []struct {
						Text         string `xml:",chardata"`
						Type         string `xml:"type,attr"`
						ID           string `xml:"id,attr"`
						Name         string `xml:"name,attr"`
						Friendlyname string `xml:"friendlyname,attr"`
						Value        string `xml:"value,attr"`
						Bayesaverage string `xml:"bayesaverage,attr"`
					} `xml:"rank"`
				} `xml:"ranks"`
			} `xml:"rating"`
		} `xml:"stats"`
		Status struct {
			Text             string `xml:",chardata"`
			Own              string `xml:"own,attr"`
			Prevowned        string `xml:"prevowned,attr"`
			Fortrade         string `xml:"fortrade,attr"`
			Want             string `xml:"want,attr"`
			Wanttoplay       string `xml:"wanttoplay,attr"`
			Wanttobuy        string `xml:"wanttobuy,attr"`
			Wishlist         string `xml:"wishlist,attr"`
			Preordered       string `xml:"preordered,attr"`
			Lastmodified     string `xml:"lastmodified,attr"`
			Wishlistpriority string `xml:"wishlistpriority,attr"`
		} `xml:"status"`
		Numplays        string `xml:"numplays"`
		Comment         string `xml:"comment"`
		Wishlistcomment string `xml:"wishlistcomment"`
		Version         struct {
			Text string `xml:",chardata"`
			Item struct {
				Text          string       `xml:",chardata"`
				Type          string       `xml:"type,attr"`
				ID            string       `xml:"id,attr"`
				Thumbnail     string       `xml:"thumbnail"`
				Image         string       `xml:"image"`
				Link          []LinkStruct `xml:"link"`
				Name          NameStruct   `xml:"name"`
				Yearpublished SimpleString `xml:"yearpublished"`
				Productcode   SimpleString `xml:"productcode"`
				Width         SimpleString `xml:"width"`
				Length        SimpleString `xml:"length"`
				Depth         SimpleString `xml:"depth"`
				Weight        SimpleString `xml:"weight"`
			} `xml:"item"`
		} `xml:"version"`
		Originalname string `xml:"originalname"`
	} `xml:"item"`
}

// GetCollectionOptions is the option used to handle the collection request
type GetCollectionOptions struct {
	version          bool
	subtype          string
	excludesubtype   string
	brief            bool
	stats            bool
	own              bool
	rated            bool
	played           bool
	comment          bool
	trade            bool
	want             bool
	wishlist         bool
	wishlistPriority int
	preordered       bool
	wanttoplay       bool
	wanttobuy        bool
	prevowned        bool
	hasparts         bool
	wantparts        bool
	minrating        int
	rating           int
	minbggrating     int
	bggrating        int
	minplays         int
	maxplays         int
	showprivate      bool
	collid           int
	modifiedsince    *time.Time
}

func (c *GetCollectionOptions) toMap() map[string]string {
	result := map[string]string{}
	if c.version {
		result["version"] = "1"
	}

	if c.subtype != "" {
		result["subtype"] = c.subtype
	}

	if c.excludesubtype != "" {
		result["excludesubtype"] = c.excludesubtype
	}

	return result
}

// CollectionOptionSetter is the option setter for the get collection
type CollectionOptionSetter func(*GetCollectionOptions)

// SetVersion Returns version info for each item in your collection.
func SetVersion(version bool) CollectionOptionSetter {
	return func(co *GetCollectionOptions) {
		co.version = version
	}
}

// SetSubType  Specifies which collection you want to retrieve.
// TYPE may be boardgame, boardgameexpansion, boardgameaccessory,
// rpgitem, rpgissue, or videogame; the default is boardgame
func SetSubType(subType ItemType) CollectionOptionSetter {
	return func(co *GetCollectionOptions) {
		co.subtype = string(subType)
	}
}

// SetExcludeSubtype  Specifies which subtype you want to exclude from the results.
func SetExcludeSubtype(subType ItemType) CollectionOptionSetter {
	return func(co *GetCollectionOptions) {
		co.excludesubtype = string(subType)
	}
}

// GetCollection is to get the collections of a user
func (bgg *BGG) GetCollection(ctx context.Context, username string, options ...CollectionOptionSetter) (interface{}, error) {
	opt := GetCollectionOptions{}

	for i := range options {
		options[i](&opt)
	}

	args := opt.toMap()
	args["username"] = username

	u := bgg.buildURL(collectionPath, args)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// If there is a cookie
	for i := range bgg.cookies {
		req.AddCookie(bgg.cookies[i])
	}

	var (
		resp  *http.Response
		delay = time.Second
	)
	for i := 1; ; i++ {
		resp, err = bgg.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("http call failed: %w", err)
		}

		if resp.StatusCode == http.StatusAccepted {
			resp.Body.Close() // we don't need it

			delay += time.Duration(i) * time.Second
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	var result collectionItems
	if err = dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	return result, nil
}
