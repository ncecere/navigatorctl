// cmd/model_list.go

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type ModelListItem struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelListResponse struct {
	Data   []ModelListItem `json:"data"`
	Object string          `json:"object"`
}

var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available models",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api-url")
		apiKey, _ := cmd.Flags().GetString("api-key")
		output, _ := cmd.Flags().GetString("output")
		if apiURL == "" || apiKey == "" {
			fmt.Fprintln(os.Stderr, "API URL and API Key are required")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/models?return_wildcard_routes=false&include_model_access_groups=false", apiURL)
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

		var result ModelListResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decode response:", err)
			os.Exit(1)
		}

		switch output {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(result.Data); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("%-25s | %-10s | %-20s\n", "ID", "OWNER", "CREATED")
			fmt.Println("-------------------------+------------+----------------------")
			for _, m := range result.Data {
				created := time.Unix(m.Created, 0).Format("2006-01-02 15:04:05")
				fmt.Printf("%-25s | %-10s | %-20s\n", m.ID, m.OwnedBy, created)
			}
		}
	},
}

func init() {
	modelListCmd.Flags().String("output", "table", "Output format: table or json")
	modelCmd.AddCommand(modelListCmd)
}
