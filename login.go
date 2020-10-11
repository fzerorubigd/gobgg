package gobgg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

const loginPath = "login/api/v1"

// Login tries to login into the bgg using the credentials and returns the cookies required for next calls
func (bgg *BGG) Login(ctx context.Context, username, password string) error {
	payload := map[string]interface{}{
		"credentials": map[string]string{
			"username": username,
			"password": password,
		},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create the payload: %w", err)
	}

	u := bgg.buildURL(loginPath, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create the request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := bgg.client.Do(req)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	d, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(d))

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return fmt.Errorf("maybe, invalid username/password")
	}

	bgg.cookies = resp.Cookies()

	return nil
}
