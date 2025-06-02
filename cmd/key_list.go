// cmd/key_list.go

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type KeyListResponse struct {
	Keys        []string `json:"keys"`
	TotalCount  int      `json:"total_count"`
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
}

var keyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api-url")
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiURL == "" || apiKey == "" {
			fmt.Fprintln(os.Stderr, "API URL and API Key are required")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/key/list?page=1&size=100&return_full_object=true&include_team_keys=true&sort_order=desc", apiURL)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create request:", err)
			os.Exit(1)
		}
		req.Header.Set("accept", "application/json")
		req.Header.Set("x-litellm-api-key", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Request failed:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "API error: %s\n%s\n", resp.Status, string(body))
			os.Exit(1)
		}

		type KeyObject struct {
			Token     string   `json:"token"`
			KeyName   string   `json:"key_name"`
			KeyAlias  string   `json:"key_alias"`
			Spend     float64  `json:"spend"`
			Models    []string `json:"models"`
			TeamID    *string  `json:"team_id"`
			CreatedAt string   `json:"created_at"`
		}

		type KeyListFullResponse struct {
			Keys        []KeyObject `json:"keys"`
			TotalCount  int         `json:"total_count"`
			CurrentPage int         `json:"current_page"`
			TotalPages  int         `json:"total_pages"`
		}

		var fullResult KeyListFullResponse
		if err := json.NewDecoder(resp.Body).Decode(&fullResult); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decode response:", err)
			os.Exit(1)
		}

		// Output as table matching user keys output
		fmt.Printf("  %-10s | %-10s | %-4s | %-7s | %-15s | %s\n", "KEY NAME", "ALIAS", "TEAM", "SPEND", "MODELS", "CREATED")
		fmt.Println("-------------+------------+------+---------+-----------------+----------------------")
		for _, key := range fullResult.Keys {
			// Mask key name as sk-...Oktg
			masked := key.KeyName
			if len(masked) > 8 {
				masked = masked[:2] + "-..." + masked[len(masked)-4:]
			}
			alias := key.KeyAlias
			if alias == "" {
				alias = "-"
			}
			team := "-"
			if key.TeamID != nil && *key.TeamID != "" {
				team = *key.TeamID
			}
			spend := fmt.Sprintf("$%.2f", key.Spend)
			models := "-"
			if len(key.Models) > 0 {
				models = key.Models[0]
			}
			created := key.CreatedAt
			if len(created) > 19 {
				created = created[:19]
			}
			fmt.Printf("  %-10s | %-10s | %-4s | %-7s | %-15s | %s\n", masked, alias, team, spend, models, created)
		}
	},
}

func init() {
	keyCmd.AddCommand(keyListCmd)
}
