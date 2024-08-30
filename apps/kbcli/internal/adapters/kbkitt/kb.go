package kbkitt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
)

type Setup struct {
	URL string
}

type Client struct {
	client *http.Client
	host   string
}

type NewKBResponse struct {
	ID string `json:"id"`
}

// url patterns
const (
	kbURL        = "%s/kbs"
	getKBByIDURL = "%s/kbs/%s"
)

// http values
const (
	appJsonContentType = "application/json"
	keyParam           = "key"
	keywordParam       = "keyword"
	limitParam         = "limit"
	offsetParam        = "offset"
)

func NewClient(settings Setup) *Client {
	newClient := Client{
		client: newHTTPClient(settings),
		host:   settings.URL,
	}

	return &newClient
}

func newHTTPClient(settings Setup) *http.Client {
	newHTTPClient := http.Client{}

	return &newHTTPClient
}

func (c *Client) Create(ctx context.Context, newKB kbs.NewKB) (string, error) {
	postBody, err := json.Marshal(newKB)
	if err != nil {
		return "", fmt.Errorf("unable to marsal kb data: %w")
	}

	resp, err := c.client.Post(c.getKBURL(), appJsonContentType, bytes.NewBuffer(postBody))
	if err != nil {
		return "", fmt.Errorf("unable to create new kb: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response after trying to create a new kb: %w", err)
	}

	if isNotSuccess(resp.StatusCode) {
		return "", fmt.Errorf("server failed to create new kb: %s", string(respBody))
	}

	var newKBResponse NewKBResponse
	err = json.Unmarshal(respBody, &newKBResponse)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshall new kb response: %w", err)
	}

	return newKBResponse.ID, nil
}

func (c *Client) Search(ctx context.Context, filter kbs.KBQueryFilter) ([]kbs.KBItem, error) {
	req, err := http.NewRequest(http.MethodGet, c.getKBURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build get kb request: %w", err)
	}

	q := req.URL.Query()
	q.Add(keyParam, filter.Key)
	q.Add(keywordParam, filter.Keyword)
	q.Add(limitParam, fmt.Sprintf("%d", filter.Limit))
	q.Add(offsetParam, fmt.Sprintf("%d", filter.Offset))

	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get kb by key: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response after trying to get a kb by key: %w", err)
	}

	if isNotSuccess(resp.StatusCode) {
		return nil, fmt.Errorf("server failed to get kb by key: %s", string(respBody))
	}

	var kbResponse []kbs.KBItem
	err = json.Unmarshal(respBody, &kbResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall kb response: %w", err)
	}

	return kbResponse, nil
}

func (c *Client) Get(ctx context.Context, id string) (*kbs.KB, error) {
	resp, err := c.client.Get(c.getGetKBByIDURL(id))
	if err != nil {
		return nil, fmt.Errorf("unable to get kb: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response after trying to get a kb: %w", err)
	}

	if isNotSuccess(resp.StatusCode) {
		return nil, fmt.Errorf("server failed to get kb: %s", string(respBody))
	}

	var kbResponse kbs.KB
	err = json.Unmarshal(respBody, &kbResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall kb response: %w", err)
	}

	return &kbResponse, nil
}

func (c *Client) getKBURL() string {
	return fmt.Sprintf(kbURL, c.host)
}

func (c *Client) getGetKBByIDURL(id string) string {
	return fmt.Sprintf(getKBByIDURL, c.host, id)
}

func isNotSuccess(statusCode int) bool {
	return !(statusCode >= 200 && statusCode < 300)
}
