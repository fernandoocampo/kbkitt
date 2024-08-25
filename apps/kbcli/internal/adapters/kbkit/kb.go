package kbkitt

import "net/http"

type Setup struct {
	URL string
}

type Client struct {
	client *http.Client
}

func NewClient(settings Setup) *Client {
	newClient := Client{
		client: newHTTPClient(settings),
	}

	return &newClient
}

func newHTTPClient(settings Setup) *http.Client {
	newHTTPClient := http.Client{}

	return &newHTTPClient
}
