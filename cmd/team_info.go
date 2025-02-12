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

var teamInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display team information",
	Long: `Display detailed information about a team including:
- Team ID and alias
- Available models
- Current spend
- Creation date
	
Example:
  # Using team ID
  navigatorctl team info --team-id 0dbaa4dd-8523-4e05-8d43-91b7dd80f671
  
  # Using team alias
  navigatorctl team info --team-id CLINE
  
  # Using JSON output
  navigatorctl team info --team-id CLINE --output json`,
	Run: showTeamInfo,
}

func init() {
	teamCmd.AddCommand(teamInfoCmd)
}

func showTeamInfo(cmd *cobra.Command, args []string) {
	teamID := getTeamIdentifier(cmd)
	format := getOutputFormat(cmd)

	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	team, err := client.GetTeamInfo(teamID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting team info: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "json":
		outputTeamInfoJSON(team)
	case "table":
		outputTeamInfoTable(team)
	}
}

func outputTeamInfoJSON(team *api.Team) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(team); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputTeamInfoTable(team *api.Team) {
	// Create table for general info
	infoTable := tablewriter.NewWriter(os.Stdout)
	infoTable.SetHeader([]string{"Field", "Value"})
	infoTable.SetBorder(false)
	infoTable.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	infoTable.SetAlignment(tablewriter.ALIGN_LEFT)
	infoTable.SetAutoWrapText(false)

	infoTable.Append([]string{"Team ID", team.TeamID})
	infoTable.Append([]string{"Alias", team.TeamAlias})
	infoTable.Append([]string{"Spend", fmt.Sprintf("$%.2f", team.Spend)})
	infoTable.Append([]string{"Created At", team.CreatedAt})
	infoTable.Append([]string{"Models", fmt.Sprintf("%v", team.Models)})

	fmt.Println("Team Information:")
	infoTable.Render()
}
