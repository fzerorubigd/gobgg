package gobgg

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

// searchItems is the result of the search in xmlapi2 bgg
type searchItems struct {
	XMLName    xml.Name `xml:"items"`
	Text       string   `xml:",chardata"`
	Total      string   `xml:"total,attr"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Item       []struct {
		Text          string       `xml:",chardata"`
		Type          string       `xml:"type,attr"`
		ID            int64        `xml:"id,attr"`
		Name          []NameStruct `xml:"name"`
		YearPublished struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"yearpublished"`
	} `xml:"item"`
}

// SearchResult is the result for the search
type SearchResult struct {
	ID             int64
	Name           string
	AlternateNames []string
	Type           ItemType
	YearPublished  int // Zero means no data
}

const searchPath = "xmlapi2/search"

// SearchOption is used to handle func option ins earch api
type SearchOption struct {
	types []string
	exact bool
}

// SearchOptionSetter is used to handle the func option in search api
type SearchOptionSetter func(*SearchOption)

// SearchExact set the exact argument for bgg
func SearchExact() SearchOptionSetter {
	return func(opt *SearchOption) {
		opt.exact = true
	}
}

// SearchTypes set the valid types for the api
func SearchTypes(types ...ItemType) SearchOptionSetter {
	return func(opt *SearchOption) {
		opt.types = make([]string, len(types))
		for i := range types {
			opt.types[i] = string(types[i])
		}
	}
}

// Search using search api of the bgg, it get the list of requested items
func (bgg *BGG) Search(ctx context.Context, query string, setter ...SearchOptionSetter) ([]SearchResult, error) {
	opt := SearchOption{}
	for i := range setter {
		setter[i](&opt)
	}

	args := map[string]string{
		"query": query,
	}
	if opt.exact {
		args["exact"] = "1"
	}

	if len(opt.types) > 0 {
		args["type"] = strings.Join(opt.types, ",")
	}

	u := bgg.buildURL(searchPath, args)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := bgg.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	var result searchItems
	if err = dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	ret := make([]SearchResult, len(result.Item))
	for i := range result.Item {
		ret[i] = SearchResult{
			ID:            result.Item[i].ID,
			Type:          ItemType(result.Item[i].Type),
			YearPublished: int(safeInt(result.Item[i].YearPublished.Value)),
		}

		ret[i].Name, ret[i].AlternateNames = nameStructToString(result.Item[i].Name)
	}

	return ret, nil
}
