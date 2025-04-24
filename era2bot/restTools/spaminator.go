package restTools

import (
	"encoding/json"
	"io"
	"net/http"
)

func CheckSpaminator(author *string) string {
	// Define the URL
	url := "https://spaminator.me/api/bl/all.json"

	resp, err := http.Get(url)
	if err != nil {
		return "5skipped: spaminator not available (failed completely)"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Fatal("Failed to read response body:", err)
		return "5skipped: spaminator not available (response body)"
	}

	// Define a structure that matches the JSON format
	var data struct {
		Result []string `json:"result"`
	}

	// Parse JSON into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		// log.Fatal("JSON unmarshal failed:", err)
		return "5skipped: spaminator not available (json not parseable)"
	}

	// Check if target exists in the result
	found := false
	for _, val := range data.Result {
		if val == *author {
			found = true
			break
		}
	}

	if found {
		return "4author blacklisted"
	} else {
		return "2author is okay"
	}
}
