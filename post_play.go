package gobgg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
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

func (bgg *BGG) PostPlay(ctx context.Context, play Play) error {
	if len(bgg.cookies) == 0 {
		return fmt.Errorf("call login first")
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
		return fmt.Errorf("create payload failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create the request: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	for i := range bgg.cookies {
		req.AddCookie(bgg.cookies[i])
	}
	resp, err := bgg.client.Do(req)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	d, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(d))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed with status code")
	}

	return nil
}
