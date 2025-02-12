package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ncecere/navigatorctl/pkg/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "List API keys associated with a user",
	Long: `Display all API keys associated with a user, including:
- Key name and alias
- Spend and budget information
- Associated team
- Available models
	
Example:
  # Using user ID
  navigatorctl user keys --user-id test-user

  # Using email
  navigatorctl user keys --email user@example.com

  # Using JSON output
  navigatorctl user keys --user-id test-user --output json`,
	Run: showUserKeys,
}

func init() {
	userCmd.AddCommand(userKeysCmd)
}

func showUserKeys(cmd *cobra.Command, args []string) {
	identifier := getUserIdentifier(cmd)
	if identifier == "" {
		fmt.Fprintln(os.Stderr, "Error: either --user-id or --email is required")
		os.Exit(1)
	}

	format := getOutputFormat(cmd)

	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	response, err := client.GetUserInfo(identifier)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user keys: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "json":
		outputUserKeysJSON(response)
	case "table":
		outputUserKeysTable(response)
	}
}

func outputUserKeysJSON(response *api.UserResponse) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response.Keys); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputUserKeysTable(response *api.UserResponse) {
	if len(response.Keys) == 0 {
		fmt.Println("No API keys found for user")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key Name", "Alias", "Team", "Spend", "Models", "Created"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)

	for _, key := range response.Keys {
		// Format team info
		team := "-"
		if key.TeamID != "" {
			for _, t := range response.Teams {
				if t.TeamID == key.TeamID {
					team = getOrDefault(t.TeamAlias, key.TeamID)
					break
				}
			}
		}

		// Format models
		models := "all-team-models"
		if len(key.Models) > 0 && key.Models[0] != "all-team-models" {
			models = strings.Join(key.Models, ", ")
		}

		// Format creation date
		created := key.CreatedAt
		if t, err := time.Parse(time.RFC3339, key.CreatedAt); err == nil {
			created = t.Format("2006-01-02 15:04:05")
		}

		table.Append([]string{
			key.KeyName,
			getOrDefault(key.KeyAlias, "-"),
			team,
			fmt.Sprintf("$%.2f", key.Spend),
			models,
			created,
		})
	}

	table.Render()
}
