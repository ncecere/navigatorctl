package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ncecere/navigatorctl/pkg/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "List API keys for a team",
	Long: `List all API keys associated with a team and their details.
	
Example:
  # Using team ID
  navigatorctl team keys --team-id 8fb54f04-e5e3-4409-9dd0-262091e5a671
  
  # Using team alias
  navigatorctl team keys --team-alias CHAT
  
  # Using JSON output
  navigatorctl team keys --team-alias CHAT --output json`,
	Run: listKeys,
}

func init() {
	teamCmd.AddCommand(listKeysCmd)
}

func listKeys(cmd *cobra.Command, args []string) {
	teamID := getTeamIdentifier(cmd)
	format := getOutputFormat(cmd)

	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	keys, err := client.ListTeamKeys(teamID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing team keys: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "json":
		outputKeysJSON(keys)
	case "table":
		outputKeysTable(keys)
	}
}

func outputKeysJSON(keys []api.KeyResponse) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(keys); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputKeysTable(keys []api.KeyResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key Name", "Alias", "User ID", "Spend", "Models", "Created At"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)

	for _, key := range keys {
		// Parse the timestamp and format it nicely
		createdAt := key.Info.CreatedAt
		if t, err := time.Parse(time.RFC3339, key.Info.CreatedAt); err == nil {
			createdAt = t.Format("2006-01-02 15:04:05")
		}

		// Format models list
		models := "all-team-models"
		if len(key.Info.Models) > 0 && key.Info.Models[0] != "all-team-models" {
			models = fmt.Sprintf("%d models", len(key.Info.Models))
		}

		// Format user ID
		userID := key.Info.UserID
		if userID == "" {
			userID = "-"
		}

		table.Append([]string{
			key.Info.KeyName,
			key.Info.KeyAlias,
			userID,
			fmt.Sprintf("$%.2f", key.Info.Spend),
			models,
			createdAt,
		})
	}

	table.Render()
}
