package gobgg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

const (
	topPages         = "browse/boardgame/page/%d"
	hotnessPage      = "https://api.geekdo.com/api/hotness"
	geekListPage     = "https://api.geekdo.com/api/listitems"
	trendBestSeller  = "https://api.geekdo.com/api/trends/ownership"
	trendMostPlays   = "https://api.geekdo.com/api/trends/plays"
	trendsTrendPlays = "https://api.geekdo.com/api/trends/plays_delta"
)

type TrendInterval string

const (
	TrendIntervalWeek  TrendInterval = "week"
	TrendIntervalMonth TrendInterval = "month"
)

type hotnessStruct struct {
	Items []struct {
		Objecttype    string `json:"objecttype"`
		Objectid      string `json:"objectid"`
		RepImageid    string `json:"rep_imageid"`
		Delta         int    `json:"delta"`
		Href          string `json:"href"`
		Name          string `json:"name"`
		ID            string `json:"id"`
		Type          string `json:"type"`
		Imageurl      string `json:"imageurl"`
		Images        any    `json:"images"`
		Yearpublished string `json:"yearpublished"`
		Rank          string `json:"rank,omitempty"`
		Description   string `json:"description"`
	} `json:"items"`
}

type trendsStruct struct {
	Items []struct {
		ID   string `json:"id"`
		Item struct {
			Type           string `json:"type"`
			ID             string `json:"id"`
			Name           string `json:"name"`
			Href           string `json:"href"`
			Label          string `json:"label"`
			Labelpl        string `json:"labelpl"`
			HasAngularLink bool   `json:"hasAngularLink"`
			Descriptors    []struct {
				Name         string `json:"name"`
				DisplayValue string `json:"displayValue"`
			} `json:"descriptors"`
			Breadcrumbs []any `json:"breadcrumbs"`
			ImageSets   struct {
				Square100 struct {
					Src   string `json:"src"`
					Src2X string `json:"src@2x"`
				} `json:"square100"`
				Mediacard100 struct {
					Src   string `json:"src"`
					Src2X string `json:"src@2x"`
				} `json:"mediacard100"`
				Mediacard struct {
					Src   string `json:"src"`
					Src2X string `json:"src@2x"`
				} `json:"mediacard"`
			} `json:"imageSets"`
			Imageid       int `json:"imageid"`
			NameSortIndex int `json:"nameSortIndex"`
		} `json:"item"`
		Rank        int    `json:"rank"`
		Description string `json:"description"`
		Delta       int    `json:"delta"`
		Appearances int    `json:"appearances"`
	} `json:"items"`
	Interval string    `json:"interval"`
	EndDate  time.Time `json:"endDate"`
}

type geekListResult struct {
	Data []struct {
		Type   string `json:"type"`
		ID     string `json:"id"`
		Listid string `json:"listid"`
		Item   struct {
			Type           string `json:"type"`
			ID             string `json:"id"`
			Name           string `json:"name"`
			Href           string `json:"href"`
			Label          string `json:"label"`
			Labelpl        string `json:"labelpl"`
			HasAngularLink bool   `json:"hasAngularLink"`
			Descriptors    []struct {
				Name         string `json:"name"`
				DisplayValue string `json:"displayValue"`
			} `json:"descriptors"`
			Breadcrumbs   []any `json:"breadcrumbs"`
			ImageSets     any   `json:"imageSets"`
			Imageid       int   `json:"imageid"`
			NameSortIndex int   `json:"nameSortIndex"`
		} `json:"item"`
		Postdate        time.Time `json:"postdate"`
		Editdate        time.Time `json:"editdate"`
		Body            string    `json:"body"`
		BodyXML         string    `json:"bodyXml"`
		Author          int       `json:"author"`
		Href            string    `json:"href"`
		Imageid         int       `json:"imageid"`
		ImageOverridden bool      `json:"imageOverridden"`
		LinkedImage     any       `json:"linkedImage"`
		RollsEnabled    bool      `json:"rollsEnabled"`
		Links           []struct {
			Rel string `json:"rel"`
			URI string `json:"uri"`
		} `json:"links"`
		RollsCount int `json:"rollsCount"`
		Stats      struct {
			Average float64 `json:"average"`
			Rank    int     `json:"rank"`
		} `json:"stats"`
	} `json:"data"`
	Pagination struct {
		Pageid  int `json:"pageid"`
		PerPage int `json:"perPage"`
		Total   int `json:"total"`
	} `json:"pagination"`
}

type ListItem struct {
	ID          int64
	Name        string
	Description string
}

type IDDelta struct {
	ID    int64
	Delta int
}

type TrendOutput struct {
	ID          int64
	Delta       int
	Appearances int
}

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

func (bgg *BGG) Hotness(ctx context.Context, count int) ([]IDDelta, error) {
	if count < 1 || count > 50 {
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

func (bgg *BGG) trends(ctx context.Context, u string) ([]TrendOutput, error) {
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
	result := trendsStruct{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing failed: %w", err)
	}

	final := make([]TrendOutput, len(result.Items))
	for i := range result.Items {
		final[i].ID = safeInt(result.Items[i].Item.ID)
		final[i].Delta = result.Items[i].Delta
		final[i].Appearances = result.Items[i].Appearances
	}

	return final, nil
}

// BestSellers returns the best sellers of the week starting from the given date. The start day should be Monday, so the
// function will calculate the previous Monday if the given date is not Monday.
func (bgg *BGG) BestSellers(ctx context.Context, start time.Time) ([]TrendOutput, error) {
	start = getPreviousDay(start, time.Monday)
	u := bgg.buildURL(trendBestSeller, map[string]string{
		"interval":  string(TrendIntervalWeek),
		"startDate": start.Format("2006-01-02"),
	})

	return bgg.trends(ctx, u)
}

func (bgg *BGG) MostPlays(ctx context.Context, interval TrendInterval, start time.Time) ([]TrendOutput, error) {
	switch interval {
	case TrendIntervalWeek:
		start = getPreviousDay(start, time.Monday)
	case TrendIntervalMonth:
		start = getStartOfTheMonth(start)
	default:
		return nil, fmt.Errorf("invalid interval: %q", interval)
	}
	u := bgg.buildURL(trendMostPlays, map[string]string{
		"interval":  string(interval),
		"startDate": start.Format("2006-01-02"),
	})

	return bgg.trends(ctx, u)
}

func (bgg *BGG) TrendingPlays(ctx context.Context, interval TrendInterval, start time.Time) ([]TrendOutput, error) {
	switch interval {
	case TrendIntervalWeek:
		start = getPreviousDay(start, time.Monday)
	case TrendIntervalMonth:
		start = getStartOfTheMonth(start)
	default:
		return nil, fmt.Errorf("invalid interval: %q", interval)
	}
	u := bgg.buildURL(trendsTrendPlays, map[string]string{
		"interval":  string(interval),
		"startDate": start.Format("2006-01-02"),
	})

	return bgg.trends(ctx, u)
}

func (bgg *BGG) GeekList(ctx context.Context, geekID int64) ([]ListItem, error) {
	var final []ListItem
	page := 1
	for {
		u := bgg.buildURL(geekListPage, map[string]string{
			"page":   fmt.Sprint(page),
			"listid": fmt.Sprint(geekID),
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
		result := geekListResult{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			return nil, fmt.Errorf("JSON parsing failed: %w", err)
		}
		if len(result.Data) == 0 {
			break
		}
		for i := range result.Data {
			final = append(final, ListItem{
				ID:          safeInt(result.Data[i].Item.ID),
				Name:        result.Data[i].Item.Name,
				Description: result.Data[i].Body,
			})
		}
		page++
	}
	return final, nil
}
