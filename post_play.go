package gobgg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type createPlayer struct {
	Name          string  `json:"name"`
	Username      string  `json:"username"`
	Userid        int64   `json:"userid,omitempty"`
	Avatarfile    string  `json:"avatarfile,omitempty"`
	Avatar        bool    `json:"avatar,omitempty"`
	Selected      bool    `json:"selected"`
	Color         string  `json:"color"`
	Score         string  `json:"score"`
	Win           bool    `json:"win,omitempty"`
	New           bool    `json:"new,omitempty"`
	Disambiguator float64 `json:"disambiguator,omitempty"`
}

type createPlayPayload struct {
	Players        []createPlayer `json:"players"`
	Quantity       int            `json:"quantity"`
	Date           time.Time      `json:"date"`
	Twitter        bool           `json:"twitter"`
	Locationfilter string         `json:"locationfilter"`
	Location       string         `json:"location"`
	Minutes        int            `json:"minutes"`
	Hours          int            `json:"hours"`
	Incomplete     bool           `json:"incomplete"`
	Comments       string         `json:"comments"`
	Userfilter     string         `json:"userfilter"`
	Objecttype     string         `json:"objecttype"`
	Objectid       string         `json:"objectid"`
	Playdate       string         `json:"playdate"`
	Length         int            `json:"length"`
	Ajax           int            `json:"ajax"`
	Action         string         `json:"action"`
}

type createPlayResponse struct {
	Playid   string `json:"playid,omitempty"`
	Numplays int    `json:"numplays,omitempty"`
	HTML     string `json:"html,omitempty"`
	Error    string `json:"error,omitempty"`
}

// PostPlay save a play record, you should be logged in, and it returns the number of plays after you save this one
func (bgg *BGG) PostPlay(ctx context.Context, play *Play) (int, error) {
	if len(bgg.GetActiveCookies()) == 0 {
		return 0, fmt.Errorf("call Login first")
	}
	payload := createPlayPayload{
		Playdate:   play.Date.Format(bggTimeFormat),
		Comments:   play.Comment,
		Length:     int(play.Length.Minutes()),
		Twitter:    false,
		Minutes:    int(play.Length.Minutes()),
		Location:   play.Location,
		Objectid:   fmt.Sprint(play.Item.ID),
		Hours:      0,
		Quantity:   1,
		Action:     "save",
		Date:       time.Now(),
		Players:    nil,
		Objecttype: string(play.Item.Type),
		Ajax:       1,
	}

	for _, py := range play.Players {
		payload.Players = append(payload.Players, createPlayer{
			Name:          py.Name,
			Username:      py.UserName,
			Userid:        safeInt(py.UserID),
			Avatarfile:    "",
			Avatar:        false,
			Selected:      false,
			Color:         py.Color,
			Score:         fmt.Sprint(py.Score),
			Win:           py.Win,
			New:           py.New,
			Disambiguator: 0,
		})
	}

	u := bgg.buildURL("geekplay.php", nil)
	b, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("create payload failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("failed to create the request: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	bgg.requestCookies(req)

	resp, err := bgg.do(req)
	if err != nil {
		return 0, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("failed with status code")
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read the payload: %w", err)
	}

	var cpr createPlayResponse
	if err := json.Unmarshal(b, &cpr); err != nil {
		return 0, fmt.Errorf("invalid json response: %w", err)
	}

	if cpr.Error != "" {
		return 0, fmt.Errorf("create play failed: %q", cpr.Error)
	}

	play.ID = safeInt(cpr.Playid)

	return cpr.Numplays, nil
}
