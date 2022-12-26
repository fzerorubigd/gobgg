package gobgg

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

const (
	topPages = "browse/boardgame/page/%d"
)

func (bgg *BGG) TopPages(ctx context.Context, page int) ([]int64, error) {
	u := bgg.buildURL(fmt.Sprintf(topPages, page), map[string]string{})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	resp, err := bgg.client.Do(req)
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
				result = append(result, getID(getAttr(node.Attr, "href"), "boardgame"))
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
