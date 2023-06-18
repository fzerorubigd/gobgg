package gobgg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

const (
	topPages    = "browse/boardgame/page/%d"
	hotnessPage = "https://api.geekdo.com/api/hotness"
)

func (bgg *BGG) TopPages(ctx context.Context, page int) ([]int64, error) {
	u := bgg.buildURL(fmt.Sprintf(topPages, page), map[string]string{})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	resp, err := bgg.do(req)
	if err != nil {
		return nil, fmt.Errorf("get request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html failed: %w", err)
	}

	var result []int64
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			if getAttr(node.Attr, "class") == "primary" {
				id := getID(getAttr(node.Attr, "href"), "boardgame")
				if id > 0 {
					result = append(result, id)
				}
			}
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)

	return result, nil
}

type hotnessStruct struct {
	Items []struct {
		Objecttype    string      `json:"objecttype"`
		Objectid      string      `json:"objectid"`
		RepImageid    string      `json:"rep_imageid"`
		Delta         int         `json:"delta"`
		Href          string      `json:"href"`
		Name          string      `json:"name"`
		ID            string      `json:"id"`
		Type          string      `json:"type"`
		Imageurl      string      `json:"imageurl"`
		Images        interface{} `json:"images"`
		Yearpublished string      `json:"yearpublished"`
		Rank          string      `json:"rank,omitempty"`
		Description   string      `json:"description"`
	} `json:"items"`
}

type IDDelta struct {
	ID    int64
	Delta int
}

func (bgg *BGG) Hotness(ctx context.Context, count int) ([]IDDelta, error) {
	if count < 1 || count > 100 {
		count = 50
	}
	u := bgg.buildURL(hotnessPage, map[string]string{
		"geeksite":   "boardgame",
		"objecttype": "thing",
		"showcount":  fmt.Sprint(count),
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	resp, err := bgg.do(req)
	if err != nil {
		return nil, fmt.Errorf("get request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}
	result := hotnessStruct{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing failed: %w", err)
	}

	final := make([]IDDelta, len(result.Items))
	for i := range result.Items {
		final[i].ID = safeInt(result.Items[i].ID)
		final[i].Delta = result.Items[i].Delta
	}

	return final, nil
}
