package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users and their information",
	Long: `User management commands allow you to:
- View user information and settings
- List user's teams and API keys
- View user's spend and budget information`,
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Global flags for user commands
	userCmd.PersistentFlags().StringP("user-id", "u", "", "User ID to perform operations on")
	userCmd.PersistentFlags().StringP("email", "e", "", "User email to perform operations on")
	userCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table, json)")

	// Bind flags to viper
	if err := viper.BindPFlag("user.id", userCmd.PersistentFlags().Lookup("user-id")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("user.email", userCmd.PersistentFlags().Lookup("email")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("output.format", userCmd.PersistentFlags().Lookup("output")); err != nil {
		panic(err)
	}
}

func getUserIdentifier(cmd *cobra.Command) string {
	userID := viper.GetString("user.id")
	if userID != "" {
		return userID
	}

	email := viper.GetString("user.email")
	if email != "" {
		return email
	}

	cmd.Help()
	return ""
}
