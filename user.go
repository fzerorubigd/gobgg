package gobgg

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

const (
	userPath = "xmlapi2/user"
)

type userResponse struct {
	XMLName    xml.Name `xml:"user"`
	Text       string   `xml:",chardata"`
	ID         int64    `xml:"id,attr"`
	Name       string   `xml:"name,attr"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Firstname  struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"firstname"`
	Lastname struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"lastname"`
	Avatarlink struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"avatarlink"`
	Yearregistered struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"yearregistered"`
	Lastlogin struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"lastlogin"`
	Stateorprovince struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"stateorprovince"`
	Country struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"country"`
	Webaddress struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"webaddress"`
	Xboxaccount struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"xboxaccount"`
	Wiiaccount struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"wiiaccount"`
	Psnaccount struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"psnaccount"`
	Battlenetaccount struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"battlenetaccount"`
	Steamaccount struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"steamaccount"`
	Traderating struct {
		Text  string `xml:",chardata"`
		Value string `xml:"value,attr"`
	} `xml:"traderating"`
}

// User is a single bgg user
type User struct {
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Year       int    `json:"year"`
	AvatarLink string `json:"avatar_link,omitempty"`
}

// GetUser return the user from BGG if exists
func (bgg *BGG) GetUser(ctx context.Context, username string) (*User, error) {
	u := bgg.buildURL(userPath, map[string]string{"name": username})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	if bgg.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bgg.bearerToken)
	}

	resp, err := bgg.do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	var result userResponse
	if err = decode(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	usr := User{
		UserID:     result.ID,
		UserName:   result.Name,
		FirstName:  result.Firstname.Value,
		LastName:   result.Lastname.Value,
		Year:       int(safeInt(result.Yearregistered.Value)),
		AvatarLink: result.Avatarlink.Value,
	}

	return &usr, nil
}
