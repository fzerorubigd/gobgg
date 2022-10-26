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
		Items []map[string]interface{} `json:"items"`
	}

	rankRequest struct {
		Item map[string]interface{} `json:"item"`
	}
)

// SetRank tries to add rank for an item (experimental)
func (bgg *BGG) SetRank(ctx context.Context, objectID int64, rate float64) error {
	if rate <= 0 || rate > 10 {
		return fmt.Errorf("invalid rate range: %f", rate)
	}
	name := bgg.GetActiveUsername()
	if len(bgg.GetActiveCookies()) == 0 || name == "" {
		return fmt.Errorf("call login first")
	}

	usr, err := bgg.GetUser(ctx, name)
	if err != nil {
		return err
	}

	params := map[string]string{
		"objectid":   fmt.Sprint(objectID),
		"objecttype": "thing",
		"userid":     fmt.Sprint(usr.UserID),
	}
	u := bgg.buildURL(apiCollectionUrl, params)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("failed to create the request: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	bgg.requestCookies(req)

	resp, err := bgg.client.Do(req)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}
	var items rankResponse
	if err := json.Unmarshal(b, &items); err != nil {
		return fmt.Errorf("decoding JSON failed: %w", err)
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
	b, err = json.Marshal(update)
	if err != nil {
		return fmt.Errorf("marshaling json failed: %w", err)
	}

	url := bgg.buildURL(fmt.Sprintf(apiCollectionItemUrl, collid), nil)
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create the request: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	bgg.requestCookies(req)

	resp2, err := bgg.client.Do(req)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode < 200 || resp2.StatusCode >= 300 {
		return fmt.Errorf("failed with status code")
	}

	return nil
}
