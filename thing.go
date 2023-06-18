package gobgg

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
)

const (
	thingPath = "xmlapi2/thing"
	rankPath  = "api/collectionstatsgraph"
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
	ID           int64   `json:"id,omitempty"`
	Name         string  `json:"name,omitempty"`
	FriendlyName string  `json:"friendly_name,omitempty"`
	Rank         int     `json:"rank,omitempty"`
	BayesAverage float64 `json:"bayes_average,omitempty"`
}

// CollectionItem is the item in collection
type CollectionItem struct {
	ID            int64    `json:"id,omitempty"`
	CollID        int64    `json:"coll_id"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Type          ItemType `json:"type,omitempty"`
	YearPublished int      `json:"year_published,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	Image         string   `json:"image,omitempty"`

	CollectionStatus []string `json:"collection_status,omitempty"`
}

// ThingResult is the result for the thing api
type ThingResult struct {
	ID             int64    `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	AlternateNames []string `json:"alternate_names,omitempty"`
	Type           ItemType `json:"type,omitempty"`
	YearPublished  int      `json:"year_published,omitempty"`

	Thumbnail string `json:"thumbnail,omitempty"`
	Image     string `json:"image,omitempty"`

	MinPlayers int `json:"min_players,omitempty"`
	MaxPlayers int `json:"max_players,omitempty"`

	SuggestedPlayerCount []SuggestedPlayerCount

	// TODO: int?
	MinAge string `json:"min_age,omitempty"`

	PlayTime    string `json:"play_time,omitempty"`
	MinPlayTime string `json:"min_play_time,omitempty"`
	MaxPlayTime string `json:"max_play_time,omitempty"`

	Description string `json:"description,omitempty"`

	Links map[string][]Link `json:"links,omitempty"`

	UsersRated   int     `json:"users_rated,omitempty"`
	AverageRate  float64 `json:"average_rate,omitempty"`
	BayesAverage float64 `json:"bayes_average,omitempty"`

	UsersOwned    int     `json:"users_owned,omitempty"`
	UsersTrading  int     `json:"users_trading,omitempty"`
	UsersWanting  int     `json:"users_wanting,omitempty"`
	UsersWishing  int     `json:"users_wishing,omitempty"`
	NumComments   int     `json:"num_comments,omitempty"`
	NumWeight     int     `json:"num_weight,omitempty"`
	AverageWeight float64 `json:"average_weight,omitempty"`

	RankTotal int                   `json:"rank_total,omitempty"`
	Family    map[string]FamilyRank `json:"family,omitempty"`
}

// GetThings is the get things API entry point
func (bgg *BGG) GetThings(ctx context.Context, setters ...GetOptionSetter) ([]ThingResult, error) {
	opt := GetThingOption{}

	for i := range setters {
		setters[i](&opt)
	}

	if len(opt.ids) == 0 {
		return nil, errors.New("at least one id is required")
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

	resp, err := bgg.do(req)
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
		spc, err := getSuggestedPoll(result.Item[i].Poll)
		if err != nil {
			return nil, err
		}
		ret[i] = ThingResult{
			ID:                   result.Item[i].ID,
			Type:                 ItemType(result.Item[i].Type),
			YearPublished:        int(safeInt(result.Item[i].Yearpublished.Value)),
			Description:          html.UnescapeString(result.Item[i].Description),
			Thumbnail:            result.Item[i].Thumbnail,
			Image:                result.Item[i].Image,
			MinPlayers:           int(safeInt(result.Item[i].Minplayers.Value)),
			MaxPlayers:           int(safeInt(result.Item[i].Maxplayers.Value)),
			SuggestedPlayerCount: spc,
			MinAge:               result.Item[i].Minage.Value,
			PlayTime:             result.Item[i].Playingtime.Value,
			MinPlayTime:          result.Item[i].Minplaytime.Value,
			MaxPlayTime:          result.Item[i].Maxplaytime.Value,
			UsersRated:           int(safeInt(result.Item[i].Statistics.Ratings.Usersrated.Value)),
			AverageRate:          safeFloat64(result.Item[i].Statistics.Ratings.Average.Value),
			BayesAverage:         safeFloat64(result.Item[i].Statistics.Ratings.Bayesaverage.Value),
			UsersOwned:           int(safeInt(result.Item[i].Statistics.Ratings.Owned.Value)),
			UsersTrading:         int(safeInt(result.Item[i].Statistics.Ratings.Trading.Value)),
			UsersWanting:         int(safeInt(result.Item[i].Statistics.Ratings.Wanting.Value)),
			UsersWishing:         int(safeInt(result.Item[i].Statistics.Ratings.Wishing.Value)),
			NumComments:          int(safeInt(result.Item[i].Statistics.Ratings.Numcomments.Value)),
			NumWeight:            int(safeInt(result.Item[i].Statistics.Ratings.Numweights.Value)),
			AverageWeight:        safeFloat64(result.Item[i].Statistics.Ratings.Averageweight.Value),
			Family:               make(map[string]FamilyRank),
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

// RankBreakDown shows the rank break down in BGG website
type RankBreakDown [10]int64

func (rb RankBreakDown) Total() int64 {
	var total int64
	for i := range rb {
		total += rb[i]
	}

	return total
}
func (rb RankBreakDown) Average() float64 {
	var (
		total int64
		av    float64
	)
	for i := range rb {
		total += rb[i]
		av += float64(int64(i+1) * rb[i])
	}

	return av / float64(total)
}

func (rb RankBreakDown) BayesianAverage(added int64) float64 {
	m := float64(added) * 5.5
	total := added
	for i := range rb {
		m += float64(i+1) * float64(rb[i])
		total += rb[i]
	}

	return m / float64(total)
}

type rankBreakDownResponse struct {
	Type       string                 `json:"type"`
	Options    map[string]interface{} `json:"options"`
	Formatters []interface{}          `json:"formatters"`
	Data       struct {
		Cols []interface{} `json:"cols"`
		Rows []struct {
			C []struct {
				V interface{} `json:"v"`
				F interface{} `json:"f"`
			} `json:"c"`
		} `json:"rows"`
	} `json:"data"`
}

func (bgg *BGG) GetRankBreakDown(ctx context.Context, gameID int64) (RankBreakDown, error) {
	args := map[string]string{
		"objectid":   fmt.Sprint(gameID),
		"objecttype": "thing",
		"type":       "BarChart",
	}
	rbd := RankBreakDown{}
	u := bgg.buildURL(rankPath, args)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return rbd, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := bgg.do(req)
	if err != nil {
		return rbd, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()
	var result rankBreakDownResponse
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return rbd, fmt.Errorf("read data failed: %w", err)
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return rbd, fmt.Errorf("JSON decoding failed: %w", err)
	}

	for _, row := range result.Data.Rows {
		if len(row.C) < 2 {
			return rbd, fmt.Errorf("[RESPONSE] invalid row data")
		}

		num := safeIntInterface(row.C[0].V)
		val := safeIntInterface(row.C[1].V)
		if num <= 10 && num > 0 {
			rbd[num-1] = val
		}
	}

	return rbd, nil
}
