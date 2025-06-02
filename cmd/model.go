// cmd/model.go

package cmd

import (
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage and inspect available models",
	Long:  "List, inspect, and check health of available models.",
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
