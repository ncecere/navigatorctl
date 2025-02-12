package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client handles API communication
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// resolveTeamIdentifier gets team ID from either ID or alias
func (c *Client) resolveTeamIdentifier(identifier string) (string, error) {
	// If it looks like a UUID, use it directly
	if len(identifier) == 36 && identifier[8] == '-' {
		return identifier, nil
	}

	// Otherwise, try to find the team by alias
	teams, err := c.ListTeams()
	if err != nil {
		return "", err
	}

	for _, team := range teams {
		if team.TeamAlias == identifier {
			return team.TeamID, nil
		}
	}

	return "", fmt.Errorf("team not found: %s", identifier)
}

// ListTeamMembers gets all members in a team
func (c *Client) ListTeamMembers(identifier string) ([]TeamMember, error) {
	teamID, err := c.resolveTeamIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/team/info?team_id=%s", c.BaseURL, teamID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var response TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return response.TeamInfo.MembersWithRoles, nil
}

// AddTeamMember adds a new member to a team
func (c *Client) AddTeamMember(identifier string, member TeamMember) (*TeamResponse, error) {
	teamID, err := c.resolveTeamIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	request := AddMemberRequest{
		Member: []TeamMember{member},
		TeamID: teamID,
	}

	url := fmt.Sprintf("%s/team/member_add", c.BaseURL)
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var response TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &response, nil
}

// RemoveTeamMember removes a member from a team
func (c *Client) RemoveTeamMember(identifier string, member TeamMember) (*TeamResponse, error) {
	teamID, err := c.resolveTeamIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	request := RemoveMemberRequest{
		UserID:    member.UserID,
		UserEmail: member.UserEmail,
		TeamID:    teamID,
	}

	url := fmt.Sprintf("%s/team/member_delete", c.BaseURL)
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var response TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &response, nil
}

// GetKeyInfo gets detailed information about a specific key
func (c *Client) GetKeyInfo(keyID string) (*KeyResponse, error) {
	url := fmt.Sprintf("%s/key/info?key=%s", c.BaseURL, keyID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var keyResponse KeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&keyResponse); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &keyResponse, nil
}

// ListTeamKeys gets all API keys for a team
func (c *Client) ListTeamKeys(identifier string) ([]KeyResponse, error) {
	teamID, err := c.resolveTeamIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/key/list?page=1&size=100&team_id=%s", c.BaseURL, teamID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var listResponse KeyListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Get detailed info for each key
	var keys []KeyResponse
	for _, keyID := range listResponse.Keys {
		keyInfo, err := c.GetKeyInfo(keyID)
		if err != nil {
			return nil, fmt.Errorf("getting key info for %s: %w", keyID, err)
		}
		keys = append(keys, *keyInfo)
	}

	return keys, nil
}

// ListTeams gets all teams
func (c *Client) ListTeams() ([]Team, error) {
	url := fmt.Sprintf("%s/team/list", c.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var teams []Team
	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return teams, nil
}

// GetTeamInfo gets detailed information about a team by ID or alias
func (c *Client) GetTeamInfo(identifier string) (*Team, error) {
	teams, err := c.ListTeams()
	if err != nil {
		return nil, err
	}

	for _, team := range teams {
		if team.TeamID == identifier || team.TeamAlias == identifier {
			return &team, nil
		}
	}

	return nil, fmt.Errorf("team not found: %s", identifier)
}

// GetUserInfo gets detailed information about a user by ID or email
func (c *Client) GetUserInfo(identifier string) (*UserResponse, error) {
	var url string
	if strings.Contains(identifier, "@") {
		url = fmt.Sprintf("%s/user/info?email=%s", c.BaseURL, identifier)
	} else {
		url = fmt.Sprintf("%s/user/info?user_id=%s", c.BaseURL, identifier)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var response UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &response, nil
}

// ListUsers gets all users
func (c *Client) ListUsers() (*UserResponse, error) {
	url := fmt.Sprintf("%s/user/info", c.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	var response UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &response, nil
}
