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

var (
	listMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "List members in a team",
		Long: `List all members in a team and their roles.
	
Example:
  # Using team ID
  navigatorctl team members --team-id 0dbaa4dd-8523-4e05-8d43-91b7dd80f671
  
  # Using team alias
  navigatorctl team members --team-alias CLINE
  
  # Using JSON output
  navigatorctl team members --team-alias CLINE --output json`,
		Run: listMembers,
	}

	addMemberCmd = &cobra.Command{
		Use:   "add-member",
		Short: "Add a member to a team",
		Long: `Add a new member to a team with specified role.
	
Example:
  # Using team ID
  navigatorctl team add-member --team-id 0dbaa4dd-8523-4e05-8d43-91b7dd80f671 --user-id user_456 --role admin
  
  # Using team alias
  navigatorctl team add-member --team-alias CLINE --user-id user_456 --role member`,
		Run: addMember,
	}

	removeMemberCmd = &cobra.Command{
		Use:   "remove-member",
		Short: "Remove a member from a team",
		Long: `Remove a member from a team.
	
Example:
  # Using team ID
  navigatorctl team remove-member --team-id 0dbaa4dd-8523-4e05-8d43-91b7dd80f671 --user-id user_456
  
  # Using team alias
  navigatorctl team remove-member --team-alias CLINE --user-id user_456`,
		Run: removeMember,
	}
)

func init() {
	teamCmd.AddCommand(listMembersCmd)
	teamCmd.AddCommand(addMemberCmd)
	teamCmd.AddCommand(removeMemberCmd)

	// Add flags for member management
	addMemberCmd.Flags().StringP("user-id", "u", "", "User ID to add")
	addMemberCmd.Flags().StringP("email", "e", "", "User email address")
	addMemberCmd.Flags().StringP("role", "r", "", "Role to assign (admin/user)")
	addMemberCmd.MarkFlagRequired("user-id")
	addMemberCmd.MarkFlagRequired("role")

	removeMemberCmd.Flags().StringP("user-id", "u", "", "User ID to remove")
	removeMemberCmd.Flags().StringP("email", "e", "", "User email address")
	removeMemberCmd.MarkFlagRequired("user-id")
}

func addMember(cmd *cobra.Command, args []string) {
	teamID := getTeamIdentifier(cmd)
	userID, _ := cmd.Flags().GetString("user-id")
	role, _ := cmd.Flags().GetString("role")

	if role != "admin" && role != "user" {
		fmt.Fprintln(os.Stderr, "Error: role must be either 'admin' or 'user'")
		os.Exit(1)
	}

	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	email, _ := cmd.Flags().GetString("email")

	member := api.TeamMember{
		UserID:    userID,
		UserEmail: email,
		Role:      role,
	}

	response, err := client.AddTeamMember(teamID, member)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding team member: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully added user %s to team %s (%s) with role %s\n",
		userID, response.TeamID, response.TeamAlias, role)
}

func removeMember(cmd *cobra.Command, args []string) {
	teamID := getTeamIdentifier(cmd)
	userID, _ := cmd.Flags().GetString("user-id")

	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	email, _ := cmd.Flags().GetString("email")

	member := api.TeamMember{
		UserID:    userID,
		UserEmail: email,
	}

	response, err := client.RemoveTeamMember(teamID, member)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing team member: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully removed user %s from team %s (%s)\n",
		userID, response.TeamID, response.TeamAlias)
}

func listMembers(cmd *cobra.Command, args []string) {
	teamID := getTeamIdentifier(cmd)
	format := getOutputFormat(cmd)

	// Get API client from root command
	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	members, err := client.ListTeamMembers(teamID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing team members: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "json":
		outputJSON(members)
	case "table":
		outputTable(members)
	}
}

func outputJSON(members []api.TeamMember) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(members); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputTable(members []api.TeamMember) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"User ID", "Email", "Role"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, member := range members {
		// Format email
		email := member.UserEmail
		if email == "" {
			email = "-"
		}

		table.Append([]string{
			member.UserID,
			email,
			member.Role,
		})
	}

	table.Render()
}
