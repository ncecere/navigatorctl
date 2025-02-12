package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ncecere/navigatorctl/pkg/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List teams a user belongs to",
	Long: `Display all teams that a user is a member of, including:
- Team ID and alias
- User's role in each team
- Available models
	
Example:
  # Using user ID
  navigatorctl user teams --user-id test-user

  # Using email
  navigatorctl user teams --email user@example.com

  # Using JSON output
  navigatorctl user teams --user-id test-user --output json`,
	Run: showUserTeams,
}

func init() {
	userCmd.AddCommand(userTeamsCmd)
}

func showUserTeams(cmd *cobra.Command, args []string) {
	identifier := getUserIdentifier(cmd)
	if identifier == "" {
		fmt.Fprintln(os.Stderr, "Error: either --user-id or --email is required")
		os.Exit(1)
	}

	format := getOutputFormat(cmd)

	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	response, err := client.GetUserInfo(identifier)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user teams: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "json":
		outputUserTeamsJSON(response)
	case "table":
		outputUserTeamsTable(response)
	}
}

func outputUserTeamsJSON(response *api.UserResponse) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response.Teams); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputUserTeamsTable(response *api.UserResponse) {
	if len(response.Teams) == 0 {
		fmt.Println("User is not a member of any teams")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Team ID", "Alias", "Role", "Models", "Spend"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)

	for _, team := range response.Teams {
		for _, member := range team.MembersWithRoles {
			if member.UserID == response.UserInfo.UserID {
				models := "all-team-models"
				if len(team.Models) > 0 && team.Models[0] != "all-team-models" {
					models = fmt.Sprintf("%d models", len(team.Models))
				}

				table.Append([]string{
					team.TeamID,
					getOrDefault(team.TeamAlias, "-"),
					member.Role,
					models,
					fmt.Sprintf("$%.2f", team.Spend),
				})
			}
		}
	}

	table.Render()
}
