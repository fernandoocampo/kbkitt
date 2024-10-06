package webs

import (
	"fmt"
	"io"
	"net/http"
)

func GetWebMediaFile(urlpath string) ([]byte, error) {
	resp, err := http.DefaultClient.Get(urlpath)
	if err != nil {
		return nil, fmt.Errorf("unable to get web media resource: %w", err)
	}

	defer resp.Body.Close()

	mediaData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to get web media resource: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unable to get web media resource: %v", mediaData)
	}

	return mediaData, nil
}

func GetMediaContentType(urlpath string) ([]string, error) {
	resp, err := http.DefaultClient.Head(urlpath)
	if err != nil {
		return nil, fmt.Errorf("unable to get web media content type: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unable to get web media content type: %d", resp.StatusCode)
	}

	value, ok := resp.Header["Content-Type"]
	if !ok {
		return nil, nil
	}

	return value, nil
}
