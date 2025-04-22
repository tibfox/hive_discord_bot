package hiveTools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Vote struct {
	Voter string `json:"voter"`
}

type VotesResponse struct {
	Result []Vote `json:"result"`
}

type PostResultCreated struct {
	Created string `json:"created"`
}

type RPCResponseCreated struct {
	Result PostResultCreated `json:"result"`
	ID     int               `json:"id"`
}

func PostAge(postAuthor *string, postPermlink *string) (string, error) {

	// Construct the API request body
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "condenser_api.get_content",
		"params":  []interface{}{postAuthor, postPermlink},
		"id":      1,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "5skipped: failed to marshal request body", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make the HTTP POST request to the Hive API (Appbase RPC)
	resp, err := http.Post("https://api.hive.blog", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "5skipped api.hive.blog not reachable", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return "5skipped api.hive.blog not reachable (non-OK response)", fmt.Errorf("received non-OK response from Hive API: %v", resp.Status)
	}

	// Parse the JSON response
	var response RPCResponseCreated
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "5skipped: failed to decode JSON", fmt.Errorf("failed to decode JSON response: %v", err)
	}
	layout := "2006-01-02T15:04:05"
	createdTime, err := time.Parse(layout, response.Result.Created)
	if err != nil {
		fmt.Println("%v", err)
		return "5skipped: failed to get creation date", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if time.Since(createdTime) > 24*time.Hour {
		return "4post is > 24h old", nil
	}

	// Account has not voted
	return "2post is recent enough", nil
}

func HasVoted(postAuthor string, postPermlink string, account string) (string, error) {
	fmt.Println("Author:", postAuthor)
	fmt.Println("Permlink:", postPermlink)
	// Construct the API request body
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "condenser_api.get_active_votes",
		"params":  []interface{}{postAuthor, postPermlink},
		"id":      1,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "5skipped: failed to marshal request body", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make the HTTP POST request to the Hive API (Appbase RPC)
	resp, err := http.Post("https://api.hive.blog", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "5skipped api.hive.blog not reachable", fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return "5skipped api.hive.blog not reachable (non-OK response)", fmt.Errorf("received non-OK response from Hive API: %v", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "5skipped: failed to read response body", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var votesResponse VotesResponse
	if err := json.Unmarshal(respBody, &votesResponse); err != nil {
		return "5skipped: failed to decode JSON", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	// Debug: Print the list of voters to see if the account exists
	fmt.Println("Voters List:")
	for _, vote := range votesResponse.Result {
		fmt.Println(vote.Voter) // Debugging: Print all voters
	}

	// Loop through the votes and check if the account has voted
	for _, vote := range votesResponse.Result {
		if strings.EqualFold(strings.ToLower(vote.Voter), strings.ToLower(account)) {
			return "4already voted", nil // Account has voted on the post
		}
	}

	// Account has not voted
	return "2not voted", nil
}

type Account struct {
	Name        string `json:"name"`
	VotingPower int    `json:"voting_power"`
}

type AccountResponse struct {
	Result []Account `json:"result"`
}

func GetVotingPowerPercent(username string) (string, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "condenser_api.get_accounts",
		"params":  [][]string{{username}},
		"id":      1,
	}
	minVP := 80.0

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "5skipped: failed to marshal request", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post("https://api.hive.blog", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "5skipped: request failed", fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	var accountResp AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountResp); err != nil {
		return "5skipped: failed to decode response", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(accountResp.Result) == 0 {
		return fmt.Sprintf("5skipped: no account data found for %s", username), fmt.Errorf("no account data found for %s", username)
	}

	rawPower := accountResp.Result[0].VotingPower
	percent := float64(rawPower) / 100.0 // Convert to percentage
	if minVP >= percent {
		return fmt.Sprintf("4failed: We don't have enough VP (currently %.2f%% VP)", percent), nil
	} else {
		return fmt.Sprintf("2we have %.2f%% VP", percent), nil
	}

}
