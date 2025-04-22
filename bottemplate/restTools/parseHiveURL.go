package restTools

import (
	"fmt"
	"strings"
)

func ParseHiveLink(url string) (string, string, error) {
	atIndex := strings.Index(url, "@")
	if atIndex == -1 {
		return "", "", fmt.Errorf("link is not in a valid format")
	}

	rest := url[atIndex+1:]
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("link is not in a valid format")
	}

	return parts[0], parts[1], nil
}
