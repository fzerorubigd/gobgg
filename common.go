package gobgg

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

// ItemType is the item type for the search api
type ItemType string

const (
	// RPGItemType for rpg
	RPGItemType ItemType = "rpgitem"
	// VideGameType for video game
	VideGameType ItemType = "videogame"
	// BoardGameType for board game
	BoardGameType ItemType = "boardgame"
	// BoardGameAccessoryType for accessory
	BoardGameAccessoryType ItemType = "boardgameaccessory"
	// BoardGameExpansionType for expansion
	BoardGameExpansionType ItemType = "boardgameexpansion"
)

// SimpleString is the string
type SimpleString struct {
	Text  string `xml:",chardata"`
	Value string `xml:"value,attr"`
}

func (s *SimpleString) String() string {
	if s == nil {
		return ""
	}
	return s.Value
}

// NameStruct is the name from the api
type NameStruct struct {
	Text      string `xml:",chardata"`
	Type      string `xml:"type,attr"`
	Sortindex string `xml:"sortindex,attr"`
	Value     string `xml:"value,attr"`
}

// PollStruct is the poll
type PollStruct struct {
	Text       string `xml:",chardata"`
	Name       string `xml:"name,attr"`
	Title      string `xml:"title,attr"`
	Totalvotes string `xml:"totalvotes,attr"`
	Results    []struct {
		Text       string `xml:",chardata"`
		Numplayers string `xml:"numplayers,attr"`
		Result     []struct {
			Text     string `xml:",chardata"`
			Value    string `xml:"value,attr"`
			Numvotes string `xml:"numvotes,attr"`
			Level    string `xml:"level,attr"`
		} `xml:"result"`
	} `xml:"results"`
}

type Statistics struct {
	Text    string `xml:",chardata"`
	Page    string `xml:"page,attr"`
	Ratings struct {
		Text       string `xml:",chardata"`
		Usersrated struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"usersrated"`
		Average struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"average"`
		Bayesaverage struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"bayesaverage"`
		Ranks struct {
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
		Stddev struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"stddev"`
		Median struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"median"`
		Owned struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"owned"`
		Trading struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"trading"`
		Wanting struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"wanting"`
		Wishing struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"wishing"`
		Numcomments struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"numcomments"`
		Numweights struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"numweights"`
		Averageweight struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"averageweight"`
	} `xml:"ratings"`
}

type Poll struct {
}

// PollItem is a single poll  item in a poll
type PollItem struct {
	Name       string
	Title      string
	TotalVotes int
}

// LinkStruct is for the link for the things
type LinkStruct struct {
	Text    string `xml:",chardata"`
	Type    string `xml:"type,attr"`
	ID      int64  `xml:"id,attr"`
	Value   string `xml:"value,attr"`
	Inbound string `xml:"inbound,attr"`
}

// Link is the link
type Link struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Plays struct {
	Total    int64  `json:"total,omitempty"`
	Page     int64  `json:"page,omitempty"`
	UserName string `json:"user_name,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
	Items    []Play `json:"items,omitempty"`
}

type Play struct {
	ID         int64         `json:"id,omitempty"`
	Date       time.Time     `json:"date,omitempty"`
	Quantity   int64         `json:"quantity,omitempty"`
	Length     time.Duration `json:"length,omitempty"`
	Incomplete bool          `json:"incomplete,omitempty"`
	NowInStats bool          `json:"now_in_stats,omitempty"`
	Location   string        `json:"location,omitempty"`
	Comment    string        `json:"comment,omitempty"`
	Item       Item          `json:"item,omitempty"`
	Players    []Player      `json:"players,omitempty"`
}

type Item struct {
	Name string   `json:"name,omitempty"`
	Type ItemType `json:"type,omitempty"`
	ID   int64    `json:"id,omitempty"`
}

type Player struct {
	UserName      string `json:"user_name,omitempty"`
	UserID        string `json:"user_id,omitempty"`
	Name          string `json:"name,omitempty"`
	StartPosition string `json:"start_position,omitempty"`
	Color         string `json:"color,omitempty"`
	Score         int64  `json:"score,omitempty"`
	New           bool   `json:"new,omitempty"`
	Rating        string `json:"rating,omitempty"`
	Win           bool   `json:"win,omitempty"`
}

func nameStructToString(args []NameStruct) (string, []string) {
	var (
		primary   string
		alternate []string
	)

	for _, name := range args {
		switch {
		case name.Type == "primary":
			primary = name.Value
		case name.Type == "alternate":
			alternate = append(alternate, name.Value)
		default:
			log.Printf("Name type %q is not handled, please report it as an issue", name.Type)
		}
	}

	return primary, alternate
}

func linksMap(links []LinkStruct) map[string][]Link {
	m := make(map[string][]Link)
	for _, lnk := range links {
		ln := append(m[lnk.Type], Link{
			ID:   lnk.ID,
			Name: lnk.Value,
		})
		m[lnk.Type] = ln
	}

	return m
}

func safeInt(str string) int64 {
	if str == "" {
		return 0
	}

	i64, _ := strconv.ParseInt(str, 10, 0)
	return i64
}

func safeFloat64(str string) float64 {
	if str == "" {
		return 0
	}

	i64, _ := strconv.ParseFloat(str, 64)
	return i64
}

func safeDate(str string) time.Time {
	ts, err := time.Parse(bggTimeFormat, str)
	if err != nil {
		return time.Time{}
	}

	return ts
}

type bggError struct {
	XMLName xml.Name `xml:"error"`
	Text    string   `xml:",chardata"`
	Message string   `xml:"message"`
}

func decode(r io.Reader, in interface{}) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}

	err = xml.Unmarshal(buf, in)
	if err == nil {
		return nil
	}
	var errType bggError
	if err2 := xml.Unmarshal(buf, &errType); err2 != nil {
		return err
	}

	return fmt.Errorf("error from bgg: %q", errType.Message)
}
