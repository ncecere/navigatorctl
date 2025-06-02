// cmd/model_health.go

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type HealthEndpoint struct {
	ApiBase                   string `json:"api_base"`
	ApiVersion                string `json:"api_version"`
	CustomProvider            string `json:"custom_llm_provider"`
	XMsRegion                 string `json:"x-ms-region"`
	XRateLimitRemainingReqs   string `json:"x-ratelimit-remaining-requests"`
	XRateLimitRemainingTokens string `json:"x-ratelimit-remaining-tokens"`
}

type ModelHealthResponse struct {
	HealthyEndpoints   []HealthEndpoint `json:"healthy_endpoints"`
	UnhealthyEndpoints []HealthEndpoint `json:"unhealthy_endpoints"`
	HealthyCount       int              `json:"healthy_count"`
	UnhealthyCount     int              `json:"unhealthy_count"`
}

var modelHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Show health and endpoint status for a specific model",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api-url")
		apiKey, _ := cmd.Flags().GetString("api-key")
		model, _ := cmd.Flags().GetString("model")
		output, _ := cmd.Flags().GetString("output")
		if apiURL == "" || apiKey == "" || model == "" {
			fmt.Fprintln(os.Stderr, "API URL, API Key, and --model are required")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/health?model=%s", apiURL, model)
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

		var result ModelHealthResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decode response:", err)
			os.Exit(1)
		}

		switch output {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(result); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Println("Healthy Endpoints:")
			if len(result.HealthyEndpoints) == 0 {
				fmt.Println("  None")
			} else {
				fmt.Printf("%-40s | %-10s | %-10s | %-10s | %-10s\n", "API BASE", "REGION", "REQ LEFT", "TOKENS LEFT", "PROVIDER")
				fmt.Println("------------------------------------------+------------+------------+------------+------------")
				for _, ep := range result.HealthyEndpoints {
					fmt.Printf("%-40s | %-10s | %-10s | %-10s | %-10s\n",
						ep.ApiBase, ep.XMsRegion, ep.XRateLimitRemainingReqs, ep.XRateLimitRemainingTokens, ep.CustomProvider)
				}
			}
			fmt.Println()
			fmt.Println("Unhealthy Endpoints:")
			if len(result.UnhealthyEndpoints) == 0 {
				fmt.Println("  None")
			} else {
				fmt.Printf("%-40s | %-10s | %-10s | %-10s | %-10s\n", "API BASE", "REGION", "REQ LEFT", "TOKENS LEFT", "PROVIDER")
				fmt.Println("------------------------------------------+------------+------------+------------+------------")
				for _, ep := range result.UnhealthyEndpoints {
					fmt.Printf("%-40s | %-10s | %-10s | %-10s | %-10s\n",
						ep.ApiBase, ep.XMsRegion, ep.XRateLimitRemainingReqs, ep.XRateLimitRemainingTokens, ep.CustomProvider)
				}
			}
		}
	},
}

func init() {
	modelHealthCmd.Flags().String("model", "", "Model ID to check health for")
	modelHealthCmd.Flags().String("output", "table", "Output format: table or json")
	modelCmd.AddCommand(modelHealthCmd)
}
