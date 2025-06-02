// cmd/model_info.go

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type ModelInfoParams struct {
	InputCostPerToken  float64 `json:"input_cost_per_token"`
	OutputCostPerToken float64 `json:"output_cost_per_token"`
	ApiBase            string  `json:"api_base"`
	ApiVersion         string  `json:"api_version"`
	CustomProvider     string  `json:"custom_llm_provider"`
}

type ModelInfoDetails struct {
	ID                string `json:"id"`
	BaseModel         string `json:"base_model"`
	Tier              string `json:"tier"`
	Mode              string `json:"mode"`
	MaxTokens         int    `json:"max_tokens"`
	LitellmProvider   string `json:"litellm_provider"`
	SupportsVision    bool   `json:"supports_vision"`
	SupportsFunction  bool   `json:"supports_function_calling"`
	SupportsTool      bool   `json:"supports_tool_choice"`
	SupportsStreaming bool   `json:"supports_native_streaming"`
}

type ModelInfoItem struct {
	ModelName     string           `json:"model_name"`
	LitellmParams ModelInfoParams  `json:"litellm_params"`
	ModelInfo     ModelInfoDetails `json:"model_info"`
}

type ModelInfoResponse struct {
	Data []ModelInfoItem `json:"data"`
}

var modelInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show detailed info for all models",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api-url")
		apiKey, _ := cmd.Flags().GetString("api-key")
		output, _ := cmd.Flags().GetString("output")
		modelFilter, _ := cmd.Flags().GetString("model")
		if apiURL == "" || apiKey == "" {
			fmt.Fprintln(os.Stderr, "API URL and API Key are required")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/model/info", apiURL)
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

		var result ModelInfoResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decode response:", err)
			os.Exit(1)
		}

		// Filter by model if --model is provided
		filtered := result.Data
		if modelFilter != "" {
			filtered = []ModelInfoItem{}
			for _, m := range result.Data {
				if m.ModelName == modelFilter || m.ModelInfo.ID == modelFilter {
					filtered = append(filtered, m)
				}
			}
		}

		switch output {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(filtered); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
		default:
			// Helper to truncate long values for table columns
			trunc := func(s string, n int) string {
				if len(s) > n {
					return s[:n-1] + "…"
				}
				return s
			}
			fmt.Printf("%-18s | %-8s | %-8s | %-10s | %-12s | %-6s | %-6s | %-6s | %-6s\n", "MODEL", "TIER", "MODE", "MAX TOKENS", "PROVIDER", "VISION", "FUNC", "TOOL", "STREAM")
			fmt.Println("--------------------+----------+----------+------------+--------------+--------+--------+--------+--------")
			for _, m := range filtered {
				maxTokens := "-"
				if m.ModelInfo.MaxTokens > 0 {
					maxTokens = fmt.Sprintf("%d", m.ModelInfo.MaxTokens)
				}
				fmt.Printf("%-18s | %-8s | %-8s | %-10s | %-12s | %-6t | %-6t | %-6t | %-6t\n",
					trunc(m.ModelName, 18),
					trunc(m.ModelInfo.Tier, 8),
					trunc(m.ModelInfo.Mode, 8),
					maxTokens,
					trunc(m.ModelInfo.LitellmProvider, 12),
					m.ModelInfo.SupportsVision,
					m.ModelInfo.SupportsFunction,
					m.ModelInfo.SupportsTool,
					m.ModelInfo.SupportsStreaming,
				)
			}
		}
	},
}

func init() {
	modelInfoCmd.Flags().String("output", "table", "Output format: table or json")
	modelInfoCmd.Flags().String("model", "", "Model name or ID to filter")
	modelCmd.AddCommand(modelInfoCmd)
}
