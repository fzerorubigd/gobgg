package gobgg

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

const (
	personsPath = "xmlapi2/person"
)

type personItem struct {
	XMLName    xml.Name `xml:"items"`
	Text       string   `xml:",chardata"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Item       struct {
		Text      string `xml:",chardata"`
		Type      string `xml:"type,attr"`
		ID        string `xml:"id,attr"`
		Thumbnail string `xml:"thumbnail"`
		Image     string `xml:"image"`
	} `xml:"item"`
}

// PersonImage is the persons image and thumbnail
type PersonImage struct {
	ID        int64
	Thumbnail string
	Image     string
}

func (bgg *BGG) PersonImage(ctx context.Context, id int64) (*PersonImage, error) {
	u := bgg.buildURL(playsPath, map[string]string{
		"id": fmt.Sprint(id),
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := bgg.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	var pr personItem
	if err = decode(resp.Body, &pr); err != nil {
		return nil, fmt.Errorf("XML decoding failed: %w", err)
	}

	return &PersonImage{
		ID:        id,
		Thumbnail: pr.Item.Thumbnail,
		Image:     pr.Item.Image,
	}, nil
}
