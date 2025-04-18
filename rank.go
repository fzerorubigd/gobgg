package gobgg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Rank API is based on the site API, and not the official api

const (
	apiCollectionUrl     = "api/collections"
	apiCollectionItemUrl = "api/collectionitems/%s"
)

type (
	rankResponse struct {
		Items []map[string]any `json:"items"`
	}

	rankRequest struct {
		Item map[string]any `json:"item"`
	}
)

func (bgg *BGG) myCollections(ctx context.Context, objectID int64) (*rankResponse, error) {
	name := bgg.GetActiveUsername()
	if len(bgg.GetActiveCookies()) == 0 || name == "" {
		return nil, fmt.Errorf("call Login first")
	}

	usr, err := bgg.GetUser(ctx, name)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"objectid":   fmt.Sprint(objectID),
		"objecttype": "thing",
		"userid":     fmt.Sprint(usr.UserID),
	}
	u := bgg.buildURL(apiCollectionUrl, params)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the request: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	bgg.requestCookies(req)

	resp, err := bgg.do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	var items rankResponse
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, fmt.Errorf("decoding JSON failed: %w", err)
	}

	return &items, nil
}

// SetRank tries to add rank for an item (experimental)
func (bgg *BGG) SetRank(ctx context.Context, objectID int64, rate float64) error {
	if rate <= 0 || rate > 10 {
		return fmt.Errorf("invalid rate range: %f", rate)
	}

	items, err := bgg.myCollections(ctx, objectID)
	if err != nil {
		return err
	}

	if len(items.Items) == 0 {
		return fmt.Errorf("no item found")
	}

	item := items.Items[0]
	collid, ok := item["collid"].(string)
	if !ok {
		return fmt.Errorf("the response has no collid")
	}

	item["rating"] = rate

	update := rankRequest{
		Item: item,
	}
	b, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("marshaling json failed: %w", err)
	}

	url := bgg.buildURL(fmt.Sprintf(apiCollectionItemUrl, collid), nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create the request: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	bgg.requestCookies(req)

	resp, err := bgg.do(req)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed with status code")
	}

	return nil
}
