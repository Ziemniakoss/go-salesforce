// Package salesforce provides a minimal client for the Salesforce REST/Tooling
// Query APIs, with support for API versioning and "query more" pagination.
package gosalesforce

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type QueryResult struct {
	TotalSize      int               `json:"totalSize"`
	Done           bool              `json:"done"`
	NextRecordsURL string            `json:"nextRecordsUrl,omitempty"`
	Records        []json.RawMessage `json:"records"`
}

func (c *SfConnectionWithToken) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *SfConnectionWithToken) version() string {
	if c.APIVersion != "" {
		return c.APIVersion
	}
	return DefaultAPIVersion
}

func (c *SfConnectionWithToken) Query(ctx context.Context, soql string, tooling bool) (*QueryResult, error) {
	endpoint := fmt.Sprintf("/services/data/v%s/query", c.version())
	if tooling {
		endpoint = fmt.Sprintf("/services/data/v%s/tooling/query", c.version())
	}

	u := c.InstanceURL + endpoint + "?" + url.Values{"q": {soql}}.Encode()
	return c.doQueryRequest(ctx, u)
}

func (c *SfConnectionWithToken) QueryMore(ctx context.Context, nextRecordsURL string) (*QueryResult, error) {
	u := c.InstanceURL + nextRecordsURL
	return c.doQueryRequest(ctx, u)
}

// QueryAll runs soql and transparently pages through every "query more"
// continuation, returning the combined records. Use this when you just want
// every row and don't need to control pagination yourself.
func (c *SfConnectionWithToken) QueryAll(ctx context.Context, soql string, tooling bool) ([]json.RawMessage, error) {
	result, err := c.Query(ctx, soql, tooling)
	if err != nil {
		return nil, err
	}

	all := result.Records
	for !result.Done && result.NextRecordsURL != "" {
		result, err = c.QueryMore(ctx, result.NextRecordsURL)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Records...)
	}
	return all, nil
}

func (c *SfConnectionWithToken) doQueryRequest(ctx context.Context, fullURL string) (*QueryResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("salesforce: building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("salesforce: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var sfErrs []SfError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&sfErrs); jsonErr == nil && len(sfErrs) > 0 {
			return nil, fmt.Errorf("salesforce: %s (status %d): %s [%s]",
				fullURL, resp.StatusCode, sfErrs[0].Message, sfErrs[0].ErrorCode)
		}
		return nil, fmt.Errorf("salesforce: %s returned status %d", fullURL, resp.StatusCode)
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("salesforce: decoding response: %w", err)
	}
	return &result, nil
}
