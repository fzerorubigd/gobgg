package gobgg

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strings"
)

const (
	thingPath = "xmlapi2/thing"
)

type thingItems struct {
	XMLName    xml.Name `xml:"items"`
	Text       string   `xml:",chardata"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Item       []struct {
		Text          string       `xml:",chardata"`
		Type          string       `xml:"type,attr"`
		ID            int64        `xml:"id,attr"`
		Thumbnail     string       `xml:"thumbnail"`
		Image         string       `xml:"image"`
		Name          []NameStruct `xml:"name"`
		Description   string       `xml:"description"`
		Yearpublished SimpleString `xml:"yearpublished"`
		Minplayers    SimpleString `xml:"minplayers"`
		Maxplayers    SimpleString `xml:"maxplayers"`
		Poll          []PollStruct `xml:"poll"`
		Playingtime   SimpleString `xml:"playingtime"`
		Minplaytime   SimpleString `xml:"minplaytime"`
		Maxplaytime   SimpleString `xml:"maxplaytime"`
		Minage        SimpleString `xml:"minage"`
		Link          []LinkStruct `xml:"link"`
		Videos        struct {
			Text  string `xml:",chardata"`
			Total string `xml:"total,attr"`
			Video []struct {
				Text     string `xml:",chardata"`
				ID       string `xml:"id,attr"`
				Title    string `xml:"title,attr"`
				Category string `xml:"category,attr"`
				Language string `xml:"language,attr"`
				Link     string `xml:"link,attr"`
				Username string `xml:"username,attr"`
				Userid   string `xml:"userid,attr"`
				Postdate string `xml:"postdate,attr"`
			} `xml:"video"`
		} `xml:"videos"`
		Versions struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text      string `xml:",chardata"`
				Type      string `xml:"type,attr"`
				ID        string `xml:"id,attr"`
				Thumbnail string `xml:"thumbnail"`
				Image     string `xml:"image"`
				Link      []struct {
					Text    string `xml:",chardata"`
					Type    string `xml:"type,attr"`
					ID      string `xml:"id,attr"`
					Value   string `xml:"value,attr"`
					Inbound string `xml:"inbound,attr"`
				} `xml:"link"`
				Name []struct {
					Text      string `xml:",chardata"`
					Type      string `xml:"type,attr"`
					Sortindex string `xml:"sortindex,attr"`
					Value     string `xml:"value,attr"`
				} `xml:"name"`
				Yearpublished SimpleString `xml:"yearpublished"`
				Productcode   SimpleString `xml:"productcode"`
				Width         SimpleString `xml:"width"`
				Length        SimpleString `xml:"length"`
				Depth         SimpleString `xml:"depth"`
				Weight        SimpleString `xml:"weight"`
			} `xml:"item"`
		} `xml:"versions"`
		Comments struct {
			Text       string `xml:",chardata"`
			Page       string `xml:"page,attr"`
			Totalitems string `xml:"totalitems,attr"`
			Comment    []struct {
				Text     string `xml:",chardata"`
				Username string `xml:"username,attr"`
				Rating   string `xml:"rating,attr"`
				Value    string `xml:"value,attr"`
			} `xml:"comment"`
		} `xml:"comments"`
		Marketplacelistings struct {
			Text    string `xml:",chardata"`
			Listing []struct {
				Text     string       `xml:",chardata"`
				Listdate SimpleString `xml:"listdate"`
				Price    struct {
					Text     string `xml:",chardata"`
					Currency string `xml:"currency,attr"`
					Value    string `xml:"value,attr"`
				} `xml:"price"`
				Condition SimpleString `xml:"condition"`
				Notes     SimpleString `xml:"notes"`
				Link      struct {
					Text  string `xml:",chardata"`
					Href  string `xml:"href,attr"`
					Title string `xml:"title,attr"`
				} `xml:"link"`
			} `xml:"listing"`
		} `xml:"marketplacelistings"`
		Statistics Statistics `xml:"statistics"`
	} `xml:"item"`
}

// GetThingOption is the options for the GetThing api
type GetThingOption struct {
	ids []string
}

// GetOptionSetter is the option setter for the GetThing api
type GetOptionSetter func(*GetThingOption)

// GetThingIDs is for setting IDs
func GetThingIDs(ids ...int64) GetOptionSetter {
	return func(opt *GetThingOption) {
		for i := range ids {
			opt.ids = append(opt.ids, fmt.Sprint(ids[i]))
		}
	}
}

type FamilyRank struct {
	ID                 int64
	Name, FriendlyName string
	Rank               int
	BayesAverage       float64
}

// ThingResult is the result for the thing api
type ThingResult struct {
	ID             int64
	Name           string
	AlternateNames []string
	Type           ItemType
	YearPublished  int

	Thumbnail string
	Image     string

	MinPlayers int
	MaxPlayers int

	// TODO: int?
	MinAge string

	PlayTime    string
	MinPlayTime string
	MaxPlayTime string

	Description string

	Links map[string][]Link

	UsersRated   int
	AverageRate  float64
	BayesAverage float64

	UsersOwned    int
	UsersTrading  int
	UsersWanting  int
	UsersWishing  int
	NumComments   int
	NumWeight     int
	AverageWeight float64

	RankTotal int
	Family    map[string]FamilyRank
}

// GetThings is the get things API entry point
func (bgg *BGG) GetThings(ctx context.Context, setters ...GetOptionSetter) ([]ThingResult, error) {
	opt := GetThingOption{}

	for i := range setters {
		setters[i](&opt)
	}

	args := map[string]string{
		"id":    strings.Join(opt.ids, ","),
		"stats": "1",
	}

	u := bgg.buildURL(thingPath, args)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := bgg.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()
	var result thingItems
	if err = decode(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	ret := make([]ThingResult, len(result.Item))
	for i := range result.Item {
		ret[i] = ThingResult{
			ID:            result.Item[i].ID,
			Type:          ItemType(result.Item[i].Type),
			YearPublished: int(safeInt(result.Item[i].Yearpublished.Value)),
			Description:   html.UnescapeString(result.Item[i].Description),
			Thumbnail:     result.Item[i].Thumbnail,
			Image:         result.Item[i].Image,
			MinPlayers:    int(safeInt(result.Item[i].Minplayers.Value)),
			MaxPlayers:    int(safeInt(result.Item[i].Maxplayers.Value)),
			MinAge:        result.Item[i].Minage.Value,
			PlayTime:      result.Item[i].Playingtime.Value,
			MinPlayTime:   result.Item[i].Minplaytime.Value,
			MaxPlayTime:   result.Item[i].Maxplaytime.Value,
			UsersRated:    int(safeInt(result.Item[i].Statistics.Ratings.Usersrated.Value)),
			AverageRate:   safeFloat64(result.Item[i].Statistics.Ratings.Average.Value),
			BayesAverage:  safeFloat64(result.Item[i].Statistics.Ratings.Bayesaverage.Value),
			UsersOwned:    int(safeInt(result.Item[i].Statistics.Ratings.Owned.Value)),
			UsersTrading:  int(safeInt(result.Item[i].Statistics.Ratings.Trading.Value)),
			UsersWanting:  int(safeInt(result.Item[i].Statistics.Ratings.Wanting.Value)),
			UsersWishing:  int(safeInt(result.Item[i].Statistics.Ratings.Wishing.Value)),
			NumComments:   int(safeInt(result.Item[i].Statistics.Ratings.Numcomments.Value)),
			NumWeight:     int(safeInt(result.Item[i].Statistics.Ratings.Numweights.Value)),
			AverageWeight: safeFloat64(result.Item[i].Statistics.Ratings.Averageweight.Value),
			Family:        make(map[string]FamilyRank),
		}

		for _, r := range result.Item[i].Statistics.Ratings.Ranks.Rank {
			if r.Type == "subtype" && r.Name == "boardgame" {
				ret[i].RankTotal = int(safeInt(r.Value))
				continue
			}

			if r.Type == "family" {
				ret[i].Family[r.Name] = FamilyRank{
					ID:           safeInt(r.ID),
					Name:         r.Name,
					FriendlyName: r.Friendlyname,
					Rank:         int(safeInt(r.Value)),
					BayesAverage: safeFloat64(r.Bayesaverage),
				}
			}
		}

		ret[i].Name, ret[i].AlternateNames = nameStructToString(result.Item[i].Name)
		ret[i].Links = linksMap(result.Item[i].Link)
	}

	return ret, nil
}
