package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all teams",
	Long: `List all teams with their IDs and names.

Example:
  navigatorctl team list
  navigatorctl team list --output json`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getAPIClient()
		teams, err := client.ListTeams()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing teams: %v\n", err)
			os.Exit(1)
		}

		outputFormat := getOutputFormat(cmd)
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(teams); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
		case "table":
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Team ID", "Alias", "Models", "Spend"})

			for _, team := range teams {
				table.Append([]string{
					team.TeamID,
					team.TeamAlias,
					fmt.Sprintf("%d models", len(team.Models)),
					fmt.Sprintf("$%.2f", team.Spend),
				})
			}

			table.Render()
		}
	},
}

func init() {
	teamCmd.AddCommand(listCmd)
}
