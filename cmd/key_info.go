// cmd/key_info.go

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type KeyInfo struct {
	KeyName   string  `json:"key_name"`
	KeyAlias  string  `json:"key_alias"`
	Spend     float64 `json:"spend"`
	CreatedAt string  `json:"created_at"`
}

type KeyInfoResponse struct {
	Key  string  `json:"key"`
	Info KeyInfo `json:"info"`
}

var keyInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get API key info",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api-url")
		apiKey, _ := cmd.Flags().GetString("api-key")
		key, _ := cmd.Flags().GetString("key")
		if apiURL == "" || apiKey == "" || key == "" {
			fmt.Fprintln(os.Stderr, "API URL, API Key, and --key are required")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/key/info?key=%s", apiURL, key)
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

		var result KeyInfoResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to decode response:", err)
			os.Exit(1)
		}

		// Output key info (name, alias, spend, created_at)
		fmt.Println("Key Info:")
		fmt.Println("Name: ", result.Info.KeyName)
		fmt.Println("Alias:", result.Info.KeyAlias)
		fmt.Println("Spend:", result.Info.Spend)
		fmt.Println("Created At:", result.Info.CreatedAt)
	},
}

func init() {
	keyInfoCmd.Flags().String("key", "", "API key string to get info for")
	keyCmd.AddCommand(keyInfoCmd)
}
