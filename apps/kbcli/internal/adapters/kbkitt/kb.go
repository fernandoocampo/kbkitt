package kbkitt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/kbs"
)

type Setup struct {
	URL            string
	RequestTimeout time.Duration
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

// default values
const (
	requestTimeoutDefault = 5 * time.Second
)

func NewClient(settings Setup) *Client {
	newClient := Client{
		client: newHTTPClient(settings),
		host:   settings.URL,
	}

	return &newClient
}

func newHTTPClient(settings Setup) *http.Client {
	newHTTPClient := http.Client{
		Timeout: settings.getTimeout(),
	}

	return &newHTTPClient
}

func (c *Client) Create(ctx context.Context, newKB kbs.NewKB) (string, error) {
	postBody, err := json.Marshal(newKB)
	if err != nil {
		return "", fmt.Errorf("unable to marsal kb data: %w")
	}

	request, err := http.NewRequest(http.MethodPost, c.getKBURL(), bytes.NewBuffer(postBody))
	if err != nil {
		return "", fmt.Errorf("unable to create new kb request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return "", kbs.NewServerErrorWithWrapper("unable to create new kb", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response after trying to create a new kb: %w", err)
	}

	if isServerError(resp.StatusCode) {
		return "", kbs.NewServerError(fmt.Sprintf("server failed to create new kb: %s", string(respBody)))
	}

	if isClientError(resp.StatusCode) {
		return "", kbs.NewClientError(fmt.Sprintf("invalid request: %s", string(respBody)))
	}

	var newKBResponse NewKBResponse
	err = json.Unmarshal(respBody, &newKBResponse)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshall new kb response: %w", err)
	}

	return newKBResponse.ID, nil
}

func (c *Client) Search(ctx context.Context, filter kbs.KBQueryFilter) (*kbs.SearchResult, error) {
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

	var kbResponse kbs.SearchResult
	err = json.Unmarshal(respBody, &kbResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall kb response: %w", err)
	}

	return &kbResponse, nil
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

func (c *Client) Update(ctx context.Context, kb *kbs.KB) error {
	postBody, err := json.Marshal(kb)
	if err != nil {
		return fmt.Errorf("unable to marsal kb data: %w")
	}

	request, err := http.NewRequest(http.MethodPatch, c.getKBURL(), bytes.NewBuffer(postBody))
	if err != nil {
		return fmt.Errorf("unable to create update kb request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return kbs.NewServerErrorWithWrapper("unable to update new kb", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response after trying to update kb: %w", err)
	}

	if isServerError(resp.StatusCode) {
		return kbs.NewServerError(fmt.Sprintf("server failed to update kb: %s", string(respBody)))
	}

	if isClientError(resp.StatusCode) {
		return kbs.NewClientError(fmt.Sprintf("invalid request: %s", string(respBody)))
	}

	return nil
}

func (c *Client) getKBURL() string {
	return fmt.Sprintf(kbURL, c.host)
}

func (c *Client) getGetKBByIDURL(id string) string {
	return fmt.Sprintf(getKBByIDURL, c.host, id)
}

func (s *Setup) getTimeout() time.Duration {
	if s.RequestTimeout <= 0 {
		return requestTimeoutDefault
	}

	return s.RequestTimeout
}

func isClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func isServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

func isNotSuccess(statusCode int) bool {
	return !(statusCode >= 200 && statusCode < 300)
}
