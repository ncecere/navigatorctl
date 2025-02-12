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

var userInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display user information",
	Long: `Display detailed information about a user including:
- User ID and email
- Team memberships
- API keys
- Spend and budget information
- Created and updated dates

Example:
  # Using user ID
  navigatorctl user info --user-id test-user

  # Using email
  navigatorctl user info --email user@example.com

  # Using JSON output
  navigatorctl user info --user-id test-user --output json`,
	Run: showUserInfo,
}

func init() {
	userCmd.AddCommand(userInfoCmd)
}

func showUserInfo(cmd *cobra.Command, args []string) {
	identifier := getUserIdentifier(cmd)
	if identifier == "" {
		fmt.Fprintln(os.Stderr, "Error: either --user-id or --email is required")
		os.Exit(1)
	}

	format := getOutputFormat(cmd)

	client := api.NewClient(viper.GetString("api.url"), viper.GetString("api.key"))

	response, err := client.GetUserInfo(identifier)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user info: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "json":
		outputUserInfoJSON(response)
	case "table":
		outputUserInfoTable(response)
	}
}

func outputUserInfoJSON(response *api.UserResponse) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputUserInfoTable(response *api.UserResponse) {
	if response.UserInfo == nil {
		fmt.Println("No user information available")
		return
	}

	// User Info Table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)

	// Format dates
	createdAt := response.UserInfo.CreatedAt
	if t, err := time.Parse(time.RFC3339, response.UserInfo.CreatedAt); err == nil {
		createdAt = t.Format("2006-01-02 15:04:05")
	}
	updatedAt := response.UserInfo.UpdatedAt
	if t, err := time.Parse(time.RFC3339, response.UserInfo.UpdatedAt); err == nil {
		updatedAt = t.Format("2006-01-02 15:04:05")
	}

	table.Append([]string{"User ID", response.UserInfo.UserID})
	table.Append([]string{"Email", getOrDefault(response.UserInfo.UserEmail, "-")})
	table.Append([]string{"Role", response.UserInfo.UserRole})
	table.Append([]string{"Spend", fmt.Sprintf("$%.2f", response.UserInfo.Spend)})
	if response.UserInfo.MaxBudget > 0 {
		table.Append([]string{"Max Budget", fmt.Sprintf("$%.2f", response.UserInfo.MaxBudget)})
	}
	table.Append([]string{"Created At", createdAt})
	table.Append([]string{"Updated At", updatedAt})

	fmt.Println("User Information:")
	table.Render()
	fmt.Println()

}

func getOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
