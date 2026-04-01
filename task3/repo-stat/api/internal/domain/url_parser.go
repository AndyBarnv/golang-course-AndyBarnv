package domain

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseGitHubURL(rawURL string) (string, string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid url format")
	}

	path := strings.Trim(parsedURL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid github repository url format")
	}

	return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
}
