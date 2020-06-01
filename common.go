package gobgg

import (
	"strconv"

	"github.com/go-acme/lego/log"
)

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

// Poll is a single poll in bgg
type Poll struct {
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
	ID   int64
	Name string
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
