package gobgg

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

const collectionPath = "/xmlapi2/collection"

// CollectionType is the bgg collection type
type CollectionType string

const (
	// CollectionTypeOwn is the owned items
	CollectionTypeOwn CollectionType = "own"
	// CollectionTypeRated is rated items
	CollectionTypeRated CollectionType = "rated"
	// CollectionTypePlayed is played items
	CollectionTypePlayed CollectionType = "played"
	// CollectionTypeComment is commented items
	CollectionTypeComment CollectionType = "comment"
	// CollectionTypeTrade is in trade list item
	CollectionTypeTrade CollectionType = "trade"
	// CollectionTypeWant is in want list item
	CollectionTypeWant CollectionType = "want"
	// CollectionTypeWishList is the wishlist items
	CollectionTypeWishList CollectionType = "wishlist"
	// CollectionTypePreorder is the pre orders item
	CollectionTypePreorder CollectionType = "preorder"
	// CollectionTypeWantToPlay is want to play items
	CollectionTypeWantToPlay CollectionType = "wanttoplay"
	// CollectionTypeWantToBuy is want to buy items
	CollectionTypeWantToBuy CollectionType = "wanttobuy"
	// CollectionTypePrevOwned is previously owned items
	CollectionTypePrevOwned CollectionType = "prevowned"
	// CollectionTypeHasParts has parts item
	CollectionTypeHasParts CollectionType = "hasparts"
	// CollectionTypeWantParts want parts item
	CollectionTypeWantParts CollectionType = "wantparts"
)

type collectionItems struct {
	XMLName    xml.Name `xml:"items"`
	Text       string   `xml:",chardata"`
	Totalitems string   `xml:"totalitems,attr"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Pubdate    string   `xml:"pubdate,attr"`
	Item       []struct {
		Text       string `xml:",chardata"`
		Objecttype string `xml:"objecttype,attr"`
		Objectid   int64  `xml:"objectid,attr"`
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
			Own              int    `xml:"own,attr"`
			Prevowned        int    `xml:"prevowned,attr"`
			Fortrade         int    `xml:"fortrade,attr"`
			Want             int    `xml:"want,attr"`
			Wanttoplay       int    `xml:"wanttoplay,attr"`
			Wanttobuy        int    `xml:"wanttobuy,attr"`
			Wishlist         int    `xml:"wishlist,attr"`
			Preordered       int    `xml:"preordered,attr"`
			Lastmodified     string `xml:"lastmodified,attr"`
			Wishlistpriority int    `xml:"wishlistpriority,attr"`
		} `xml:"status"`
		Numplays        int    `xml:"numplays"`
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
	version        bool
	subtype        string
	excludesubtype string
	brief          bool
	stats          bool
	options        []CollectionType
	minrating      int
	rating         int
	minbggrating   int
	bggrating      int
	minplays       int
	maxplays       int
	showprivate    bool
	collid         int
	modifiedsince  *time.Time
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

	if c.minbggrating > 0 {
		result["minbggrating"] = fmt.Sprint(c.minbggrating)
	}

	if c.bggrating > 0 {
		result["bggrating"] = fmt.Sprint(c.bggrating)
	}

	if c.minrating > 0 {
		result["minrating"] = fmt.Sprint(c.minrating)
	}

	if c.rating > 0 {
		result["rating"] = fmt.Sprint(c.rating)
	}

	for i := range c.options {
		result[string(c.options[i])] = "1"
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

// SetCollectionTypes returns the collection types
func SetCollectionTypes(typ ...CollectionType) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		options.options = typ
	}
}

// SetMinRating is the minimum personal rating for this item for this user
func SetMinRating(rate int) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		if rate > 0 && rate <= 10 {
			options.minrating = rate
		}
	}
}

// SetRating is the exact personal rating for this item for this user
func SetRating(rate int) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		if rate > 0 && rate <= 10 {
			options.rating = rate
		}
	}
}

// SetBGGRating is the exact bgg rating for this item
func SetBGGRating(rate int) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		if rate > 0 && rate <= 10 {
			options.bggrating = rate
		}
	}
}

// SetMinBGGRating is the minimum bgg rating for this item
func SetMinBGGRating(rate int) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		if rate > 0 && rate <= 10 {
			options.minbggrating = rate
		}
	}
}

// GetCollection is to get the collections of a user
func (bgg *BGG) GetCollection(ctx context.Context, username string, options ...CollectionOptionSetter) ([]ThingResult, error) {
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

		if resp.StatusCode == http.StatusOK {
			break
		}
		resp.Body.Close() // we don't need it
		if resp.StatusCode != http.StatusAccepted {
			return nil, fmt.Errorf("invalid status: %q", resp.Status)
		}

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
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	var result collectionItems
	if err = dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	ret := make([]ThingResult, len(result.Item))
	for i := range result.Item {
		ret[i] = ThingResult{
			ID:            result.Item[i].Objectid,
			Type:          ItemType(result.Item[i].Objecttype),
			YearPublished: int(safeInt(result.Item[i].Yearpublished)),
			Description:   strings.Trim(html.UnescapeString(result.Item[i].Text), "\n\t "),
			Thumbnail:     result.Item[i].Thumbnail,
			Image:         result.Item[i].Image,
		}

		ret[i].Name = result.Item[i].Name.Text
	}
	return ret, nil
}
