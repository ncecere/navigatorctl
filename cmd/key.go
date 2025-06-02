// cmd/key.go

package cmd

import (
	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage API keys",
	Long:  "List and get information about API keys.",
}

func init() {
	rootCmd.AddCommand(keyCmd)
}
