package gobgg

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	playsPath     = "xmlapi2/plays"
	bggTimeFormat = "2006-01-02"
)

// playsResponse is the response for the plays
type playsResponse struct {
	XMLName    xml.Name `xml:"plays"`
	Text       string   `xml:",chardata"`
	Username   string   `xml:"username,attr"`
	UserID     string   `xml:"userid,attr"`
	Total      string   `xml:"total,attr"`
	Page       string   `xml:"page,attr"`
	TermsOfUse string   `xml:"termsofuse,attr"`
	Play       []struct {
		Text       string `xml:",chardata"`
		ID         string `xml:"id,attr"`
		Date       string `xml:"date,attr"`
		Quantity   string `xml:"quantity,attr"`
		Length     string `xml:"length,attr"`
		Incomplete string `xml:"incomplete,attr"`
		NowInStats string `xml:"nowinstats,attr"`
		Location   string `xml:"location,attr"`
		Item       struct {
			Text       string `xml:",chardata"`
			Name       string `xml:"name,attr"`
			ObjectType string `xml:"objecttype,attr"`
			ObjectID   string `xml:"objectid,attr"`
			Subtypes   struct {
				Text    string         `xml:",chardata"`
				Subtype []SimpleString `xml:"subtype"`
			} `xml:"subtypes"`
		} `xml:"item"`
		Players struct {
			Text   string `xml:",chardata"`
			Player []struct {
				Text          string `xml:",chardata"`
				Username      string `xml:"username,attr"`
				Userid        string `xml:"userid,attr"`
				Name          string `xml:"name,attr"`
				StartPosition string `xml:"startposition,attr"`
				Color         string `xml:"color,attr"`
				Score         string `xml:"score,attr"`
				New           string `xml:"new,attr"`
				Rating        string `xml:"rating,attr"`
				Win           string `xml:"win,attr"`
			} `xml:"player"`
		} `xml:"players"`
		Comments string `xml:"comments"`
	} `xml:"play"`
}

// PlaysOption is used to handle func option ins plays api
type PlaysOption struct {
	userName         string
	gameID           int
	minDate, maxDate time.Time

	page int
}

// SetUserName for the plays
func SetUserName(name string) PlaysOptionSetter {
	return func(opt *PlaysOption) {
		opt.userName = name
	}
}

// SetGameID to set the game id
func SetGameID(id int) PlaysOptionSetter {
	return func(opt *PlaysOption) {
		opt.gameID = id
	}
}

// SetPageNumber set the current page
func SetPageNumber(page int) PlaysOptionSetter {
	return func(opt *PlaysOption) {
		opt.page = page
	}
}

// SetDateRangeMin set the min date
func SetDateRangeMin(min time.Time) PlaysOptionSetter {
	return func(option *PlaysOption) {
		option.minDate = min
	}
}

// SetDateRangeMax set the min date
func SetDateRangeMax(max time.Time) PlaysOptionSetter {
	return func(option *PlaysOption) {
		option.maxDate = max
	}
}

// PlaysOptionSetter is used to handle the func option in plays api
type PlaysOptionSetter func(*PlaysOption)

// Plays using plays api of the bgg, it get the list of requested items
func (bgg *BGG) Plays(ctx context.Context, setter ...PlaysOptionSetter) (*Plays, error) {
	opt := PlaysOption{}
	for i := range setter {
		setter[i](&opt)
	}

	args := map[string]string{}

	if opt.gameID == 0 && opt.userName == "" {
		return nil, errors.New("at least game id or username should be there")
	}

	if opt.userName != "" {
		args["username"] = opt.userName
	}

	if opt.gameID > 0 {
		args["id"] = fmt.Sprint(opt.gameID)
	}

	if opt.page > 0 {
		args["page"] = fmt.Sprint(opt.page)
	}

	if !opt.minDate.IsZero() {
		args["mindate"] = opt.minDate.Format(bggTimeFormat)
	}

	if !opt.maxDate.IsZero() {
		args["maxdate"] = opt.maxDate.Format(bggTimeFormat)
	}

	u := bgg.buildURL(playsPath, args)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := bgg.do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	var pr playsResponse
	if err = decode(resp.Body, &pr); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	result := Plays{
		Total:    safeInt(pr.Total),
		Page:     safeInt(pr.Page),
		UserName: pr.Username,
		UserID:   safeInt(pr.UserID),
		Items:    make([]Play, 0, len(pr.Play)),
	}

	for _, ply := range pr.Play {
		item := Play{
			ID:         safeInt(ply.ID),
			Date:       safeDate(ply.Date),
			Quantity:   safeInt(ply.Quantity),
			Length:     time.Duration(safeInt(ply.Length)) * time.Second,
			Incomplete: safeInt(ply.Incomplete) != 0,
			NowInStats: safeInt(ply.NowInStats) != 0,
			Location:   ply.Location,
			Comment:    ply.Comments,
			Item: Item{
				Name: ply.Item.Name,
				Type: ItemType(ply.Item.ObjectType),
				ID:   safeInt(ply.Item.ObjectID),
			},
			Players: make([]Player, 0, len(ply.Players.Player)),
		}

		for _, plr := range ply.Players.Player {
			item.Players = append(item.Players, Player{
				UserName:      plr.Username,
				UserID:        plr.Userid,
				Name:          plr.Name,
				StartPosition: plr.StartPosition,
				Color:         plr.Color,
				Score:         safeInt(plr.Score),
				New:           safeInt(plr.New) != 0,
				Rating:        plr.Rating,
				Win:           safeInt(plr.Win) != 0,
			})
		}

		result.Items = append(result.Items, item)
	}

	return &result, nil
}
