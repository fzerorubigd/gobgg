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

const (
	// WishListPriorityMustHave BGA definition
	WishListPriorityMustHave = iota + 1
	// WishListPriorityLoveToHave BGA definition
	WishListPriorityLoveToHave
	// WishListPriorityLikeToHave BGA definition
	WishListPriorityLikeToHave
	// WishListPriorityThinkingAboutIt BGA definition
	WishListPriorityThinkingAboutIt
	// WishListPriorityDoNotBuy BGA definition
	WishListPriorityDoNotBuy
)

var priorityToText = map[int]string{
	WishListPriorityMustHave:        "musthave",
	WishListPriorityLoveToHave:      "lovetohave",
	WishListPriorityLikeToHave:      "liketohave",
	WishListPriorityThinkingAboutIt: "thinkingaboutit",
	WishListPriorityDoNotBuy:        "donotbuy",
}

type collectionStatus struct {
	Text             string `xml:",chardata"`
	Own              int    `xml:"own,attr"`
	Prevowned        int    `xml:"prevowned,attr"`
	Fortrade         int    `xml:"fortrade,attr"`
	Want             int    `xml:"want,attr"`
	Wanttoplay       int    `xml:"wanttoplay,attr"`
	Wanttotrade      int    `xml:"wanttotrade,attr"`
	Wanttobuy        int    `xml:"wanttobuy,attr"`
	Wishlist         int    `xml:"wishlist,attr"`
	Preordered       int    `xml:"preordered,attr"`
	Lastmodified     string `xml:"lastmodified,attr"`
	Wishlistpriority int    `xml:"wishlistpriority,attr"`
}

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
		Collid     int64  `xml:"collid,attr"`
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
		Status          collectionStatus `xml:"status"`
		Numplays        int              `xml:"numplays"`
		Comment         string           `xml:"comment"`
		Wishlistcomment string           `xml:"wishlistcomment"`
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
	// brief          bool
	// stats          bool
	options      []CollectionType
	minrating    int
	rating       int
	minbggrating int
	bggrating    int
	minplays     int
	maxplays     int
	// showprivate   bool
	ids           []int64
	collID        int64
	modifiedsince *time.Time
}

func (c *GetCollectionOptions) toMap() map[string]string {
	result := map[string]string{}

	setIf := func(cond bool, key, value string) {
		if cond {
			result[key] = value
		}
	}
	setIf(c.version, "version", "1")
	setIf(c.subtype != "", "subtype", c.subtype)
	setIf(c.excludesubtype != "", "excludesubtype", c.excludesubtype)
	setIf(c.minbggrating > 0, "minbggrating", fmt.Sprint(c.minbggrating))
	setIf(c.bggrating > 0, "bggrating", fmt.Sprint(c.bggrating))
	setIf(c.minrating > 0, "minrating", fmt.Sprint(c.minrating))
	setIf(c.rating > 0, "rating", fmt.Sprint(c.rating))
	setIf(c.minplays > 0, "minplays", fmt.Sprint(c.minplays))
	setIf(c.maxplays > 0, "maxplays", fmt.Sprint(c.maxplays))
	setIf(c.collID > 0, "collid", fmt.Sprint(c.collID))
	if c.modifiedsince != nil {
		result["modifiedsince"] = c.modifiedsince.Format("06-01-02")
	}

	for i := range c.options {
		result[string(c.options[i])] = "1"
	}

	if len(c.ids) > 0 {
		st := make([]string, 0, len(c.ids))
		for i := range c.ids {
			if c.ids[i] > 0 {
				st = append(st, fmt.Sprint(c.ids[i]))
			}
		}
		result["id"] = strings.Join(st, ",")
	}

	return result
}

// CollectionOptionSetter is the option setter for the get collection
type CollectionOptionSetter func(*GetCollectionOptions)

func SetIDs(ids ...int64) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		options.ids = ids
	}
}

func SetCollID(collID int64) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		options.collID = collID
	}
}

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

// SetModifiedSince to set the modified since flag
func SetModifiedSince(t time.Time) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		options.modifiedsince = &t
	}
}

// SetMinPlays show games with min plays
func SetMinPlays(plays int) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		options.minplays = plays
	}
}

// SetMaxPlays show games with max plays
func SetMaxPlays(plays int) CollectionOptionSetter {
	return func(options *GetCollectionOptions) {
		options.maxplays = plays
	}
}

func statusToStringArray(status *collectionStatus, played int) []string {
	var result []string
	setIf := func(cond bool, txt string) {
		if cond {
			result = append(result, txt)
		}
	}
	setIf(status.Own != 0, string(CollectionTypeOwn))
	setIf(status.Want != 0, string(CollectionTypeWant))
	setIf(status.Wanttobuy != 0, string(CollectionTypeWantToBuy))
	setIf(status.Wanttoplay != 0, string(CollectionTypeWantToPlay))
	setIf(status.Wanttotrade != 0, string(CollectionTypeTrade))
	setIf(status.Wishlist != 0, string(CollectionTypeWishList))
	prio, ok := priorityToText[status.Wishlistpriority]
	setIf(ok, prio)
	setIf(status.Preordered != 0, string(CollectionTypePreorder))
	setIf(status.Prevowned != 0, string(CollectionTypePrevOwned))
	setIf(status.Fortrade != 0, string(CollectionTypeTrade))
	setIf(played > 0, string(CollectionTypePlayed))
	return result
}

// GetCollection is to get the collections of a user
func (bgg *BGG) GetCollection(ctx context.Context, username string, options ...CollectionOptionSetter) ([]CollectionItem, error) {
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

	bgg.requestCookies(req)

	var (
		resp  *http.Response
		delay = time.Second
	)
	for i := 1; ; i++ {
		resp, err = bgg.do(req)
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

	var result collectionItems
	if err = decode(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	ret := make([]CollectionItem, len(result.Item))
	for i := range result.Item {
		ret[i] = CollectionItem{
			ID:               result.Item[i].Objectid,
			CollID:           result.Item[i].Collid,
			Name:             result.Item[i].Name.Text,
			Description:      strings.Trim(html.UnescapeString(result.Item[i].Text), "\n\t "),
			Type:             ItemType(result.Item[i].Objecttype),
			YearPublished:    int(safeInt(result.Item[i].Yearpublished)),
			Thumbnail:        result.Item[i].Thumbnail,
			Image:            result.Item[i].Image,
			CollectionStatus: statusToStringArray(&result.Item[i].Status, result.Item[i].Numplays),
		}
	}
	return ret, nil
}
