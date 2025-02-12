package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage teams and their members",
	Long: `Team management commands allow you to:
- List all teams and their information
- List, add, and remove team members
- List team API keys
- View team information and budgets

Examples:
  # List all teams
  navigatorctl team list
  navigatorctl team list --output json

  # List team members (using ID or alias)
  navigatorctl team members --team-id 0dbaa4dd-8523-4e05-8d43-91b7dd80f671
  navigatorctl team members --team-id CLINE

  # Add a team member (using ID or alias)
  navigatorctl team add-member --team-id CLINE --user-id user_456 --role admin

  # View team information (using ID or alias)
  navigatorctl team info --team-id CLINE
  navigatorctl team info --team-id 0dbaa4dd-8523-4e05-8d43-91b7dd80f671`,
}

func init() {
	rootCmd.AddCommand(teamCmd)

	// Global flags for team commands
	teamCmd.PersistentFlags().StringP("team-id", "t", "", "Team ID to perform operations on")
	teamCmd.PersistentFlags().StringP("team-alias", "a", "", "Team alias to perform operations on")
	teamCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table, json)")

	// Bind flags to viper
	if err := viper.BindPFlag("team.id", teamCmd.PersistentFlags().Lookup("team-id")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding team-id flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("team.alias", teamCmd.PersistentFlags().Lookup("team-alias")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding team-alias flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("output.format", teamCmd.PersistentFlags().Lookup("output")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding output flag: %v\n", err)
		os.Exit(1)
	}
}

func getTeamIdentifier(cmd *cobra.Command) string {
	teamID := viper.GetString("team.id")
	teamAlias := viper.GetString("team.alias")

	if teamID == "" && teamAlias == "" {
		fmt.Fprintln(os.Stderr, "Error: either --team-id or --team-alias is required")
		cmd.Help()
		os.Exit(1)
	}

	// Prefer team ID if both are provided
	if teamID != "" {
		return teamID
	}
	return teamAlias
}

func getOutputFormat(cmd *cobra.Command) string {
	// First try to get from flag
	format, _ := cmd.Flags().GetString("output")
	if format == "" {
		// If not in flag, try config
		format = viper.GetString("output.format")
	}
	// Default to table if still empty
	if format == "" {
		format = "table"
	}

	if format != "table" && format != "json" {
		fmt.Fprintf(os.Stderr, "Error: invalid output format '%s'. Must be 'table' or 'json'\n", format)
		cmd.Help()
		os.Exit(1)
	}
	return format
}
